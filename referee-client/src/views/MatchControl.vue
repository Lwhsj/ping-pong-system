<template>
  <section class="match-control">
    <div class="control-header">
      <div>
        <p class="section-kicker">Live Match</p>
        <h2>比赛控制台</h2>
      </div>

      <div class="status-cluster">
        <span class="status-pill">状态：{{ statusText }}</span>
        <span v-if="isRecording" class="status-pill recording">录制中</span>
        <span v-if="loading" class="status-pill processing">处理中...</span>
      </div>
    </div>

    <div v-if="error || videoError" class="message-stack">
      <div v-if="error" class="alert alert-error">API 错误：{{ error }}</div>
      <div v-if="videoError" class="alert alert-error">视频错误：{{ videoError }}</div>
    </div>

    <div class="scoreboard">
      <div class="player-score">
        <p class="player-label">选手 1</p>
        <h3>{{ matchInfo.player1Name || '选手 1' }}</h3>
        <div class="score">{{ score.p1 }}</div>
        <button @click="scorePoint('player1')" :disabled="scoreLocked" class="btn-score">
          记为得分
        </button>
      </div>

      <div class="match-center">
        <div class="vs">VS</div>
        <div class="info-card">
          <span>回合</span>
          <strong>{{ rallyNumber }}</strong>
        </div>
        <div class="info-card server-card">
          <span>当前发球方</span>
          <strong>{{ currentServerName || '待同步' }}</strong>
        </div>
      </div>

      <div class="player-score accent">
        <p class="player-label">选手 2</p>
        <h3>{{ matchInfo.player2Name || '选手 2' }}</h3>
        <div class="score">{{ score.p2 }}</div>
        <button @click="scorePoint('player2')" :disabled="scoreLocked" class="btn-score">
          记为得分
        </button>
      </div>
    </div>

    <div class="actions">
      <button @click="finishMatch" :disabled="finishLocked" class="btn-danger">
        结束比赛
      </button>
      <button @click="exportExcel" :disabled="matchInfo.status !== 'finished'" class="btn-secondary">
        导出 Excel
      </button>
    </div>

    <section class="agent-panel" aria-label="AI 复盘助手">
      <div class="agent-header">
        <div>
          <p class="section-kicker">AI Match Desk</p>
          <h3>AI 复盘助手</h3>
        </div>
        <span class="agent-status">{{ agentStatusText }}</span>
      </div>

      <div class="agent-grid">
        <div class="analysis-box">
          <div v-if="showFinalAnalysisPrompt" class="final-analysis-callout">
            <div>
              <strong>比赛已结束</strong>
              <span>可以生成一份完整赛后复盘。</span>
            </div>
            <button type="button" @click="generateFinalAnalysis" :disabled="agentLoading">
              生成完整复盘
            </button>
          </div>

          <label for="analysis-focus">复盘关注点</label>
          <textarea
            id="analysis-focus"
            v-model="analysisQuestion"
            rows="3"
            placeholder="例如：重点分析后半段连续丢分、发球轮表现、关键分处理"
          ></textarea>

          <div class="quick-actions" aria-label="复盘关注点快捷选择">
            <button
              v-for="item in analysisPrompts"
              :key="item"
              type="button"
              class="quick-chip"
              @click="analysisQuestion = item"
            >
              {{ item }}
            </button>
          </div>

          <button @click="analyzeMatch" :disabled="agentLoading" class="btn-agent">
            {{ agentLoading ? '生成中...' : '生成 AI 复盘' }}
          </button>

          <div v-if="analysisError" class="alert alert-error agent-alert">{{ analysisError }}</div>

          <div v-if="analysis" class="analysis-result">
            <div v-if="analysisOutdated" class="alert alert-warning">
              该复盘基于 {{ analysisScoreLabel }} 生成，当前比分已变为 {{ currentScoreLabel }}。
              <button type="button" @click="analyzeMatch" :disabled="agentLoading">重新生成</button>
            </div>

            <div class="summary-card">
              <span>总结</span>
              <p>{{ analysis.summary || analysis.raw_text }}</p>
              <small v-if="analysisSnapshot">{{ analysisSnapshot }}</small>
              <button type="button" class="copy-review" @click="copyAnalysis">
                {{ copyStatus || '复制复盘' }}
              </button>
            </div>

            <div v-if="analysis.strengths?.length" class="insight-section">
              <h4>优势</h4>
              <ul>
                <li v-for="item in analysis.strengths" :key="item">{{ item }}</li>
              </ul>
            </div>

            <div v-if="analysis.weaknesses?.length" class="insight-section">
              <h4>问题</h4>
              <ul>
                <li v-for="item in analysis.weaknesses" :key="item">{{ item }}</li>
              </ul>
            </div>

            <div v-if="analysis.key_moments?.length" class="key-moments">
              <h4>关键回合</h4>
              <div v-for="moment in analysis.key_moments" :key="`${moment.rally_number}-${moment.reason}`" class="moment-row">
                <strong>#{{ moment.rally_number }}</strong>
                <span>{{ moment.reason }}</span>
              </div>
            </div>

            <div v-if="analysis.training_suggestions?.length" class="insight-section training">
              <h4>训练建议</h4>
              <ul>
                <li v-for="item in analysis.training_suggestions" :key="item">{{ item }}</li>
              </ul>
            </div>

            <div class="follow-up-panel">
              <h4>继续追问</h4>
              <div class="quick-actions">
                <button
                  v-for="item in followUpQuestions"
                  :key="item"
                  type="button"
                  class="quick-chip"
                  @click="askSuggestedQuestion(item)"
                >
                  {{ item }}
                </button>
              </div>
            </div>
          </div>

          <div v-else class="empty-agent-state">
            完成几个回合后，可以生成基于比分走势、发球轮和连续得分的复盘。
          </div>
        </div>

        <div class="chat-box">
          <div class="chat-header">
            <h4>比赛问答</h4>
            <span>{{ chatMessages.length }} 条</span>
          </div>

          <div class="chat-log">
            <div v-if="!chatMessages.length" class="chat-empty">
              <span>可以直接点一个问题开始。</span>
              <div class="quick-actions">
                <button
                  v-for="item in suggestedQuestions"
                  :key="item"
                  type="button"
                  class="quick-chip"
                  @click="askSuggestedQuestion(item)"
                >
                  {{ item }}
                </button>
              </div>
            </div>
            <div v-for="message in chatMessages" :key="message.id" class="chat-message" :class="message.role">
              <span>{{ message.role === 'user' ? '我' : 'AI' }}</span>
              <p>{{ message.content }}</p>
              <small v-if="message.snapshot">{{ message.snapshot }}</small>
            </div>
          </div>

          <div v-if="chatError" class="alert alert-error chat-alert">{{ chatError }}</div>

          <form class="chat-form" @submit.prevent="sendAgentQuestion">
            <input
              v-model="chatQuestion"
              type="text"
              placeholder="围绕这场比赛提问"
              :disabled="chatLoading"
            />
            <button type="submit" :disabled="chatLoading || !chatQuestion.trim()" class="btn-secondary">
              {{ chatLoading ? '发送中' : '发送' }}
            </button>
          </form>
        </div>
      </div>
    </section>
  </section>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useVideoRecorder } from '../composables/useVideoRecorder'
import api from '../services/api'

const route = useRoute()
const router = useRouter()
const { startRecording, stopRecording, stopAndGetBlob, isRecording, error: videoError } = useVideoRecorder()

const matchInfo = reactive({
  id: '',
  player1Id: '',
  player2Id: '',
  player1Name: '',
  player2Name: '',
  firstServerId: '',
  status: ''
})

const score = reactive({
  p1: 0,
  p2: 0
})

const rallyNumber = ref(1)
const currentServerId = ref('')
const loading = ref(false)
const error = ref(null)
const analysisQuestion = ref('')
const analysis = ref(null)
const agentLoading = ref(false)
const analysisError = ref(null)
const analysisSnapshot = ref('')
const analysisScore = reactive({
  p1: null,
  p2: null
})
const copyStatus = ref('')
const chatQuestion = ref('')
const chatMessages = ref([])
const chatLoading = ref(false)
const chatError = ref(null)
let chatMessageId = 0

const analysisPrompts = [
  '重点分析发球表现',
  '找出连续失分的原因',
  '分析关键分处理',
  '给出下一场训练建议',
  '判断当前谁更占优'
]

const liveSuggestedQuestions = [
  '当前谁更占优？',
  '下一分发球方是谁？',
  '刚才连续丢分可能是什么原因？',
  '现在应该优先稳住什么？'
]

const finishedSuggestedQuestions = [
  '这场输赢关键是什么？',
  '下一场应该重点练什么？',
  '谁的发球轮表现更好？',
  '这场最需要复盘哪几个回合？'
]

const defaultSuggestedQuestions = [
  '这场谁的发球优势更明显？',
  '我在哪些回合开始丢节奏？',
  '后半段连续丢分的原因是什么？',
  '下一场应该重点练什么？'
]

const followUpQuestions = [
  '为什么这些回合是转折点？',
  '只看选手 1 的问题',
  '给我 3 个可执行训练动作'
]

const currentServerName = computed(() => {
  if (currentServerId.value === matchInfo.player1Id) return matchInfo.player1Name
  if (currentServerId.value === matchInfo.player2Id) return matchInfo.player2Name
  return ''
})

const statusText = computed(() => {
  if (matchInfo.status === 'finished') return '已结束'
  if (['ongoing', 'started', 'active', 'in_progress'].includes(matchInfo.status)) return '进行中'
  return matchInfo.status || '进行中'
})

const isFinished = computed(() => matchInfo.status === 'finished')

const agentStatusText = computed(() => {
  if (agentLoading.value || chatLoading.value) return '思考中'
  if (analysisOutdated.value) return '需更新'
  if (analysis.value) return '已有复盘'
  return '待生成'
})

const currentScoreLabel = computed(() => `${score.p1}-${score.p2}`)

const analysisScoreLabel = computed(() => {
  if (analysisScore.p1 === null || analysisScore.p2 === null) return ''
  return `${analysisScore.p1}-${analysisScore.p2}`
})

const analysisOutdated = computed(() => (
  Boolean(analysis.value) &&
  analysisScore.p1 !== null &&
  analysisScore.p2 !== null &&
  (analysisScore.p1 !== score.p1 || analysisScore.p2 !== score.p2)
))

const suggestedQuestions = computed(() => {
  if (isFinished.value) return finishedSuggestedQuestions
  if (['ongoing', 'started', 'active', 'in_progress'].includes(matchInfo.status)) return liveSuggestedQuestions
  return defaultSuggestedQuestions
})

const showFinalAnalysisPrompt = computed(() => isFinished.value && !analysis.value && !agentLoading.value)

const scoreLocked = computed(() => loading.value || matchInfo.status === 'finished')

const finishLocked = computed(() => loading.value || matchInfo.status === 'finished')

const numericId = (value) => {
  const parsed = Number(value)
  return Number.isFinite(parsed) && parsed > 0 ? parsed : ''
}

const applyMatchInfo = (match) => {
  matchInfo.id = numericId(match.id)
  matchInfo.player1Id = numericId(match.player1_id || match.player1Id || match.player1ID)
  matchInfo.player2Id = numericId(match.player2_id || match.player2Id || match.player2ID)
  matchInfo.player1Name = match.player1_name || match.player1Name || ''
  matchInfo.player2Name = match.player2_name || match.player2Name || ''
  matchInfo.firstServerId = numericId(match.first_server || match.firstServerId)
  matchInfo.status = match.status || ''
  currentServerId.value = matchInfo.firstServerId
}

const loadStoredMatch = () => {
  const storedMatch = localStorage.getItem('currentMatch')
  if (!storedMatch) return null
  try {
    return JSON.parse(storedMatch)
  } catch (err) {
    console.warn('Invalid stored match info', err)
    localStorage.removeItem('currentMatch')
    return null
  }
}

const loadMatchInfo = async () => {
  const routeMatchId = Number(route.params.id || 0)
  if (routeMatchId > 0) {
    const response = await api.getMatch(routeMatchId)
    applyMatchInfo(response.data)
    localStorage.setItem('currentMatch', JSON.stringify(matchInfo))
    return true
  }

  const storedMatch = loadStoredMatch()
  if (storedMatch?.id) {
    applyMatchInfo(storedMatch)
    return true
  }
  return false
}

const fetchMatchState = async () => {
  try {
    const res = await api.getMatchCurrent(matchInfo.id)
    const data = res.data
    rallyNumber.value = data.rally_number + 1 || 1
    score.p1 = data.score_p1 || 0
    score.p2 = data.score_p2 || 0
    matchInfo.player1Name = data.player1_name || matchInfo.player1Name
    matchInfo.player2Name = data.player2_name || matchInfo.player2Name

    // Update current server from API response
    if (data.server === 'player1') {
      currentServerId.value = matchInfo.player1Id
    } else if (data.server === 'player2') {
      currentServerId.value = matchInfo.player2Id
    }
  } catch (e) {
    console.warn('Could not sync with server', e)
    error.value = '同步比赛状态失败'
  }
}

onMounted(async () => {
  try {
    const hasMatch = await loadMatchInfo()
    if (!hasMatch) {
      router.push('/')
      return
    }
    await fetchMatchState()
  } catch (err) {
    console.error('Failed to load match info', err)
    error.value = '加载比赛信息失败'
    router.push('/')
    return
  }

  // Start video recording
  await startRecording()
})

onUnmounted(() => {
  stopRecording()
})

const scorePoint = async (scorer) => {
  if (loading.value) return

  try {
    loading.value = true
    error.value = null

    // 1. Stop current recording and get full video
    console.log('Stopping recording for rally', rallyNumber.value)
    const videoBlob = await stopAndGetBlob()

    // 2. Restart recording for the NEXT rally before upload work continues.
    const recordingRestarted = await startRecording()
    if (!recordingRestarted) {
      videoError.value = videoError.value || '下一回合录制启动失败，请检查摄像头权限'
    }

    let videoFileName = ''

    if (videoBlob) {
       // Upload video
       const formData = new FormData()
       const fileName = `match_${matchInfo.id}_rally${rallyNumber.value}.webm`
       formData.append('file', videoBlob, fileName)
       formData.append('matchId', matchInfo.id)

       try {
         // Upload video
         const res = await api.uploadVideo(formData)
         videoFileName = res.data // Server returns filename
         console.log('Video uploaded:', videoFileName)
       } catch (uploadErr) {
         console.error('Video upload failed', uploadErr)
         videoError.value = '视频上传失败，比分会继续记录：' + uploadErr.message
       }
    }

    // 3. Send to API (No optimistic update)
    const rallyData = {
      match_id: matchInfo.id,
      set_number: 1, // Default to set 1 for now
      rally_number: rallyNumber.value,
      scorer: scorer, // 'player1' or 'player2'
      server: currentServerId.value === matchInfo.player1Id ? 'player1' : 'player2',
      timestamp: new Date().toISOString(),
      video_file: videoFileName
    }

    await api.recordRally(rallyData)

    // 3. Sync state from server
    await fetchMatchState()

  } catch (err) {
    console.error('Error scoring point:', err)
    error.value = '记录得分失败：' + err.message
  } finally {
    loading.value = false
  }
}

const finishMatch = async () => {
  if (loading.value || matchInfo.status === 'finished') return

  if (confirm('您确定要结束比赛吗？')) {
    let finished = false
    try {
      loading.value = true
      await api.finishMatch(matchInfo.id)

      matchInfo.status = 'finished'
      finished = true
      // Update local storage
      localStorage.setItem('currentMatch', JSON.stringify(matchInfo))

    } catch (err) {
      console.error('Failed to finish match:', err)
      error.value = '结束比赛失败：' + apiErrorMessage(err, '请检查服务器是否可用。')
    } finally {
      loading.value = false
      if (finished) {
        stopRecording()
      }
    }
  }
}

const exportExcel = async () => {
  try {
    const response = await api.exportMatch(matchInfo.id)
    // Create blob link to download
    const url = window.URL.createObjectURL(new Blob([response.data]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', `match_${matchInfo.id}.xlsx`)
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
  } catch (err) {
    console.error('Export failed:', err)
    alert('导出失败。请确认服务器是否可用。')
  }
}

const analyzeMatch = async () => {
  const snapshot = `基于请求时比分 ${score.p1}-${score.p2}，第 ${rallyNumber.value} 回合前的数据生成`
  try {
    agentLoading.value = true
    analysisError.value = null
    analysis.value = null
    analysisSnapshot.value = ''
    analysisScore.p1 = null
    analysisScore.p2 = null
    copyStatus.value = ''

    const response = await api.analyzeMatch(matchInfo.id, analysisQuestion.value)
    analysis.value = response.data
    analysisSnapshot.value = snapshot
    analysisScore.p1 = score.p1
    analysisScore.p2 = score.p2
  } catch (err) {
    console.error('Agent analysis failed:', err)
    analysisError.value = formatAgentError(err)
  } finally {
    agentLoading.value = false
  }
}

const generateFinalAnalysis = () => {
  analysisQuestion.value = '比赛已结束，请生成完整赛后复盘，重点说明胜负关键、发球表现、连续得失分和下一场训练建议。'
  analyzeMatch()
}

const copyAnalysis = async () => {
  if (!analysis.value) return

  const text = formatAnalysisText()
  try {
    await navigator.clipboard.writeText(text)
    copyStatus.value = '已复制'
  } catch (err) {
    console.error('Copy analysis failed:', err)
    copyStatus.value = '复制失败'
  }
  window.setTimeout(() => {
    copyStatus.value = ''
  }, 1800)
}

const formatAnalysisText = () => {
  const lines = [
    `比赛复盘：${matchInfo.player1Name || '选手 1'} vs ${matchInfo.player2Name || '选手 2'}`,
    analysisSnapshot.value,
    '',
    `总结：${analysis.value.summary || analysis.value.raw_text || ''}`,
    ...formatSection('优势', analysis.value.strengths),
    ...formatSection('问题', analysis.value.weaknesses),
    ...formatKeyMoments(analysis.value.key_moments),
    ...formatSection('训练建议', analysis.value.training_suggestions)
  ]
  return lines.filter(line => line !== undefined && line !== null).join('\n')
}

const formatSection = (title, items = []) => {
  if (!items.length) return []
  return ['', `${title}：`, ...items.map(item => `- ${item}`)]
}

const formatKeyMoments = (moments = []) => {
  if (!moments.length) return []
  return ['', '关键回合：', ...moments.map(moment => `- #${moment.rally_number}: ${moment.reason}`)]
}

const askSuggestedQuestion = (question) => {
  chatQuestion.value = question
  sendAgentQuestion()
}

const sendAgentQuestion = async () => {
  const question = chatQuestion.value.trim()
  if (!question || chatLoading.value) return
  const snapshot = `基于请求时比分 ${score.p1}-${score.p2}，第 ${rallyNumber.value} 回合前的数据回答`

  const userMessage = {
    id: ++chatMessageId,
    role: 'user',
    content: question
  }
  chatMessages.value.push(userMessage)
  chatQuestion.value = ''

  try {
    chatLoading.value = true
    chatError.value = null

    const response = await api.chatWithAgent(matchInfo.id, question)
    chatMessages.value.push({
      id: ++chatMessageId,
      role: 'agent',
      content: response.data.answer,
      snapshot
    })
  } catch (err) {
    console.error('Agent chat failed:', err)
    const message = formatAgentError(err)
    chatError.value = message
    chatMessages.value.push({
      id: ++chatMessageId,
      role: 'agent',
      content: message
    })
  } finally {
    chatLoading.value = false
  }
}

const formatAgentError = (err) => {
  const status = err.response?.status
  const backendMessage = apiErrorMessage(err, '')
  if (status === 501) return '后端 Agent 当前未启用。'
  if (status === 503) return '后端 Agent 尚未配置模型密钥。'
  if (backendMessage) return backendMessage
  return 'AI 助手请求失败，请检查后端服务和模型配置。'
}

const apiErrorMessage = (err, fallback) => (
  err.response?.data?.error ||
  err.response?.data?.message ||
  networkErrorMessage(err) ||
  fallback
)

const networkErrorMessage = (err) => {
  if (err.code === 'ECONNABORTED') return '请求超时，请检查后端服务是否可用。'
  if (err.message === 'Network Error') return '无法连接后端服务，请检查服务器是否已启动。'
  return ''
}
</script>

<style scoped>
.match-control {
  width: min(1180px, 100%);
  margin: 0 auto;
}

.control-header {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 20px;
  margin-bottom: 24px;
}

.section-kicker {
  margin: 0 0 8px;
  color: var(--court);
  font-size: 13px;
  font-weight: 900;
  letter-spacing: 0;
  text-transform: uppercase;
}

h2 {
  margin: 0;
  font-size: clamp(30px, 5vw, 50px);
  line-height: 1.08;
  letter-spacing: 0;
}

.status-cluster {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  flex-wrap: wrap;
}

.status-pill {
  min-height: 34px;
  display: inline-flex;
  align-items: center;
  padding: 0 12px;
  border: 1px solid var(--line);
  border-radius: 999px;
  color: var(--muted);
  background: rgba(255, 255, 255, 0.78);
  font-size: 13px;
  font-weight: 800;
}

.recording {
  color: #9b2c2c;
  border-color: #f0b9af;
  background: #fff0ed;
  animation: pulse 1.8s ease-in-out infinite;
}

.processing {
  color: #805a18;
  border-color: #f2d49b;
  background: #fff6dd;
}

.message-stack {
  display: grid;
  gap: 10px;
  margin-bottom: 16px;
}

.alert {
  padding: 13px 15px;
  border-radius: 8px;
  line-height: 1.5;
}

.alert-error {
  color: #8f2424;
  background: #fff0ed;
  border: 1px solid #f2c2ba;
}

.alert-warning {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  color: #73511a;
  background: #fff8dc;
  border: 1px solid #f2d49b;
}

.alert-warning button {
  min-height: 32px;
  flex: 0 0 auto;
  padding: 0 10px;
  border: 1px solid #e5bd67;
  border-radius: 8px;
  color: #17201c;
  background: #fff;
  cursor: pointer;
  font-weight: 900;
}

.scoreboard {
  display: grid;
  grid-template-columns: minmax(260px, 1fr) minmax(220px, 0.62fr) minmax(260px, 1fr);
  gap: clamp(14px, 2.5vw, 24px);
  align-items: stretch;
}

.player-score,
.match-center {
  border: 1px solid var(--line);
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.92);
  box-shadow: var(--shadow);
}

.player-score {
  min-height: 440px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  padding: clamp(22px, 4vw, 34px);
  position: relative;
  overflow: hidden;
}

.player-score::before {
  content: "";
  position: absolute;
  inset: 0 0 auto;
  height: 8px;
  background: var(--court);
}

.player-score.accent::before {
  background: var(--amber);
}

.player-label {
  margin: 0 0 10px;
  color: var(--muted);
  font-size: 13px;
  font-weight: 900;
}

h3 {
  min-height: 70px;
  margin: 0;
  font-size: clamp(24px, 4vw, 36px);
  line-height: 1.2;
  letter-spacing: 0;
  word-break: break-word;
}

.score {
  margin: 24px 0;
  color: var(--ink);
  font-family: Georgia, "Times New Roman", serif;
  font-size: clamp(96px, 16vw, 168px);
  font-weight: 900;
  line-height: 0.92;
  text-align: center;
}

.btn-score {
  width: 100%;
  min-height: 58px;
  border: 0;
  border-radius: 8px;
  color: #fff;
  background: linear-gradient(135deg, var(--court), var(--court-dark));
  box-shadow: 0 16px 30px rgba(15, 95, 82, 0.24);
  cursor: pointer;
  font-size: 18px;
  font-weight: 900;
}

.player-score.accent .btn-score {
  color: var(--ink);
  background: linear-gradient(135deg, var(--amber), var(--lime));
}

.match-center {
  display: grid;
  align-content: center;
  gap: 14px;
  padding: 22px;
  background:
    linear-gradient(180deg, rgba(15, 95, 82, 0.08), transparent),
    rgba(255, 255, 255, 0.84);
}

.vs {
  display: grid;
  place-items: center;
  width: 92px;
  height: 92px;
  margin: 0 auto 8px;
  border-radius: 50%;
  color: var(--court-dark);
  background: var(--lime);
  box-shadow: inset -8px -10px 0 rgba(23, 32, 28, 0.09);
  font-size: 24px;
  font-weight: 1000;
}

.info-card {
  display: grid;
  gap: 6px;
  padding: 16px;
  border: 1px solid var(--line);
  border-radius: 8px;
  background: #fff;
}

.info-card span {
  color: var(--muted);
  font-size: 13px;
  font-weight: 800;
}

.info-card strong {
  font-size: 24px;
  line-height: 1.25;
  word-break: break-word;
}

.server-card strong {
  color: var(--court);
}

.actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 22px;
  flex-wrap: wrap;
}

.btn-danger,
.btn-secondary {
  min-height: 46px;
  padding: 0 18px;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 900;
}

.btn-danger {
  border: 1px solid #e3a7a7;
  color: #fff;
  background: var(--danger);
  box-shadow: 0 14px 26px rgba(215, 75, 75, 0.2);
}

.btn-secondary {
  border: 1px solid var(--line);
  color: var(--ink);
  background: #fff;
}

.agent-panel {
  margin-top: 24px;
  padding: clamp(18px, 3vw, 26px);
  border: 1px solid var(--line);
  border-radius: 8px;
  background:
    linear-gradient(135deg, rgba(15, 95, 82, 0.08), transparent 42%),
    rgba(255, 255, 255, 0.9);
  box-shadow: var(--shadow);
}

.agent-header,
.chat-header {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 16px;
}

.agent-header {
  margin-bottom: 18px;
}

.agent-header h3 {
  min-height: 0;
  margin: 0;
  font-size: clamp(24px, 3vw, 34px);
}

.agent-status {
  min-height: 34px;
  display: inline-flex;
  align-items: center;
  padding: 0 12px;
  border: 1px solid rgba(15, 95, 82, 0.22);
  border-radius: 999px;
  color: var(--court-dark);
  background: rgba(207, 232, 95, 0.28);
  font-size: 13px;
  font-weight: 900;
}

.agent-grid {
  display: grid;
  grid-template-columns: minmax(0, 1.2fr) minmax(280px, 0.8fr);
  gap: 18px;
  align-items: stretch;
}

.analysis-box,
.chat-box {
  min-width: 0;
  padding: 18px;
  border: 1px solid var(--line);
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.88);
}

.analysis-box label {
  display: block;
  margin-bottom: 8px;
  color: var(--muted);
  font-size: 13px;
  font-weight: 900;
}

textarea,
.chat-form input {
  width: 100%;
  border: 1px solid var(--line);
  border-radius: 8px;
  color: var(--ink);
  background: var(--panel-soft);
  outline: none;
}

textarea {
  min-height: 96px;
  resize: vertical;
  padding: 12px 14px;
  line-height: 1.5;
}

textarea:focus,
.chat-form input:focus {
  border-color: var(--court);
  box-shadow: 0 0 0 4px rgba(15, 95, 82, 0.12);
}

.btn-agent {
  width: 100%;
  min-height: 48px;
  margin-top: 12px;
  border: 0;
  border-radius: 8px;
  color: #17201c;
  background: linear-gradient(135deg, var(--lime), var(--amber));
  box-shadow: 0 14px 28px rgba(207, 232, 95, 0.22);
  cursor: pointer;
  font-weight: 900;
}

.final-analysis-callout {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 14px;
  margin-bottom: 14px;
  padding: 14px;
  border: 1px solid rgba(207, 232, 95, 0.8);
  border-radius: 8px;
  background: #fbfde9;
}

.final-analysis-callout div {
  display: grid;
  gap: 4px;
}

.final-analysis-callout strong {
  font-size: 16px;
}

.final-analysis-callout span {
  color: var(--muted);
  line-height: 1.5;
}

.final-analysis-callout button,
.copy-review {
  min-height: 34px;
  padding: 0 12px;
  border: 1px solid rgba(15, 95, 82, 0.2);
  border-radius: 8px;
  color: var(--court-dark);
  background: #fff;
  cursor: pointer;
  font-weight: 900;
}

.agent-alert {
  margin: 12px 0 0;
}

.chat-alert {
  margin: 0;
}

.quick-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 10px;
}

.quick-chip {
  min-height: 32px;
  padding: 0 10px;
  border: 1px solid rgba(15, 95, 82, 0.18);
  border-radius: 999px;
  color: var(--court-dark);
  background: rgba(255, 255, 255, 0.84);
  cursor: pointer;
  font-size: 12px;
  font-weight: 900;
}

.quick-chip:hover {
  border-color: rgba(15, 95, 82, 0.36);
  background: rgba(207, 232, 95, 0.22);
}

.analysis-result {
  display: grid;
  gap: 14px;
  margin-top: 16px;
}

.summary-card {
  padding: 16px;
  border-radius: 8px;
  color: #fff;
  background: linear-gradient(135deg, var(--court-dark), var(--court));
}

.summary-card span {
  display: block;
  margin-bottom: 8px;
  color: rgba(255, 255, 255, 0.68);
  font-size: 13px;
  font-weight: 900;
}

.summary-card p,
.chat-message p {
  margin: 0;
  line-height: 1.65;
}

.summary-card small {
  display: block;
  margin-top: 12px;
  color: rgba(255, 255, 255, 0.68);
  line-height: 1.5;
}

.copy-review {
  margin-top: 14px;
  color: #17201c;
  background: var(--lime);
}

.insight-section,
.key-moments {
  padding: 14px;
  border: 1px solid var(--line);
  border-radius: 8px;
  background: #fff;
}

.insight-section h4,
.key-moments h4,
.chat-header h4 {
  margin: 0 0 10px;
  font-size: 16px;
}

.insight-section ul {
  display: grid;
  gap: 8px;
  margin: 0;
  padding-left: 20px;
  line-height: 1.55;
}

.training {
  border-color: rgba(207, 232, 95, 0.8);
  background: #fbfde9;
}

.follow-up-panel {
  padding: 14px;
  border: 1px solid rgba(15, 95, 82, 0.18);
  border-radius: 8px;
  background: rgba(246, 247, 242, 0.86);
}

.follow-up-panel h4 {
  margin: 0;
  font-size: 16px;
}

.moment-row {
  display: grid;
  grid-template-columns: 58px minmax(0, 1fr);
  gap: 10px;
  padding: 10px 0;
  border-top: 1px solid var(--line);
  line-height: 1.5;
}

.moment-row:first-of-type {
  border-top: 0;
}

.moment-row strong {
  color: var(--court);
}

.empty-agent-state,
.chat-empty {
  padding: 16px;
  border: 1px dashed #cdd8c6;
  border-radius: 8px;
  color: var(--muted);
  background: rgba(246, 247, 242, 0.72);
  line-height: 1.6;
}

.empty-agent-state {
  margin-top: 14px;
}

.chat-box {
  display: grid;
  grid-template-rows: auto minmax(260px, 1fr) auto;
  gap: 14px;
}

.chat-header span {
  color: var(--muted);
  font-size: 13px;
  font-weight: 800;
}

.chat-log {
  min-height: 260px;
  max-height: 520px;
  overflow: auto;
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding-right: 4px;
}

.chat-message {
  display: grid;
  gap: 6px;
  width: min(92%, 460px);
}

.chat-message span {
  color: var(--muted);
  font-size: 12px;
  font-weight: 900;
}

.chat-message p {
  padding: 12px 14px;
  border-radius: 8px;
  background: var(--panel-soft);
}

.chat-message small {
  color: var(--muted);
  font-size: 12px;
  line-height: 1.45;
}

.chat-message.user {
  align-self: flex-end;
}

.chat-message.user span {
  text-align: right;
}

.chat-message.user p {
  color: #fff;
  background: var(--court);
}

.chat-form {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  gap: 10px;
}

.chat-form input {
  min-height: 46px;
  padding: 0 14px;
}

@keyframes pulse {
  0%,
  100% {
    opacity: 1;
  }

  50% {
    opacity: 0.58;
  }
}

@media (max-width: 960px) {
  .scoreboard {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .agent-grid {
    grid-template-columns: 1fr;
  }

  .player-score {
    min-height: 300px;
  }

  .match-center {
    grid-column: 1 / -1;
    grid-template-columns: auto minmax(0, 1fr) minmax(0, 1fr);
    align-items: center;
    order: -1;
  }

  .vs {
    width: 72px;
    height: 72px;
    margin: 0;
  }

  .info-card {
    min-height: 78px;
  }

  .score {
    font-size: clamp(82px, 18vw, 128px);
  }
}

@media (max-width: 680px) {
  .control-header {
    align-items: flex-start;
    flex-direction: column;
  }

  .status-cluster,
  .actions {
    justify-content: flex-start;
  }

  .player-score {
    min-height: 320px;
  }

  .scoreboard {
    grid-template-columns: 1fr;
  }

  .match-center {
    grid-template-columns: 1fr;
  }

  .vs {
    margin: 0 auto 8px;
  }

  .agent-header {
    align-items: flex-start;
    flex-direction: column;
  }

  .chat-form {
    grid-template-columns: 1fr;
  }

  .chat-form .btn-secondary {
    width: 100%;
  }

  .final-analysis-callout {
    align-items: stretch;
    flex-direction: column;
  }
}
</style>

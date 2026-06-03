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
        <button @click="scorePoint('player1')" :disabled="loading || matchInfo.status === 'finished'" class="btn-score">
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
        <button @click="scorePoint('player2')" :disabled="loading || matchInfo.status === 'finished'" class="btn-score">
          记为得分
        </button>
      </div>
    </div>

    <div class="actions">
      <button @click="finishMatch" :disabled="matchInfo.status === 'finished'" class="btn-danger">
        结束比赛
      </button>
      <button @click="exportExcel" :disabled="matchInfo.status !== 'finished'" class="btn-secondary">
        导出 Excel
      </button>
    </div>
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

const fetchMatchState = async () => {
  try {
    const res = await api.getMatchCurrent(matchInfo.id)
    const data = res.data
    rallyNumber.value = data.rally_number + 1 || 1
    score.p1 = data.score_p1 || 0
    score.p2 = data.score_p2 || 0

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
  // Load match info from local storage (mocking backend retrieval)
  const storedMatch = localStorage.getItem('currentMatch')
  if (storedMatch) {
    const parsed = JSON.parse(storedMatch)
    Object.assign(matchInfo, parsed)
    currentServerId.value = parsed.firstServerId

    await fetchMatchState()
  } else {
    // Redirect back if no match info
    router.push('/')
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

    // 2. Restart recording immediately for the NEXT rally
    // We don't await this because we want to proceed with upload,
    // but we ensure it's triggered.
    startRecording()

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
         // Continue without video?
         videoError.value = '视频上传失败：' + uploadErr.message
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
  if (confirm('您确定要结束比赛吗？')) {
    try {
      loading.value = true
      await api.finishMatch(matchInfo.id)

      matchInfo.status = 'finished'
      // Update local storage
      localStorage.setItem('currentMatch', JSON.stringify(matchInfo))

    } catch (err) {
      console.error('Failed to finish match:', err)
      error.value = '结束比赛失败'
      // Fallback for demo
      matchInfo.status = 'finished'
      localStorage.setItem('currentMatch', JSON.stringify(matchInfo))
    } finally {
      loading.value = false
      stopRecording()
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
  } catch (err) {
    console.error('Export failed:', err)
    alert('导出失败。请确认服务器是否可用。')
  }
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
}
</style>

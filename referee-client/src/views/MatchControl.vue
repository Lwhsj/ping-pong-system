<template>
  <div class="match-control">
    <div class="header">
      <h2>比赛</h2>
      <div class="match-status">状态: {{ matchInfo.status }}</div>
      <div v-if="error" class="error-message">API 错误: {{ error }}</div>
      <div v-if="videoError" class="error-message">视频错误: {{ videoError }}</div>
      <div v-if="isRecording" class="recording-indicator">🔴 录制中</div>
      <div v-if="loading" class="loading-indicator">处理中...</div>
    </div>

    <div class="scoreboard">
      <div class="player-score">
        <h3>{{ matchInfo.player1Name }}</h3>
        <div class="score">{{ score.p1 }}</div>
        <button @click="scorePoint('player1')" :disabled="matchInfo.status === 'finished'" class="btn-score">
          得分
        </button>
      </div>

      <div class="vs">VS</div>

      <div class="player-score">
        <h3>{{ matchInfo.player2Name }}</h3>
        <div class="score">{{ score.p2 }}</div>
        <button @click="scorePoint('player2')" :disabled="matchInfo.status === 'finished'" class="btn-score">
          得分
        </button>
      </div>
    </div>

    <div class="info-panel">
      <p><strong>回合:</strong> {{ rallyNumber }}</p>
      <p><strong>当前发球方:</strong> {{ currentServerName }}</p>
    </div>

    <div class="actions">
      <button @click="finishMatch" :disabled="matchInfo.status === 'finished'" class="btn-danger">
        结束比赛
      </button>
      <button @click="exportExcel" :disabled="matchInfo.status !== 'finished'" class="btn-secondary">
        导出 Excel
      </button>
    </div>
  </div>
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
         videoError.value = '视频上传失败: ' + uploadErr.message
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
    error.value = 'Failed to record score: ' + err.message
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
      error.value = 'Failed to finish match'
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
    alert('Export failed. Server might be unavailable.')
  }
}
</script>

<style scoped>
.match-control {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
  text-align: center;
}

.header {
  margin-bottom: 30px;
  border-bottom: 1px solid #eee;
  padding-bottom: 10px;
}

.scoreboard {
  display: flex;
  justify-content: space-around;
  align-items: center;
  margin-bottom: 30px;
  background: #f9f9f9;
  padding: 20px;
  border-radius: 8px;
}

.player-score {
  text-align: center;
}

.score {
  font-size: 48px;
  font-weight: bold;
  color: #333;
  margin: 10px 0;
}

.vs {
  font-size: 24px;
  font-weight: bold;
  color: #999;
}

.btn-score {
  background-color: #28a745;
  color: white;
  padding: 10px 20px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 18px;
}

.btn-score:disabled {
  background-color: #ccc;
  cursor: not-allowed;
}

.info-panel {
  margin-bottom: 30px;
  font-size: 18px;
}

.actions {
  display: flex;
  justify-content: center;
  gap: 20px;
}

.btn-danger {
  background-color: #dc3545;
  color: white;
  padding: 10px 20px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.btn-secondary {
  background-color: #6c757d;
  color: white;
  padding: 10px 20px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.btn-danger:disabled, .btn-secondary:disabled {
  background-color: #ccc;
  cursor: not-allowed;
}

.recording-indicator {
  color: red;
  font-weight: bold;
  animation: blink 2s infinite;
  margin-left: 10px;
  display: inline-block;
}

.error-message {
  color: red;
  margin-top: 5px;
}

@keyframes blink {
  0% { opacity: 1; }
  50% { opacity: 0.5; }
  100% { opacity: 1; }
}
</style>

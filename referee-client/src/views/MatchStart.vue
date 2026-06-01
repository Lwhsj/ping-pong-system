<template>
  <div class="match-start">
    <h2>开始新比赛</h2>
    
    <div v-if="error" class="error-alert">{{ error }}</div>
    <div v-if="loading" class="loading">加载中...</div>

    <form @submit.prevent="startMatch" v-else>
      <div class="form-group">
        <label for="player1">选手 1:</label>
        <select id="player1" v-model="matchData.player1Id" required>
          <option value="" disabled>请选择选手 1</option>
          <option v-for="player in players" :key="player.id" :value="player.id">
            {{ player.name }}
          </option>
        </select>
      </div>

      <div class="form-group">
        <label for="player2">选手 2:</label>
        <select id="player2" v-model="matchData.player2Id" required>
          <option value="" disabled>请选择选手 2</option>
          <option v-for="player in players" :key="player.id" :value="player.id">
            {{ player.name }}
          </option>
        </select>
      </div>

      <div class="form-group">
        <label for="server">先发球方:</label>
        <select id="server" v-model="matchData.firstServerId" required>
          <option value="" disabled>请选择先发球方</option>
          <option v-if="matchData.player1Id" :value="matchData.player1Id">
            {{ getPlayerName(matchData.player1Id) }}
          </option>
          <option v-if="matchData.player2Id" :value="matchData.player2Id">
            {{ getPlayerName(matchData.player2Id) }}
          </option>
        </select>
      </div>

      <div class="form-group">
        <label>比赛日期:</label>
        <input type="date" v-model="matchData.date" disabled />
      </div>

      <button type="submit" class="btn-primary">开始比赛</button>
    </form>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import api from '../services/api'

const router = useRouter()

const players = ref([])
const loading = ref(false)
const error = ref(null)

const matchData = reactive({
  player1Id: '',
  player2Id: '',
  firstServerId: '',
  date: ''
})

const getPlayerName = (id) => {
  const player = players.value.find(p => p.id === id)
  return player ? player.name : ''
}

onMounted(async () => {
  // Set date to local client date formatted as yyyy-mm-dd
  const now = new Date()
  const year = now.getFullYear()
  const month = String(now.getMonth() + 1).padStart(2, '0')
  const day = String(now.getDate()).padStart(2, '0')
  matchData.date = `${year}-${month}-${day}`

  try {
    loading.value = true
    const response = await api.getPlayers()
    players.value = response.data
  } catch (err) {
    console.error('Failed to fetch players:', err)
    error.value = '加载选手列表失败。请检查服务器连接。'
    // No mock fallback as requested
  } finally {
    loading.value = false
  }
})

const startMatch = async () => {
  try {
    loading.value = true
    error.value = null
    
    const payload = {
      player1_id: matchData.player1Id,
      player2_id: matchData.player2Id,
      first_server: matchData.firstServerId
    }

    const response = await api.startMatch(payload)
    const match = response.data

    // Store match info in localStorage
    localStorage.setItem('currentMatch', JSON.stringify({
      id: match.id,
      ...matchData,
      ...match,
      player1Name: getPlayerName(matchData.player1Id),
      player2Name: getPlayerName(matchData.player2Id),
    }))

    router.push({ name: 'MatchControl', params: { id: match.id } })
  } catch (err) {
    console.error('Failed to start match:', err)
    error.value = '开始比赛失败: ' + (err.response?.data?.message || err.message)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.match-start {
  max-width: 500px;
  margin: 0 auto;
  padding: 20px;
}

.form-group {
  margin-bottom: 15px;
}

label {
  display: block;
  margin-bottom: 5px;
  font-weight: bold;
}

select, input {
  width: 100%;
  padding: 8px;
  border: 1px solid #ccc;
  border-radius: 4px;
  box-sizing: border-box; /* Ensure padding doesn't affect width */
}

.btn-primary {
  background-color: #007bff;
  color: white;
  padding: 10px 15px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  width: 100%;
  font-size: 16px;
}

.btn-primary:hover {
  background-color: #0056b3;
}

.error-alert {
  background-color: #f8d7da;
  color: #721c24;
  padding: 10px;
  margin-bottom: 15px;
  border-radius: 4px;
}

.loading {
  text-align: center;
  margin: 20px;
  font-style: italic;
  color: #666;
}
</style>

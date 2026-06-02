<template>
  <section class="match-start">
    <div class="start-copy">
      <p class="section-kicker">Match Setup</p>
      <h2>开始新比赛</h2>
      <p class="section-intro">
        选择双方选手和先发球方，系统会自动创建比赛并进入裁判控制台。
      </p>

      <div class="setup-summary" aria-label="比赛配置摘要">
        <div>
          <span>选手 1</span>
          <strong>{{ getPlayerName(matchData.player1Id) || '待选择' }}</strong>
        </div>
        <div>
          <span>选手 2</span>
          <strong>{{ getPlayerName(matchData.player2Id) || '待选择' }}</strong>
        </div>
        <div>
          <span>先发球</span>
          <strong>{{ getPlayerName(matchData.firstServerId) || '待选择' }}</strong>
        </div>
      </div>
    </div>

    <div class="setup-panel">
      <div v-if="error" class="alert alert-error">{{ error }}</div>
      <div v-if="loading" class="loading-state">
        <span class="spinner" aria-hidden="true"></span>
        正在加载...
      </div>

      <form v-else @submit.prevent="startMatch">
        <div class="form-grid">
          <div class="form-group">
            <label for="player1">选手 1</label>
            <select id="player1" v-model="matchData.player1Id" required>
              <option value="" disabled>请选择选手 1</option>
              <option v-for="player in players" :key="player.id" :value="player.id">
                {{ player.name }}
              </option>
            </select>
          </div>

          <div class="form-group">
            <label for="player2">选手 2</label>
            <select id="player2" v-model="matchData.player2Id" required>
              <option value="" disabled>请选择选手 2</option>
              <option v-for="player in players" :key="player.id" :value="player.id">
                {{ player.name }}
              </option>
            </select>
          </div>

          <div class="form-group form-group-wide">
            <label for="server">先发球方</label>
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

          <div class="form-group form-group-wide">
            <label for="match-date">比赛日期</label>
            <input id="match-date" type="date" v-model="matchData.date" disabled />
          </div>
        </div>

        <button type="submit" class="btn-primary">
          开始比赛
        </button>
      </form>
    </div>
  </section>
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
    error.value = '开始比赛失败：' + (err.response?.data?.message || err.message)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.match-start {
  width: min(1120px, 100%);
  margin: 0 auto;
  display: grid;
  grid-template-columns: minmax(280px, 0.88fr) minmax(320px, 1.12fr);
  gap: clamp(22px, 4vw, 44px);
  align-items: stretch;
}

.start-copy {
  min-height: 560px;
  padding: clamp(28px, 5vw, 54px);
  border-radius: 8px;
  color: #fff;
  background:
    linear-gradient(160deg, rgba(207, 232, 95, 0.24), transparent 36%),
    repeating-linear-gradient(90deg, rgba(255, 255, 255, 0.05) 0 1px, transparent 1px 72px),
    linear-gradient(135deg, #123f3b, #0f5f52);
  box-shadow: var(--shadow);
  position: relative;
  overflow: hidden;
}

.start-copy::after {
  content: "";
  position: absolute;
  right: -80px;
  bottom: -80px;
  width: 260px;
  height: 260px;
  border: 32px solid rgba(207, 232, 95, 0.18);
  border-radius: 50%;
}

.section-kicker {
  margin: 0 0 12px;
  color: var(--lime);
  font-size: 13px;
  font-weight: 800;
  letter-spacing: 0;
  text-transform: uppercase;
}

h2 {
  margin: 0;
  font-size: clamp(34px, 5vw, 58px);
  line-height: 1.04;
  letter-spacing: 0;
}

.section-intro {
  max-width: 420px;
  margin: 20px 0 0;
  color: rgba(255, 255, 255, 0.78);
  font-size: 17px;
  line-height: 1.8;
}

.setup-summary {
  position: relative;
  z-index: 1;
  display: grid;
  gap: 12px;
  margin-top: clamp(38px, 9vw, 120px);
}

.setup-summary div {
  display: flex;
  justify-content: space-between;
  gap: 18px;
  padding: 16px 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.22);
}

.setup-summary span {
  color: rgba(255, 255, 255, 0.68);
}

.setup-summary strong {
  text-align: right;
  font-size: 18px;
}

.setup-panel {
  align-self: center;
  padding: clamp(22px, 4vw, 36px);
  border: 1px solid var(--line);
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.9);
  box-shadow: var(--shadow);
  backdrop-filter: blur(18px);
}

.form-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 18px;
}

.form-group {
  display: grid;
  gap: 8px;
}

.form-group-wide {
  grid-column: 1 / -1;
}

label {
  color: var(--muted);
  font-size: 13px;
  font-weight: 800;
}

select,
input {
  width: 100%;
  min-height: 50px;
  padding: 0 14px;
  border: 1px solid var(--line);
  border-radius: 8px;
  color: var(--ink);
  background: var(--panel-soft);
  outline: none;
}

select:focus,
input:focus {
  border-color: var(--court);
  box-shadow: 0 0 0 4px rgba(15, 95, 82, 0.12);
}

input:disabled {
  color: var(--muted);
  background: #edf0e8;
}

.btn-primary {
  width: 100%;
  min-height: 54px;
  margin-top: 24px;
  border: 0;
  border-radius: 8px;
  color: #17201c;
  background: linear-gradient(135deg, var(--lime), var(--amber));
  box-shadow: 0 16px 34px rgba(207, 232, 95, 0.28);
  cursor: pointer;
  font-size: 17px;
  font-weight: 900;
}

.alert {
  padding: 14px 16px;
  border-radius: 8px;
  margin-bottom: 18px;
  line-height: 1.5;
}

.alert-error {
  color: #8f2424;
  background: #fff0ed;
  border: 1px solid #f2c2ba;
}

.loading-state {
  min-height: 180px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  color: var(--muted);
  font-weight: 700;
}

.spinner {
  width: 18px;
  height: 18px;
  border: 3px solid #d9e1d3;
  border-top-color: var(--court);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

@media (max-width: 860px) {
  .match-start {
    grid-template-columns: 1fr;
  }

  .start-copy {
    min-height: auto;
  }
}

@media (max-width: 560px) {
  .form-grid {
    grid-template-columns: 1fr;
  }

  .setup-panel {
    padding: 18px;
  }
}
</style>

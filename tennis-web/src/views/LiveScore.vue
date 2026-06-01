<template>
  <div class="live-score-container">
    <div v-if="liveMatches.length > 0">
      <div class="live-hero">
        <div>
          <p class="section-kicker">Live Scoreboard</p>
          <h1>实时比分</h1>
          <p>每 3 秒刷新当前比赛分数与回合记录，方便现场快速观察局势变化。</p>
        </div>
        <div class="live-pill">
          <span class="live-dot"></span>
          {{ liveMatches.length }} 场直播中
        </div>
      </div>

      <el-tabs v-model="activeMatchId" type="card" @tab-change="handleTabChange" class="match-tabs">
        <el-tab-pane
          v-for="item in liveMatches"
          :key="item.id"
          :label="`${item.player1_name} VS ${item.player2_name}`"
          :name="item.id"
        />
      </el-tabs>

      <el-card v-if="match" class="scoreboard-card" shadow="never">
        <template #header>
          <div class="card-header">
            <div>
              <span class="card-title">实时比赛记分板</span>
              <span class="card-subtitle">Match ID {{ match.match_id }}</span>
            </div>
            <el-tag type="danger" effect="dark" round>直播中</el-tag>
          </div>
        </template>

        <div class="score-display">
          <div class="score-glow"></div>
          <div class="player-col">
            <div class="player-card player-one">
              <div class="player-name">{{ match.player1_name || 'Player 1' }}</div>
              <div class="player-score">{{ match.score_p1 }}</div>
              <div v-if="match.server === 'player1' || match.server === 'Player 1' || match.server === match.player1_name" class="server-indicator">
                <el-icon><Basketball /></el-icon> 发球
              </div>
            </div>
          </div>

          <div class="vs-col">
            <div class="vs-text">VS</div>
            <div class="rally-chip">第 {{ match.rally_number + 1 }} 回合</div>
          </div>

          <div class="player-col">
            <div class="player-card player-two">
              <div class="player-name">{{ match.player2_name || 'Player 2' }}</div>
              <div class="player-score">{{ match.score_p2 }}</div>
              <div v-if="match.server === 'player2' || match.server === 'Player 2' || match.server === match.player2_name" class="server-indicator">
                <el-icon><Basketball /></el-icon> 发球
              </div>
            </div>
          </div>
        </div>

        <el-descriptions class="info-strip" border :column="3">
          <el-descriptions-item label="比赛 ID">{{ match.match_id }}</el-descriptions-item>
          <el-descriptions-item label="当前回合">{{ match.rally_number + 1 }}</el-descriptions-item>
          <el-descriptions-item label="发球方">{{ match.server }}</el-descriptions-item>
        </el-descriptions>
      </el-card>

      <el-card v-if="match" class="rally-card" shadow="never">
        <template #header>
          <div class="card-header">
            <div>
              <span class="card-title">回合记录</span>
              <span class="card-subtitle">最新回合优先展示</span>
            </div>
          </div>
        </template>
        <el-table :data="rallies" class="rally-table" style="width: 100%" max-height="500" stripe size="default" :default-sort="{ prop: 'rally_number', order: 'descending' }">
          <el-table-column prop="rally_number" label="回合" width="80" align="center" sortable />
          <el-table-column label="比分 (P1 - P2)" align="center" width="120">
            <template #default="scope">
              <span class="score-line">{{ scope.row.score_p1 }} - {{ scope.row.score_p2 }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="scorer" label="得分方" align="center">
            <template #default="scope">
              <el-tag :type="scope.row.scorer === 'player1' ? 'primary' : 'warning'" round>
                {{ scope.row.scorer === 'player1' ? (match.player1_name || 'Player 1') : (match.player2_name || 'Player 2') }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="server" label="发球方" align="center">
             <template #default="scope">
                {{ scope.row.server === 'player1' ? (match.player1_name || 'Player 1') : (scope.row.server === 'player2' ? (match.player2_name || 'Player 2') : scope.row.server) }}
             </template>
          </el-table-column>
          <el-table-column prop="timestamp" label="时间" align="center" min-width="160">
             <template #default="scope">
               {{ formatTime(scope.row.timestamp) }}
             </template>
          </el-table-column>
        </el-table>
      </el-card>
    </div>

    <el-empty v-else class="empty-state" description="当前没有正在进行的比赛" />
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getCurrentScore, getMatchList, getMatchDetail } from '@/api/match'

const route = useRoute()
const router = useRouter()
const match = ref(null)
const liveMatches = ref([])
const activeMatchId = ref(null)
const rallies = ref([])
let timer = null

const formatTime = (timeStr) => {
  if (!timeStr) return ''
  return new Date(timeStr).toLocaleString()
}

const fetchLiveMatches = async () => {
  try {
    const data = await getMatchList({ status: 'started' })
    liveMatches.value = data || []
    
    // Determine active match
    if (liveMatches.value.length > 0) {
      const queryId = Number(route.query.matchId)
      const exists = liveMatches.value.find(m => m.id === queryId)
      
      if (exists) {
        activeMatchId.value = queryId
      } else if (!activeMatchId.value) {
        // Default to the first one if no active match selected
        activeMatchId.value = liveMatches.value[0].id
      }
    } else {
        match.value = null
        activeMatchId.value = null
        rallies.value = []
    }
  } catch (error) {
    console.error('Failed to fetch live matches:', error)
  }
}

const fetchScore = async () => {
  if (!activeMatchId.value) return
  
  try {
    const [scoreData, rallyData] = await Promise.all([
        getCurrentScore(activeMatchId.value),
        getMatchDetail(activeMatchId.value)
    ])
    match.value = scoreData
    rallies.value = rallyData
  } catch (error) {
    console.error('Failed to fetch score or details:', error)
  }
}

const handleTabChange = (val) => {
  activeMatchId.value = val
  fetchScore()
  // Update URL query param without reloading to keep state shareable
  router.replace({ query: { ...route.query, matchId: val } })
}

onMounted(async () => {
  await fetchLiveMatches()
  if (activeMatchId.value) {
    fetchScore()
    // Poll for score updates
    timer = setInterval(fetchScore, 3000)
  }
  
  // Periodically refresh the list of live matches as well, 
  // in case a new match starts or current one finishes
  // But maybe 3s is too frequent for match list. 
  // Let's stick to score polling for now.
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})

// Watch activeMatchId to switch data source immediately
watch(activeMatchId, (newVal) => {
    if (newVal) {
        fetchScore()
    }
})
</script>

<style scoped>
.live-score-container {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.live-hero {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
  padding: 26px 28px;
  color: #fff;
  overflow: hidden;
  background:
    radial-gradient(circle at 86% 12%, rgba(255, 193, 66, 0.3), transparent 30%),
    linear-gradient(135deg, #141f32 0%, #163c42 100%);
  border-radius: 18px;
  box-shadow: 0 20px 50px rgba(26, 39, 58, 0.18);
}

.section-kicker {
  margin: 0 0 8px;
  color: #84e8df;
  font-size: 12px;
  font-weight: 850;
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

.live-hero h1 {
  margin: 0;
  font-size: 32px;
  font-weight: 950;
  letter-spacing: 0;
}

.live-hero p {
  max-width: 560px;
  margin: 10px 0 0;
  color: #c9d8e8;
  font-size: 14px;
}

.live-pill {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  flex-shrink: 0;
  padding: 12px 16px;
  color: #fff6cf;
  font-size: 13px;
  font-weight: 850;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.16);
  border-radius: 999px;
}

.live-dot {
  width: 9px;
  height: 9px;
  background: #ff4d62;
  border-radius: 999px;
  box-shadow: 0 0 0 7px rgba(255, 77, 98, 0.15);
}

.match-tabs :deep(.el-tabs__header) {
  margin-bottom: 0;
}

.match-tabs :deep(.el-tabs__nav) {
  overflow: hidden;
  border: 1px solid #d8e3ee;
  border-radius: 14px;
}

.match-tabs :deep(.el-tabs__item) {
  height: 44px;
  color: #647287;
  font-weight: 800;
  background: rgba(255, 255, 255, 0.8);
  border-left-color: #d8e3ee;
}

.match-tabs :deep(.el-tabs__item.is-active) {
  color: #142033;
  background: #fff6d6;
}

.scoreboard-card,
.rally-card {
  width: 100%;
  overflow: hidden;
  background: rgba(255, 255, 255, 0.94);
  border: 1px solid rgba(211, 221, 232, 0.9);
  border-radius: 18px;
  box-shadow: 0 18px 46px rgba(28, 42, 61, 0.09);
}

.scoreboard-card :deep(.el-card__header),
.rally-card :deep(.el-card__header) {
  padding: 18px 22px;
  border-bottom-color: #edf2f6;
}

.scoreboard-card :deep(.el-card__body),
.rally-card :deep(.el-card__body) {
  padding: 22px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
}

.card-title {
  display: block;
  color: #182539;
  font-size: 18px;
  font-weight: 900;
}

.card-subtitle {
  display: block;
  margin-top: 3px;
  color: #7a8797;
  font-size: 12px;
  font-weight: 700;
}

.score-display {
  position: relative;
  display: grid;
  grid-template-columns: minmax(0, 1fr) 150px minmax(0, 1fr);
  align-items: center;
  gap: 18px;
  margin-bottom: 20px;
  padding: 28px;
  overflow: hidden;
  background:
    radial-gradient(circle at 50% 0%, rgba(51, 218, 203, 0.18), transparent 28%),
    linear-gradient(135deg, rgba(20, 32, 51, 0.97), rgba(20, 61, 65, 0.95));
  border-radius: 18px;
}

.score-glow {
  position: absolute;
  inset: 18px 42%;
  background: rgba(51, 218, 203, 0.18);
  filter: blur(42px);
}

.player-card {
  position: relative;
  z-index: 1;
  min-height: 220px;
  padding: 26px 20px;
  text-align: center;
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.13);
  border-radius: 16px;
}

.player-card::before {
  position: absolute;
  top: 16px;
  left: 18px;
  width: 42px;
  height: 4px;
  content: '';
  background: #33dacb;
  border-radius: 999px;
}

.player-two::before {
  right: 18px;
  left: auto;
  background: #ffbd40;
}

.player-name {
  min-height: 34px;
  color: #d8e6f2;
  font-size: 22px;
  font-weight: 850;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.player-score {
  margin-top: 14px;
  color: #ffffff;
  font-size: 86px;
  font-weight: 950;
  line-height: 1;
  text-shadow: 0 18px 32px rgba(0, 0, 0, 0.26);
}

.vs-col {
  position: relative;
  z-index: 1;
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
}

.vs-text {
  display: grid;
  width: 88px;
  height: 88px;
  place-items: center;
  color: #142033;
  font-size: 26px;
  font-weight: 950;
  font-style: italic;
  background: linear-gradient(135deg, #fff7cf, #33dacb);
  border-radius: 999px;
  box-shadow: 0 18px 36px rgba(0, 0, 0, 0.2);
}

.rally-chip {
  padding: 7px 12px;
  color: #d8e6f2;
  font-size: 13px;
  font-weight: 850;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.13);
  border-radius: 999px;
}

.server-indicator {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  margin-top: 18px;
  padding: 7px 12px;
  color: #142033;
  font-weight: 900;
  background: #fff2b8;
  border-radius: 999px;
}

.info-strip {
  overflow: hidden;
  border-radius: 14px;
}

.info-strip :deep(.el-descriptions__label) {
  color: #65758a;
  font-weight: 850;
  background: #f7fafc;
}

.info-strip :deep(.el-descriptions__content) {
  color: #172338;
  font-weight: 850;
}

.rally-table :deep(.el-table__header th) {
  color: #5f6f82;
  background: #f7fafc;
  font-weight: 850;
}

.score-line {
  color: #142033;
  font-weight: 950;
}

.empty-state {
  min-height: 420px;
  background: rgba(255, 255, 255, 0.88);
  border: 1px solid rgba(211, 221, 232, 0.86);
  border-radius: 18px;
}

@media (max-width: 900px) {
  .live-hero {
    flex-direction: column;
    align-items: flex-start;
  }

  .score-display {
    grid-template-columns: 1fr;
  }

  .player-card {
    min-height: 180px;
  }

  .player-score {
    font-size: 72px;
  }
}

@media (max-width: 620px) {
  .scoreboard-card :deep(.el-card__body),
  .rally-card :deep(.el-card__body) {
    padding: 16px;
  }

  .score-display {
    padding: 18px;
  }

  .player-name {
    font-size: 18px;
  }
}
</style>

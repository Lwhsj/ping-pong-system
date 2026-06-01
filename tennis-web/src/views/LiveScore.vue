<template>
  <div class="live-score-container">
    <div v-if="liveMatches.length > 0">
      <el-tabs v-model="activeMatchId" type="card" @tab-change="handleTabChange" class="match-tabs">
        <el-tab-pane
          v-for="item in liveMatches"
          :key="item.id"
          :label="`${item.player1_name} VS ${item.player2_name}`"
          :name="item.id"
        />
      </el-tabs>

      <el-card v-if="match" class="scoreboard-card">
        <template #header>
          <div class="card-header">
            <span>实时比赛记分板</span>
            <el-tag type="danger" effect="dark">直播中</el-tag>
          </div>
        </template>
        
        <div class="score-display">
          <el-row :gutter="20" align="middle" justify="center">
            <!-- Player 1 -->
            <el-col :span="8" class="player-col">
              <div class="player-name">{{ match.player1_name || 'Player 1' }}</div>
              <div class="player-score">{{ match.score_p1 }}</div>
              <div v-if="match.server === 'player1' || match.server === 'Player 1' || match.server === match.player1_name" class="server-indicator">
                <el-icon><Basketball /></el-icon> 发球
              </div>
            </el-col>
            
            <!-- VS -->
            <el-col :span="4" class="vs-col">
              <div class="vs-text">VS</div>
            </el-col>
            
            <!-- Player 2 -->
            <el-col :span="8" class="player-col">
              <div class="player-name">{{ match.player2_name || 'Player 2' }}</div>
              <div class="player-score">{{ match.score_p2 }}</div>
              <div v-if="match.server === 'player2' || match.server === 'Player 2' || match.server === match.player2_name" class="server-indicator">
                <el-icon><Basketball /></el-icon> 发球
              </div>
            </el-col>
          </el-row>
        </div>

        <el-divider content-position="center">比赛信息</el-divider>

        <el-descriptions border :column="3">
          <el-descriptions-item label="比赛 ID">{{ match.match_id }}</el-descriptions-item>
          <el-descriptions-item label="当前回合">{{ match.rally_number + 1 }}</el-descriptions-item>
          <el-descriptions-item label="发球方">{{ match.server }}</el-descriptions-item>
        </el-descriptions>
      </el-card>

      <!-- Rally Details Table -->
      <el-card v-if="match" class="box-card" style="margin-top: 20px;">
        <template #header>
          <div class="card-header">
            <span>回合记录</span>
          </div>
        </template>
        <el-table :data="rallies" style="width: 100%" max-height="500" stripe size="default" :default-sort="{ prop: 'rally_number', order: 'descending' }">
          <el-table-column prop="rally_number" label="回合" width="80" align="center" sortable />
          <el-table-column label="比分 (P1 - P2)" align="center" width="120">
            <template #default="scope">
              {{ scope.row.score_p1 }} - {{ scope.row.score_p2 }}
            </template>
          </el-table-column>
          <el-table-column prop="scorer" label="得分方" align="center">
            <template #default="scope">
              <el-tag :type="scope.row.scorer === 'player1' ? 'primary' : 'warning'">
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

    <el-empty v-else description="当前没有正在进行的比赛" />
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
  /* padding-top: 20px; */
}

.match-tabs {
  margin-bottom: 20px;
}

.scoreboard-card {
  width: 100%;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 18px;
  font-weight: bold;
}

.score-display {
  padding: 40px 0;
  background: linear-gradient(to bottom, #f5f7fa, #e4e7ed);
  border-radius: 8px;
  margin-bottom: 20px;
}

.player-col {
  text-align: center;
}

.player-name {
  font-size: 24px;
  font-weight: bold;
  margin-bottom: 10px;
  color: #303133;
}

.player-score {
  font-size: 64px;
  font-weight: 900;
  color: #409EFF;
  line-height: 1;
}

.vs-col {
  text-align: center;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.vs-text {
  font-size: 32px;
  font-weight: bold;
  color: #909399;
  font-style: italic;
}

.rally-info {
  margin-top: 10px;
  font-size: 16px;
  color: #606266;
  background-color: #fff;
  padding: 4px 12px;
  border-radius: 12px;
  display: inline-block;
}

.server-indicator {
  margin-top: 10px;
  color: #E6A23C;
  font-weight: bold;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 5px;
}
</style>
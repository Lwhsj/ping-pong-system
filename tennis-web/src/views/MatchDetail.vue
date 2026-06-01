<template>
  <div class="match-detail-container" style="height: 100%;">
    <div v-if="matchId" style="height: 100%;">
      <el-row :gutter="20" style="height: 100%;">
        <!-- Stats Panel -->
        <el-col :span="10">
          <el-card class="box-card" style="margin-bottom: 20px; height: 100%;">
            <template #header>
              <div class="card-header">
                <span>数据统计</span>
              </div>
            </template>
            <div id="chart-container" style="width: 100%; height: 400px;"></div>
            <div class="stats-text" style="margin-top: 20px;">
              <el-table :data="statsData" border style="width: 100%">
                <el-table-column prop="name" label="选手" align="center" />
                <el-table-column prop="consecutive" label="最大连胜分数" align="center" />
                <el-table-column label="平均回合时间" align="center">
                  <template #default>
                    {{ stats.average_rally_time }}s
                  </template>
                </el-table-column>
              </el-table>
            </div>
          </el-card>
        </el-col>

        <!-- Rounds Panel -->
        <el-col :span="14" style="height: 100%;">
          <el-card class="box-card" style="height: 100%;">
            <template #header>
              <div class="card-header">
                <span>回合记录</span>
              </div>
            </template>
            <el-table :data="rounds" style="width: 100%" height="450" stripe size="large">
              <el-table-column prop="rally_number" label="回合" width="80" align="center" />
              <el-table-column prop="scorer" label="得分方" align="center">
                <template #default="scope">
                  <el-tag :type="scope.row.scorer === 'player1' ? 'primary' : 'warning'">
                    {{ getPlayerName(scope.row.scorer) }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="server" label="发球方" align="center">
                <template #default="scope">
                  {{ getPlayerName(scope.row.server) }}
                </template>
              </el-table-column>
              <el-table-column prop="timestamp" label="时间" align="center">
                <template #default="scope">
                  {{ formatTime(scope.row.timestamp) }}
                </template>
              </el-table-column>
              <el-table-column label="操作" align="center">
                <template #default="scope">
                  <el-button 
                    v-if="scope.row.video_file"
                    type="primary" 
                    link 
                    icon="VideoPlay"
                    @click="playVideo(scope.row.video_file)"
                  >
                    回放
                  </el-button>
                  <span v-else style="color: #909399; font-size: 12px;">无视频</span>
                </template>
              </el-table-column>
            </el-table>
          </el-card>
        </el-col>
      </el-row>
    </div>
    
    <el-empty v-else description="未选择比赛。">
       <el-button type="primary" @click="$router.push('/matches')">前往比赛列表</el-button>
    </el-empty>

    <!-- Video Player Dialog -->
    <el-dialog
      v-model="videoVisible"
      title="回合回放"
      width="60%"
      destroy-on-close
      center
      @close="stopVideo"
    >
      <div class="video-container">
        <video 
          ref="videoPlayer"
          controls 
          autoplay
          style="width: 100%; max-height: 500px;"
          :src="currentVideoUrl"
        >
          您的浏览器不支持视频播放。
        </video>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick, computed } from 'vue'
import { useRoute } from 'vue-router'
import * as echarts from 'echarts'
import { getMatchDetail, getMatchStats, getCurrentScore } from '@/api/match'

const route = useRoute()
const matchId = route.query.matchId
const rounds = ref([])
const matchInfo = ref({
  player1_name: 'Player 1',
  player2_name: 'Player 2'
})
const stats = ref({
  serve_success_rate: { player1: 0, player2: 0 },
  consecutive_score: { player1: 0, player2: 0 },
  average_rally_time: 0
})

const statsData = computed(() => {
  const s = stats.value || {}
  const cs = s.consecutive_score || { player1: 0, player2: 0 }
  return [
    {
      name: matchInfo.value.player1_name || 'Player 1',
      consecutive: cs.player1 || 0
    },
    {
      name: matchInfo.value.player2_name || 'Player 2',
      consecutive: cs.player2 || 0
    }
  ]
})

const videoVisible = ref(false)
const currentVideoUrl = ref('')
const videoPlayer = ref(null)

const fetchData = async () => {
  if (!matchId) return

  try {
    const [roundsData, fetchedStats, scoreData] = await Promise.all([
      getMatchDetail(matchId),
      getMatchStats(matchId),
      getCurrentScore(matchId)
    ])
    
    rounds.value = roundsData
    stats.value = fetchedStats
    if (scoreData) {
      matchInfo.value = scoreData
    }
    
    initChart()
  } catch (error) {
    console.error('Failed to fetch match details:', error)
  }
}

const getPlayerName = (key) => {
  if (key === 'player1' || key === 'Player 1') return matchInfo.value.player1_name || 'Player 1'
  if (key === 'player2' || key === 'Player 2') return matchInfo.value.player2_name || 'Player 2'
  return key
}

const formatTime = (timestamp) => {
  if (!timestamp) return ''
  const date = new Date(timestamp)
  const yyyy = date.getFullYear()
  const mm = String(date.getMonth() + 1).padStart(2, '0')
  const dd = String(date.getDate()).padStart(2, '0')
  const hh = String(date.getHours()).padStart(2, '0')
  const min = String(date.getMinutes()).padStart(2, '0')
  const ss = String(date.getSeconds()).padStart(2, '0')
  return `${yyyy}-${mm}-${dd} ${hh}:${min}:${ss}`
}

const playVideo = (fileName) => {
  if (!fileName) return
  // Construct API URL. Assuming API proxy is set up or base URL is known.
  // In development, Vite proxy handles /api.
  currentVideoUrl.value = `/api/video/${encodeURIComponent(fileName)}`
  videoVisible.value = true
}

const stopVideo = () => {
  if (videoPlayer.value) {
    videoPlayer.value.pause()
    videoPlayer.value.currentTime = 0
  }
  currentVideoUrl.value = ''
}

onMounted(() => {
  fetchData()
})

const initChart = () => {
  const chartDom = document.getElementById('chart-container')
  if (!chartDom) return
  // Dispose existing instance if any to avoid warning
  const existingInstance = echarts.getInstanceByDom(chartDom)
  if (existingInstance) {
    existingInstance.dispose()
  }
  
  const myChart = echarts.init(chartDom)
  const option = {
    height: '100%',
    tooltip: {},
    grid: { top: 30, bottom: 30, left: 40, right: 20 },
    xAxis: { data: [getPlayerName('player1'), getPlayerName('player2')] },
    yAxis: { max: 100 },
    series: [{
      name: '发球得分率 %',
      type: 'bar',
      data: [
        (stats.value.serve_success_rate.player1 * 100).toFixed(1),
        (stats.value.serve_success_rate.player2 * 100).toFixed(1)
      ],
      itemStyle: { color: '#409EFF' },
      barWidth: '50%'
    }]
  }
  myChart.setOption(option)
}
</script>

<style scoped>
.card-header {
  font-weight: bold;
}
</style>

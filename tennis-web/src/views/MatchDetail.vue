<template>
  <div class="match-detail-container">
    <div v-if="matchId" class="detail-content">
      <section class="detail-hero">
        <div>
          <p class="section-kicker">Match Analysis</p>
          <h1>{{ matchInfo.player1_name || 'Player 1' }} <span>VS</span> {{ matchInfo.player2_name || 'Player 2' }}</h1>
          <p>查看发球得分率、连续得分与每个回合的视频回放。</p>
        </div>
        <div class="match-id-card">
          <span>比赛 ID</span>
          <strong>{{ matchId }}</strong>
          <el-button type="primary" class="export-button" @click="downloadExcel">
            <el-icon><Download /></el-icon>
            导出 Excel
          </el-button>
        </div>
      </section>

      <div class="detail-grid">
        <el-card class="analysis-card" shadow="never">
          <template #header>
            <div class="card-header">
              <div>
                <span class="card-title">数据统计</span>
                <span class="card-subtitle">发球表现与回合效率</span>
              </div>
            </div>
          </template>
          <div id="chart-container" class="chart-container"></div>
          <div class="stats-summary">
            <div class="summary-item">
              <span class="summary-label">平均回合时间</span>
              <strong>{{ stats.average_rally_time }}s</strong>
            </div>
            <div class="summary-item">
              <span class="summary-label">{{ matchInfo.player1_name || 'Player 1' }} 连胜</span>
              <strong>{{ stats.consecutive_score?.player1 || 0 }}</strong>
            </div>
            <div class="summary-item">
              <span class="summary-label">{{ matchInfo.player2_name || 'Player 2' }} 连胜</span>
              <strong>{{ stats.consecutive_score?.player2 || 0 }}</strong>
            </div>
          </div>
          <el-table :data="statsData" class="stats-table" border style="width: 100%">
            <el-table-column prop="name" label="选手" align="center" />
            <el-table-column prop="consecutive" label="最大连胜分数" align="center" />
            <el-table-column label="平均回合时间" align="center">
              <template #default>
                {{ stats.average_rally_time }}s
              </template>
            </el-table-column>
          </el-table>
        </el-card>

        <el-card class="round-card" shadow="never">
          <template #header>
            <div class="card-header">
              <div>
                <span class="card-title">回合记录</span>
                <span class="card-subtitle">得分、发球和视频回放</span>
              </div>
            </div>
          </template>
          <el-table :data="rounds" class="round-table" style="width: 100%" height="520" stripe size="large">
            <el-table-column prop="rally_number" label="回合" width="80" align="center" />
            <el-table-column prop="scorer" label="得分方" align="center">
              <template #default="scope">
                <el-tag :type="scope.row.scorer === 'player1' ? 'primary' : 'warning'" round>
                  {{ getPlayerName(scope.row.scorer) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="server" label="发球方" align="center">
              <template #default="scope">
                {{ getPlayerName(scope.row.server) }}
              </template>
            </el-table-column>
            <el-table-column prop="timestamp" label="时间" align="center" min-width="170">
              <template #default="scope">
                {{ formatTime(scope.row.timestamp) }}
              </template>
            </el-table-column>
            <el-table-column label="操作" align="center" width="120">
              <template #default="scope">
                <el-button
                  v-if="scope.row.video_file"
                  type="primary"
                  plain
                  round
                  icon="VideoPlay"
                  @click="playVideo(scope.row.video_file)"
                >
                  回放
                </el-button>
                <span v-else class="muted-text">无视频</span>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </div>
    </div>
    
    <el-empty v-else class="empty-state" description="未选择比赛。">
       <el-button type="primary" @click="$router.push('/matches')">前往比赛列表</el-button>
    </el-empty>

    <el-dialog
      v-model="videoVisible"
      title="回合回放"
      width="60%"
      destroy-on-close
      center
      class="video-dialog"
      @close="stopVideo"
    >
      <div class="video-container">
        <video 
          ref="videoPlayer"
          controls 
          autoplay
          class="video-player"
          :src="currentVideoUrl"
        >
          您的浏览器不支持视频播放。
        </video>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import { useRoute } from 'vue-router'
import * as echarts from 'echarts'
import { ElMessage } from 'element-plus'
import { getMatchDetail, getMatchStats, getCurrentScore, exportMatch } from '@/api/match'

const route = useRoute()
const matchId = computed(() => route.query.matchId)
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
let detailRequestSeq = 0

const defaultStats = () => ({
  serve_success_rate: { player1: 0, player2: 0 },
  consecutive_score: { player1: 0, player2: 0 },
  average_rally_time: 0
})

const resetMatchData = () => {
  rounds.value = []
  stats.value = defaultStats()
  matchInfo.value = {
    player1_name: 'Player 1',
    player2_name: 'Player 2'
  }
}

const fetchData = async () => {
  const requestedMatchId = matchId.value
  if (!requestedMatchId) {
    resetMatchData()
    return
  }

  const requestSeq = ++detailRequestSeq

  try {
    const [roundsData, fetchedStats, scoreData] = await Promise.all([
      getMatchDetail(requestedMatchId),
      getMatchStats(requestedMatchId),
      getCurrentScore(requestedMatchId)
    ])
    if (requestSeq !== detailRequestSeq || String(matchId.value) !== String(requestedMatchId)) return
    
    rounds.value = roundsData
    stats.value = fetchedStats
    if (scoreData) {
      matchInfo.value = scoreData
    }
    
    initChart()
  } catch (error) {
    if (requestSeq === detailRequestSeq && String(matchId.value) === String(requestedMatchId)) {
      resetMatchData()
      initChart()
    }
    console.error('Failed to fetch match details:', error)
  }
}

const downloadExcel = async () => {
  if (!matchId.value) return

  try {
    const data = await exportMatch(matchId.value)
    const blob = data instanceof Blob
      ? data
      : new Blob([data], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' })
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `match_${matchId.value}.xlsx`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
  } catch (error) {
    console.error('Failed to export match:', error)
    ElMessage.error('导出失败，请确认后端服务是否可用')
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

watch(() => route.query.matchId, () => {
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
    color: ['#21b7a8', '#ffbd40'],
    tooltip: {
      backgroundColor: '#142033',
      borderWidth: 0,
      textStyle: { color: '#ffffff' }
    },
    grid: { top: 32, bottom: 36, left: 48, right: 20 },
    xAxis: {
      data: [getPlayerName('player1'), getPlayerName('player2')],
      axisLine: { lineStyle: { color: '#d7e1eb' } },
      axisLabel: { color: '#526276', fontWeight: 700 }
    },
    yAxis: {
      max: 100,
      splitLine: { lineStyle: { color: '#edf2f6' } },
      axisLabel: { color: '#7a8797' }
    },
    series: [{
      name: '发球得分率 %',
      type: 'bar',
      data: [
        (stats.value.serve_success_rate.player1 * 100).toFixed(1),
        (stats.value.serve_success_rate.player2 * 100).toFixed(1)
      ],
      itemStyle: {
        borderRadius: [12, 12, 4, 4],
        color: (params) => params.dataIndex === 0 ? '#21b7a8' : '#ffbd40'
      },
      barWidth: '46%'
    }]
  }
  myChart.setOption(option)
}
</script>

<style scoped>
.match-detail-container,
.detail-content {
  min-height: 100%;
}

.detail-content {
  display: flex;
  flex-direction: column;
  gap: 22px;
}

.detail-hero {
  display: flex;
  align-items: stretch;
  justify-content: space-between;
  gap: 20px;
  padding: 26px 28px;
  color: #fff;
  overflow: hidden;
  background:
    radial-gradient(circle at 88% 18%, rgba(51, 218, 203, 0.26), transparent 34%),
    linear-gradient(135deg, #141f32 0%, #173949 100%);
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

.detail-hero h1 {
  margin: 0;
  font-size: 30px;
  font-weight: 950;
  letter-spacing: 0;
}

.detail-hero h1 span {
  display: inline-grid;
  width: 54px;
  height: 34px;
  margin: 0 8px;
  place-items: center;
  color: #142033;
  font-size: 15px;
  font-style: italic;
  background: linear-gradient(135deg, #fff7cf, #33dacb);
  border-radius: 999px;
  vertical-align: middle;
}

.detail-hero p {
  max-width: 620px;
  margin: 10px 0 0;
  color: #c9d8e8;
  font-size: 14px;
}

.match-id-card {
  display: grid;
  min-width: 128px;
  padding: 18px 20px;
  place-items: center;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.16);
  border-radius: 16px;
}

.match-id-card > span {
  color: #c9d8e8;
  font-size: 12px;
  font-weight: 800;
}

.match-id-card strong {
  color: #fff2b8;
  font-size: 34px;
  font-weight: 950;
}

.export-button {
  margin-top: 10px;
  border: 0;
  background: linear-gradient(135deg, #33dacb, #ffbd40);
  color: #142033;
  font-weight: 850;
}

.detail-grid {
  display: grid;
  grid-template-columns: minmax(360px, 0.82fr) minmax(0, 1.18fr);
  gap: 22px;
  align-items: stretch;
}

.analysis-card,
.round-card {
  overflow: hidden;
  background: rgba(255, 255, 255, 0.94);
  border: 1px solid rgba(211, 221, 232, 0.9);
  border-radius: 18px;
  box-shadow: 0 18px 46px rgba(28, 42, 61, 0.09);
}

.analysis-card :deep(.el-card__header),
.round-card :deep(.el-card__header) {
  padding: 18px 22px;
  border-bottom-color: #edf2f6;
}

.analysis-card :deep(.el-card__body),
.round-card :deep(.el-card__body) {
  padding: 22px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
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

.chart-container {
  width: 100%;
  height: 350px;
  padding: 8px;
  background: linear-gradient(180deg, #f8fbfd, #ffffff);
  border: 1px solid #edf2f6;
  border-radius: 16px;
}

.stats-summary {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;
  margin: 16px 0;
}

.summary-item {
  padding: 14px;
  background: #f7fafc;
  border: 1px solid #e5edf4;
  border-radius: 14px;
}

.summary-label {
  display: block;
  color: #6b7b8f;
  font-size: 12px;
  font-weight: 800;
}

.summary-item strong {
  display: block;
  margin-top: 6px;
  color: #142033;
  font-size: 26px;
  font-weight: 950;
  line-height: 1;
}

.stats-table :deep(.el-table__header th),
.round-table :deep(.el-table__header th) {
  color: #5f6f82;
  background: #f7fafc;
  font-weight: 850;
}

.muted-text {
  color: #9aa7b6;
  font-size: 12px;
  font-weight: 700;
}

.empty-state {
  min-height: 420px;
  background: rgba(255, 255, 255, 0.88);
  border: 1px solid rgba(211, 221, 232, 0.86);
  border-radius: 18px;
}

.video-container {
  overflow: hidden;
  background: #111827;
  border-radius: 14px;
}

.video-player {
  display: block;
  width: 100%;
  max-height: 560px;
}

:deep(.video-dialog .el-dialog) {
  border-radius: 18px;
}

@media (max-width: 1120px) {
  .detail-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 760px) {
  .detail-hero {
    flex-direction: column;
    padding: 22px;
  }

  .detail-hero h1 {
    font-size: 24px;
  }

  .stats-summary {
    grid-template-columns: 1fr;
  }
}
</style>

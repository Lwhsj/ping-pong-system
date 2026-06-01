<template>
  <div class="history-page">
    <section class="page-hero">
      <div>
        <p class="section-kicker">Match Archive</p>
        <h1>比赛历史</h1>
        <p class="hero-copy">按日期、选手和状态快速定位比赛，查看每一场对阵的详细数据。</p>
      </div>
      <div class="hero-stat">
        <span class="hero-stat-number">{{ matches.length }}</span>
        <span class="hero-stat-label">场比赛</span>
      </div>
    </section>

    <section class="filter-panel">
      <el-form :inline="true" :model="filters" class="filter-form">
        <el-form-item label="日期范围">
          <el-date-picker
            v-model="filters.dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item label="选手">
          <el-input v-model="filters.player" placeholder="选手姓名" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="filters.status" placeholder="选择状态" style="width: 150px">
            <el-option label="全部" value="" />
            <el-option label="进行中" value="started" />
            <el-option label="已结束" value="finished" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" class="search-button" @click="handleFilter">
            <el-icon><Search /></el-icon>
            <span>查询</span>
          </el-button>
        </el-form-item>
      </el-form>
    </section>

    <el-card class="table-card" shadow="never">
      <template #header>
        <div class="card-header">
          <div>
            <span class="card-title">赛事列表</span>
            <span class="card-subtitle">所有匹配当前筛选条件的比赛</span>
          </div>
        </div>
      </template>

    <el-table :data="matches" class="match-table" style="width: 100%" stripe size="large">
      <el-table-column prop="id" label="ID" width="90" sortable />
      <el-table-column prop="date" label="日期" width="140" sortable />
      <el-table-column label="对阵双方" min-width="200">
        <template #default="scope">
          <div class="vs-container">
             <span class="player-name p1">{{ scope.row.player1_name }}</span>
             <span class="vs-tag">VS</span>
             <span class="player-name p2">{{ scope.row.player2_name }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="120">
        <template #default="scope">
          <el-tag
            class="status-tag"
            :type="scope.row.status === 'finished' ? 'success' : 'danger'"
            effect="light"
            round
          >
            {{ scope.row.status === 'finished' ? '已结束' : '进行中' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="150" fixed="right" align="center">
        <template #default="scope">
          <el-button 
            size="small" 
            type="primary"
            plain
            class="detail-button"
            @click="goToDetail(scope.row.id)"
          >
            <el-icon><View /></el-icon>
            详情
          </el-button>
        </template>
      </el-table-column>
    </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getMatchList } from '@/api/match'

const router = useRouter()
const filters = ref({
  dateRange: '',
  player: '',
  status: ''
})

const matches = ref([])

const fetchMatches = async () => {
  try {
    const params = {
      player: filters.value.player,
      status: filters.value.status
    }
    if (filters.value.dateRange && filters.value.dateRange.length === 2) {
      params.dateFrom = filters.value.dateRange[0]
      params.dateTo = filters.value.dateRange[1]
    }
    const data = await getMatchList(params)
    matches.value = data
  } catch (error) {
    console.error('Failed to fetch matches:', error)
  }
}

onMounted(() => {
  fetchMatches()
})

const handleFilter = () => {
  fetchMatches()
}

const goToDetail = (id) => {
  router.push({ path: '/match-detail', query: { matchId: id } })
}
</script>

<style scoped>
.history-page {
  display: flex;
  flex-direction: column;
  gap: 22px;
}

.page-hero {
  display: flex;
  align-items: stretch;
  justify-content: space-between;
  gap: 20px;
  padding: 26px 28px;
  color: #fff;
  overflow: hidden;
  background:
    linear-gradient(135deg, rgba(19, 32, 51, 0.96), rgba(25, 57, 75, 0.94)),
    radial-gradient(circle at 88% 20%, rgba(51, 218, 203, 0.32), transparent 34%);
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

.page-hero h1 {
  margin: 0;
  font-size: 32px;
  font-weight: 900;
  letter-spacing: 0;
}

.hero-copy {
  max-width: 520px;
  margin: 10px 0 0;
  color: #c9d8e8;
  font-size: 14px;
}

.hero-stat {
  display: grid;
  min-width: 138px;
  padding: 18px 20px;
  place-items: center;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.16);
  border-radius: 16px;
}

.hero-stat-number {
  color: #fff2b8;
  font-size: 38px;
  font-weight: 950;
  line-height: 1;
}

.hero-stat-label {
  margin-top: 7px;
  color: #c9d8e8;
  font-size: 13px;
  font-weight: 750;
}

.filter-panel,
.table-card {
  background: rgba(255, 255, 255, 0.92);
  border: 1px solid rgba(211, 221, 232, 0.86);
  border-radius: 16px;
  box-shadow: 0 16px 42px rgba(28, 42, 61, 0.08);
}

.filter-panel {
  padding: 20px 22px 2px;
}

.filter-form {
  display: flex;
  align-items: flex-start;
  flex-wrap: wrap;
  gap: 4px 12px;
}

.filter-form :deep(.el-form-item__label) {
  color: #58687c;
  font-weight: 800;
}

.filter-form :deep(.el-input__wrapper),
.filter-form :deep(.el-select__wrapper),
.filter-form :deep(.el-date-editor) {
  border-radius: 10px;
  box-shadow: 0 0 0 1px #d8e1ea inset;
}

.search-button {
  gap: 6px;
  min-width: 96px;
  border: none;
  background: linear-gradient(135deg, #1fb8a9, #1989d6);
  border-radius: 10px;
  font-weight: 800;
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
  font-weight: 650;
}

.table-card {
  overflow: hidden;
}

.table-card :deep(.el-card__header) {
  padding: 18px 22px;
  border-bottom-color: #edf2f6;
}

.table-card :deep(.el-card__body) {
  padding: 0;
}

.match-table :deep(.el-table__header th) {
  color: #5f6f82;
  background: #f7fafc;
  font-weight: 850;
}

.match-table :deep(.el-table__cell) {
  padding: 15px 0;
}

.vs-container {
  display: grid;
  grid-template-columns: 1fr auto 1fr;
  align-items: center;
  gap: 14px;
  width: 100%;
}

.player-name {
  min-width: 0;
  font-weight: 850;
  font-size: 15px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.player-name.p1 {
  text-align: right;
  color: #172338;
}

.player-name.p2 {
  text-align: left;
  color: #172338;
}

.vs-tag {
  color: #1d6f67;
  font-weight: 950;
  font-style: italic;
  padding: 4px 10px;
  background-color: #e8fbf8;
  border: 1px solid #c7f1eb;
  border-radius: 999px;
  font-size: 12px;
}

.status-tag {
  font-weight: 850;
}

.detail-button {
  border-radius: 999px;
  font-weight: 800;
}

@media (max-width: 760px) {
  .page-hero {
    flex-direction: column;
    padding: 22px;
  }

  .filter-panel {
    padding-bottom: 12px;
  }

  .filter-form {
    display: block;
  }
}
</style>

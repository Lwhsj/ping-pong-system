<template>
  <el-card class="box-card">
    <template #header>
      <div class="card-header">
        <span>比赛历史</span>
      </div>
    </template>
    
    <div class="filter-container">
      <el-form :inline="true" :model="filters" class="demo-form-inline">
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
          <el-button type="primary" @click="handleFilter">
            <el-icon><Search /></el-icon> 查询
          </el-button>
        </el-form-item>
      </el-form>
    </div>

    <el-table :data="matches" style="width: 100%" stripe border size="large">
      <el-table-column prop="id" label="ID" width="80" sortable />
      <el-table-column prop="date" label="日期" width="120" sortable />
      <el-table-column label="对阵双方" min-width="200">
        <template #default="scope">
          <div style="display: flex; align-items: center; justify-content: space-around;">
             <span style="font-weight: bold;">{{ scope.row.player1_name }}</span>
             <span style="color: #909399;">VS</span>
             <span style="font-weight: bold;">{{ scope.row.player2_name }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template #default="scope">
          <el-tag :type="scope.row.status === 'finished' ? 'success' : 'danger'">
            {{ scope.row.status === 'finished' ? '已结束' : '进行中' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="scope">
          <el-button 
            size="small" 
            @click="goToDetail(scope.row.id)"
          >
            详情
          </el-button>
        </template>
      </el-table-column>
    </el-table>
  </el-card>
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
.filter-container {
  margin-bottom: 24px;
}
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.vs-container {
  display: grid;
  grid-template-columns: 1fr auto 1fr;
  align-items: center;
  gap: 12px;
  width: 100%;
}

.player-name {
  font-weight: 600;
  font-size: 14px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.player-name.p1 {
  text-align: right;
  color: #409EFF;
}

.player-name.p2 {
  text-align: left;
  color: #F56C6C;
}

.vs-tag {
  color: #909399;
  font-weight: bold;
  font-style: italic;
  padding: 0 8px;
  background-color: #f4f4f5;
  border-radius: 4px;
  font-size: 12px;
}
</style>

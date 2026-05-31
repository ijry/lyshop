<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getDistributionOrders } from '@/api/distribution'
import { ElMessage } from 'element-plus'

const loading = ref(false)
const list = ref<any[]>([])
const total = ref(0)
const page = ref(1)
const size = ref(20)
const statusFilter = ref('')

const statusOptions = [
  { label: '全部', value: '' },
  { label: '待结算', value: 'pending' },
  { label: '已结算', value: 'settled' },
  { label: '已取消', value: 'cancelled' }
]

const statusMap: Record<string, { label: string; type: string }> = {
  pending: { label: '待结算', type: 'warning' },
  settled: { label: '已结算', type: 'success' },
  cancelled: { label: '已取消', type: 'info' }
}

async function loadData() {
  loading.value = true
  try {
    const res = await getDistributionOrders({
      page: page.value,
      size: size.value,
      status: statusFilter.value
    })
    list.value = res.list || []
    total.value = res.total || 0
  } catch (error: any) {
    ElMessage.error(error.message || '加载失败')
  } finally {
    loading.value = false
  }
}

function handlePageChange(val: number) {
  page.value = val
  loadData()
}

onMounted(() => loadData())
</script>

<template>
  <div class="page-container">
    <div class="toolbar">
      <el-select v-model="statusFilter" placeholder="状态筛选" @change="loadData" style="width: 150px">
        <el-option v-for="opt in statusOptions" :key="opt.value" :label="opt.label" :value="opt.value" />
      </el-select>
    </div>

    <el-table :data="list" v-loading="loading" border stripe>
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="order_id" label="订单ID" width="100" />
      <el-table-column prop="distributor_id" label="分销商ID" width="120" />
      <el-table-column prop="level" label="分销层级" width="100">
        <template #default="{ row }">{{ row.level }}级</template>
      </el-table-column>
      <el-table-column prop="order_amount" label="订单金额" width="120">
        <template #default="{ row }">¥{{ row.order_amount?.toFixed(2) }}</template>
      </el-table-column>
      <el-table-column prop="commission_rate" label="佣金比例" width="100">
        <template #default="{ row }">{{ row.commission_rate }}%</template>
      </el-table-column>
      <el-table-column prop="commission" label="佣金金额" width="120">
        <template #default="{ row }">¥{{ row.commission?.toFixed(2) }}</template>
      </el-table-column>
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="statusMap[row.status]?.type">{{ statusMap[row.status]?.label }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="settled_at" label="结算时间" width="180" />
      <el-table-column prop="created_at" label="创建时间" width="180" />
    </el-table>

    <el-pagination
      v-model:current-page="page"
      :page-size="size"
      :total="total"
      layout="total, prev, pager, next"
      @current-change="handlePageChange"
      style="margin-top: 20px; justify-content: center"
    />
  </div>
</template>

<style scoped>
.page-container {
  padding: 20px;
}
.toolbar {
  margin-bottom: 20px;
}
</style>

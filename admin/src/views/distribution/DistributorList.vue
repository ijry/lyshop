<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getDistributors, updateDistributor } from '@/api/distribution'
import { ElMessage, ElMessageBox } from 'element-plus'

const loading = ref(false)
const list = ref<any[]>([])
const total = ref(0)
const page = ref(1)
const size = ref(20)
const statusFilter = ref('')

const statusOptions = [
  { label: '全部', value: '' },
  { label: '待审核', value: 'pending' },
  { label: '已激活', value: 'active' },
  { label: '已禁用', value: 'disabled' }
]

const statusMap: Record<string, { label: string; type: string }> = {
  pending: { label: '待审核', type: 'warning' },
  active: { label: '已激活', type: 'success' },
  disabled: { label: '已禁用', type: 'info' }
}

async function loadData() {
  loading.value = true
  try {
    const res = await getDistributors({
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

async function handleStatusChange(row: any, status: string) {
  try {
    await ElMessageBox.confirm(`确认${status === 'active' ? '激活' : '禁用'}该分销商？`, '提示')
    await updateDistributor(row.id, { status })
    ElMessage.success('操作成功')
    loadData()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '操作失败')
    }
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
      <el-table-column prop="user_id" label="用户ID" width="100" />
      <el-table-column prop="real_name" label="姓名" width="120" />
      <el-table-column prop="phone" label="手机号" width="130" />
      <el-table-column prop="level" label="等级" width="80" />
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="statusMap[row.status]?.type">{{ statusMap[row.status]?.label }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="total_earnings" label="累计收益" width="120">
        <template #default="{ row }">¥{{ row.total_earnings?.toFixed(2) }}</template>
      </el-table-column>
      <el-table-column prop="available_amount" label="可提现" width="120">
        <template #default="{ row }">¥{{ row.available_amount?.toFixed(2) }}</template>
      </el-table-column>
      <el-table-column prop="frozen_amount" label="冻结金额" width="120">
        <template #default="{ row }">¥{{ row.frozen_amount?.toFixed(2) }}</template>
      </el-table-column>
      <el-table-column prop="withdrawn_amount" label="已提现" width="120">
        <template #default="{ row }">¥{{ row.withdrawn_amount?.toFixed(2) }}</template>
      </el-table-column>
      <el-table-column prop="total_customers" label="客户数" width="100" />
      <el-table-column prop="total_orders" label="订单数" width="100" />
      <el-table-column prop="created_at" label="创建时间" width="180" />
      <el-table-column label="操作" width="180" fixed="right">
        <template #default="{ row }">
          <el-button v-if="row.status === 'pending'" type="primary" size="small" @click="handleStatusChange(row, 'active')">激活</el-button>
          <el-button v-if="row.status === 'active'" type="warning" size="small" @click="handleStatusChange(row, 'disabled')">禁用</el-button>
          <el-button v-if="row.status === 'disabled'" type="success" size="small" @click="handleStatusChange(row, 'active')">启用</el-button>
        </template>
      </el-table-column>
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

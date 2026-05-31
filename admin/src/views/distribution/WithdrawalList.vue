<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getWithdrawals, approveWithdrawal, rejectWithdrawal, completeWithdrawal } from '@/api/distribution'
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
  { label: '已通过', value: 'approved' },
  { label: '已拒绝', value: 'rejected' },
  { label: '已完成', value: 'completed' }
]

const statusMap: Record<string, { label: string; type: string }> = {
  pending: { label: '待审核', type: 'warning' },
  approved: { label: '已通过', type: 'primary' },
  rejected: { label: '已拒绝', type: 'danger' },
  completed: { label: '已完成', type: 'success' }
}

async function loadData() {
  loading.value = true
  try {
    const res = await getWithdrawals({
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

async function handleApprove(row: any) {
  try {
    await ElMessageBox.confirm('确认通过该提现申请？', '提示')
    await approveWithdrawal(row.id)
    ElMessage.success('操作成功')
    loadData()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '操作失败')
    }
  }
}

async function handleReject(row: any) {
  try {
    const { value } = await ElMessageBox.prompt('请输入拒绝原因', '拒绝提现', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      inputPattern: /.+/,
      inputErrorMessage: '请输入拒绝原因'
    })
    await rejectWithdrawal(row.id, value)
    ElMessage.success('操作成功')
    loadData()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '操作失败')
    }
  }
}

async function handleComplete(row: any) {
  try {
    await ElMessageBox.confirm('确认已完成打款？', '提示')
    await completeWithdrawal(row.id)
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
      <el-table-column prop="distributor_id" label="分销商ID" width="120" />
      <el-table-column prop="amount" label="提现金额" width="120">
        <template #default="{ row }">¥{{ row.amount?.toFixed(2) }}</template>
      </el-table-column>
      <el-table-column prop="fee" label="手续费" width="100">
        <template #default="{ row }">¥{{ row.fee?.toFixed(2) }}</template>
      </el-table-column>
      <el-table-column prop="actual_amount" label="实际到账" width="120">
        <template #default="{ row }">¥{{ row.actual_amount?.toFixed(2) }}</template>
      </el-table-column>
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="statusMap[row.status]?.type">{{ statusMap[row.status]?.label }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="bank_name" label="银行" width="120" />
      <el-table-column prop="bank_account" label="账号" width="180" />
      <el-table-column prop="account_name" label="户名" width="120" />
      <el-table-column prop="reject_reason" label="拒绝原因" width="150" show-overflow-tooltip />
      <el-table-column prop="created_at" label="申请时间" width="180" />
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button v-if="row.status === 'pending'" type="primary" size="small" @click="handleApprove(row)">通过</el-button>
          <el-button v-if="row.status === 'pending'" type="danger" size="small" @click="handleReject(row)">拒绝</el-button>
          <el-button v-if="row.status === 'approved'" type="success" size="small" @click="handleComplete(row)">完成</el-button>
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

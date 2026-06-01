<template>
  <div class="page-container">
    <div class="header">
      <h2>砍价活动管理</h2>
      <el-button type="primary" @click="openCreateDialog">新增活动</el-button>
    </div>

    <el-table :data="activities" border stripe>
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="活动名称" min-width="200" />
      <el-table-column label="活动时间" min-width="300">
        <template #default="{ row }">
          {{ formatTime(row.start_at) }} ~ {{ formatTime(row.end_at) }}
        </template>
      </el-table-column>
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'info'">
            {{ row.status === 1 ? '启用' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="280" fixed="right">
        <template #default="{ row }">
          <el-button size="small" @click="goToProducts(row.id)">商品管理</el-button>
          <el-button size="small" @click="openEditDialog(row)">编辑</el-button>
          <el-button size="small" type="danger" @click="handleDelete(row.id)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      v-model:current-page="page"
      v-model:page-size="size"
      :total="total"
      layout="total, prev, pager, next"
      @current-change="loadActivities"
    />

    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑活动' : '新增活动'" width="600px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="活动名称" required>
          <el-input v-model="form.name" placeholder="请输入活动名称" />
        </el-form-item>
        <el-form-item label="开始时间" required>
          <el-date-picker
            v-model="form.start_at"
            type="datetime"
            placeholder="选择开始时间"
            value-format="YYYY-MM-DDTHH:mm:ss"
          />
        </el-form-item>
        <el-form-item label="结束时间" required>
          <el-date-picker
            v-model="form.end_at"
            type="datetime"
            placeholder="选择结束时间"
            value-format="YYYY-MM-DDTHH:mm:ss"
          />
        </el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="form.status" :active-value="1" :inactive-value="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSave">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/api/request'

const router = useRouter()

const activities = ref<any[]>([])
const total = ref(0)
const page = ref(1)
const size = ref(20)

const dialogVisible = ref(false)
const editingId = ref(0)
const form = ref({
  name: '',
  start_at: '',
  end_at: '',
  status: 1,
})

function formatTime(time: string) {
  if (!time) return '-'
  return time.replace('T', ' ').slice(0, 19)
}

async function loadActivities() {
  const data = await request.get<any>('/bargain/activities', {
    params: {
      page: page.value,
      size: size.value,
    }
  })
  activities.value = data.list || []
  total.value = data.total || 0
}

function openCreateDialog() {
  editingId.value = 0
  form.value = {
    name: '',
    start_at: '',
    end_at: '',
    status: 1,
  }
  dialogVisible.value = true
}

function openEditDialog(row: any) {
  editingId.value = row.id
  form.value = {
    name: row.name,
    start_at: row.start_at,
    end_at: row.end_at,
    status: row.status,
  }
  dialogVisible.value = true
}

async function handleSave() {
  if (!form.value.name) {
    ElMessage.error('请输入活动名称')
    return
  }
  if (!form.value.start_at || !form.value.end_at) {
    ElMessage.error('请选择活动时间')
    return
  }

  try {
    if (editingId.value) {
      await request.put(`/bargain/activities/${editingId.value}`, form.value)
      ElMessage.success('更新成功')
    } else {
      await request.post('/bargain/activities', form.value)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    loadActivities()
  } catch (error: any) {
    ElMessage.error(error.message || '操作失败')
  }
}

async function handleDelete(id: number) {
  try {
    await ElMessageBox.confirm('确定删除该活动吗？', '提示', {
      type: 'warning',
    })
    await request.delete(`/bargain/activities/${id}`)
    ElMessage.success('删除成功')
    loadActivities()
  } catch (error) {
    // 用户取消
  }
}

function goToProducts(activityId: number) {
  router.push(`/bargain/products?activity_id=${activityId}`)
}

onMounted(() => {
  loadActivities()
})
</script>

<style scoped>
.page-container {
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.el-pagination {
  margin-top: 20px;
  justify-content: flex-end;
}
</style>

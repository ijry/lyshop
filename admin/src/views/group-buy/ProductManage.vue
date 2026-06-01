<template>
  <div class="page-container">
    <div class="header">
      <h2>拼团商品管理</h2>
      <el-button type="primary" @click="openAddDialog">添加商品</el-button>
      <el-button @click="handleBatchSave">批量保存</el-button>
    </div>

    <el-form inline style="margin-bottom: 20px">
      <el-form-item label="活动">
        <el-select v-model="activityId" placeholder="选择活动" @change="loadProducts">
          <el-option
            v-for="act in activities"
            :key="act.id"
            :label="act.name"
            :value="act.id"
          />
        </el-select>
      </el-form-item>
    </el-form>

    <el-table :data="products" border stripe>
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="product_id" label="商品ID" width="100" />
      <el-table-column prop="sku_id" label="SKU ID" width="100" />
      <el-table-column label="拼团价" width="120">
        <template #default="{ row }">
          <el-input-number v-model="row.group_price" :min="0" :precision="2" size="small" />
        </template>
      </el-table-column>
      <el-table-column label="限购数量" width="120">
        <template #default="{ row }">
          <el-input-number v-model="row.limit_per_order" :min="0" size="small" />
        </template>
      </el-table-column>
      <el-table-column label="活动库存" width="120">
        <template #default="{ row }">
          <el-input-number v-model="row.total_stock_limit" :min="0" size="small" />
        </template>
      </el-table-column>
      <el-table-column prop="sold_qty" label="已售" width="100" />
      <el-table-column label="操作" width="100" fixed="right">
        <template #default="{ $index }">
          <el-button size="small" type="danger" @click="handleRemove($index)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="addDialogVisible" title="添加商品" width="500px">
      <el-form :model="addForm" label-width="100px">
        <el-form-item label="商品ID" required>
          <el-input-number v-model="addForm.product_id" :min="1" style="width: 100%" />
        </el-form-item>
        <el-form-item label="SKU ID">
          <el-input-number v-model="addForm.sku_id" :min="0" style="width: 100%" />
          <div style="font-size: 12px; color: #999; margin-top: 5px">0表示全部SKU</div>
        </el-form-item>
        <el-form-item label="拼团价" required>
          <el-input-number v-model="addForm.group_price" :min="0" :precision="2" style="width: 100%" />
        </el-form-item>
        <el-form-item label="限购数量">
          <el-input-number v-model="addForm.limit_per_order" :min="0" style="width: 100%" />
          <div style="font-size: 12px; color: #999; margin-top: 5px">0表示不限购</div>
        </el-form-item>
        <el-form-item label="活动库存">
          <el-input-number v-model="addForm.total_stock_limit" :min="0" style="width: 100%" />
          <div style="font-size: 12px; color: #999; margin-top: 5px">0表示不限库存</div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="addDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleAdd">添加</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import request from '@/api/request'

const route = useRoute()

const activities = ref<any[]>([])
const activityId = ref(0)
const products = ref<any[]>([])

const addDialogVisible = ref(false)
const addForm = ref({
  product_id: 0,
  sku_id: 0,
  group_price: 0,
  limit_per_order: 0,
  total_stock_limit: 0,
})

async function loadActivities() {
  const data = await request.get<any>('/group-buy/activities', { params: { page: 1, size: 100 } })
  activities.value = data.list || []

  // 从URL获取activity_id
  const urlActivityId = Number(route.query.activity_id)
  if (urlActivityId && activities.value.find(a => a.id === urlActivityId)) {
    activityId.value = urlActivityId
    loadProducts()
  } else if (activities.value.length > 0) {
    activityId.value = activities.value[0].id
    loadProducts()
  }
}

async function loadProducts() {
  if (!activityId.value) return

  const data = await request.get<any>('/group-buy/products', {
    params: {
      activity_id: activityId.value,
      page: 1,
      size: 1000,
    }
  })
  products.value = data.list || []
}

function openAddDialog() {
  if (!activityId.value) {
    ElMessage.error('请先选择活动')
    return
  }
  addForm.value = {
    product_id: 0,
    sku_id: 0,
    group_price: 0,
    limit_per_order: 0,
    total_stock_limit: 0,
  }
  addDialogVisible.value = true
}

function handleAdd() {
  if (!addForm.value.product_id || addForm.value.group_price <= 0) {
    ElMessage.error('请填写完整信息')
    return
  }

  // 检查是否已存在
  const exists = products.value.find(
    p => p.product_id === addForm.value.product_id && p.sku_id === addForm.value.sku_id
  )
  if (exists) {
    ElMessage.error('该商品已存在')
    return
  }

  products.value.push({ ...addForm.value })
  addDialogVisible.value = false
  ElMessage.success('添加成功，请点击批量保存')
}

function handleRemove(index: number) {
  products.value.splice(index, 1)
  ElMessage.success('删除成功，请点击批量保存')
}

async function handleBatchSave() {
  if (!activityId.value) {
    ElMessage.error('请先选择活动')
    return
  }

  try {
    await request.put(`/group-buy/activities/${activityId.value}/products`, products.value)
    ElMessage.success('保存成功')
    loadProducts()
  } catch (error: any) {
    ElMessage.error(error.message || '保存失败')
  }
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

.header h2 {
  margin: 0;
}

.header .el-button {
  margin-left: 10px;
}
</style>

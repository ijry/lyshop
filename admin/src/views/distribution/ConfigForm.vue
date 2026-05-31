<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getDistributionConfig, updateDistributionConfig } from '@/api/distribution'
import { ElMessage } from 'element-plus'

const loading = ref(false)
const saving = ref(false)
const form = ref({
  enabled: true,
  level: 2,
  level1_rate: 10,
  level2_rate: 5,
  level3_rate: 2,
  min_withdraw: 100,
  withdraw_fee_rate: 0,
  auto_approve: false,
  require_real_name: true
})

async function loadData() {
  loading.value = true
  try {
    const res = await getDistributionConfig()
    if (res) {
      form.value = res
    }
  } catch (error: any) {
    ElMessage.error(error.message || '加载失败')
  } finally {
    loading.value = false
  }
}

async function handleSave() {
  saving.value = true
  try {
    await updateDistributionConfig(form.value)
    ElMessage.success('保存成功')
  } catch (error: any) {
    ElMessage.error(error.message || '保存失败')
  } finally {
    saving.value = false
  }
}

onMounted(() => loadData())
</script>

<template>
  <div class="page-container" v-loading="loading">
    <el-form :model="form" label-width="140px" style="max-width: 600px">
      <el-form-item label="启用分销">
        <el-switch v-model="form.enabled" />
      </el-form-item>

      <el-form-item label="分销层级">
        <el-radio-group v-model="form.level">
          <el-radio :label="1">一级</el-radio>
          <el-radio :label="2">二级</el-radio>
          <el-radio :label="3">三级</el-radio>
        </el-radio-group>
      </el-form-item>

      <el-form-item label="一级佣金比例">
        <el-input-number v-model="form.level1_rate" :min="0" :max="100" :precision="2" />
        <span style="margin-left: 10px">%</span>
      </el-form-item>

      <el-form-item label="二级佣金比例" v-if="form.level >= 2">
        <el-input-number v-model="form.level2_rate" :min="0" :max="100" :precision="2" />
        <span style="margin-left: 10px">%</span>
      </el-form-item>

      <el-form-item label="三级佣金比例" v-if="form.level >= 3">
        <el-input-number v-model="form.level3_rate" :min="0" :max="100" :precision="2" />
        <span style="margin-left: 10px">%</span>
      </el-form-item>

      <el-form-item label="最低提现金额">
        <el-input-number v-model="form.min_withdraw" :min="0" :precision="2" />
        <span style="margin-left: 10px">元</span>
      </el-form-item>

      <el-form-item label="提现手续费比例">
        <el-input-number v-model="form.withdraw_fee_rate" :min="0" :max="100" :precision="2" />
        <span style="margin-left: 10px">%</span>
      </el-form-item>

      <el-form-item label="自动审核分销商">
        <el-switch v-model="form.auto_approve" />
      </el-form-item>

      <el-form-item label="要求实名认证">
        <el-switch v-model="form.require_real_name" />
      </el-form-item>

      <el-form-item>
        <el-button type="primary" @click="handleSave" :loading="saving">保存配置</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<style scoped>
.page-container {
  padding: 20px;
}
</style>

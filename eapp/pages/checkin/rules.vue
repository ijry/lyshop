<script setup lang="ts">
import { onShow } from '@dcloudio/uni-app'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { getCheckinRules, saveCheckinRules } from '@/api/checkin'

const { t } = useI18n()
const loading = ref(false)
const rules = ref<Array<{ day: number; points: number }>>([])

async function loadRules() {
  loading.value = true
  try {
    const res = await getCheckinRules()
    rules.value = Array.isArray(res) ? res.map((r: any) => ({ day: Number(r.day || 0), points: Number(r.points || 0) })) : []
  } finally {
    loading.value = false
  }
}

function addRule() {
  rules.value.push({ day: 0, points: 10 })
}

function removeRule(index: number) {
  rules.value.splice(index, 1)
}

async function save() {
  await saveCheckinRules(rules.value)
  uni.showToast({ title: '保存成功', icon: 'success' })
}

onShow(() => loadRules())
</script>

<template>
  <view class="page">
    <view class="top-bar">
      <text class="title">{{ t('checkin.rules') }}</text>
      <up-button size="mini" type="primary" @click="addRule">添加规则</up-button>
    </view>
    <view v-if="!loading && !rules.length" class="empty">暂无签到规则</view>
    <view v-for="(rule, index) in rules" :key="index" class="card">
      <view class="form-row">
        <view class="form-item">
          <text class="label">天数（0=默认）</text>
          <up-input v-model="(rule as any).day" type="number" placeholder="天数" />
        </view>
        <view class="form-item">
          <text class="label">积分</text>
          <up-input v-model="(rule as any).points" type="number" placeholder="积分" />
        </view>
        <up-button size="mini" type="error" plain @click="removeRule(index)">删除</up-button>
      </view>
    </view>
    <view v-if="rules.length" class="save-bar">
      <up-button type="primary" @click="save">{{ t('common.save') }}</up-button>
    </view>
    <view v-if="loading" class="loading">加载中…</view>
  </view>
</template>

<style scoped>
.page { min-height: 100vh; background: var(--eapp-bg); padding: 20rpx; box-sizing: border-box; display: grid; gap: 14rpx; align-content: start; }
.top-bar { display: flex; align-items: center; justify-content: space-between; }
.title { font-size: 32rpx; font-weight: 700; }
.card { background: #fff; border: 1px solid var(--eapp-border); border-radius: 20rpx; padding: 20rpx; }
.form-row { display: flex; align-items: flex-end; gap: 12rpx; }
.form-item { flex: 1; }
.label { font-size: 22rpx; color: #475569; margin-bottom: 4rpx; display: block; }
.save-bar { margin-top: 10rpx; }
.empty { text-align: center; color: var(--eapp-text-muted); padding: 80rpx 0; }
.loading { text-align: center; color: var(--eapp-text-muted); padding: 20rpx 0; }
</style>

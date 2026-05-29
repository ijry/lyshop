<script setup lang="ts">
import { reactive, ref } from 'vue'

const props = defineProps<{ show: boolean; coupon?: any }>()
const emit = defineEmits<{ (e: 'close'): void; (e: 'submit', payload: any): void }>()

const form = reactive({
  name: '', type: 1, min_amount: '', discount: '',
  total_count: '', per_limit: '1', status: 1,
  start_at: '', end_at: '', description: '',
  stack_rule: 'exclusive', target_type: 'all', target_value: '',
})

const stackRuleOptions = [
  { label: '互斥', value: 'exclusive' },
  { label: '同类可叠', value: 'same_type' },
  { label: '跨类可叠', value: 'cross_type' },
]
const targetTypeOptions = [
  { label: '全部用户', value: 'all' },
  { label: '会员等级', value: 'vip_level' },
  { label: '新用户', value: 'new_user' },
]

import { watch } from 'vue'
watch(() => props.coupon, (c) => {
  if (c) {
    form.name = String(c.name || '')
    form.type = Number(c.type || 1)
    form.min_amount = String(c.min_amount || '')
    form.discount = String(c.discount || '')
    form.total_count = String(c.total_count || '')
    form.per_limit = String(c.per_limit || '1')
    form.status = Number(c.status || 0) === 1 ? 1 : 0
    form.start_at = String(c.start_at || '')
    form.end_at = String(c.end_at || '')
    form.description = String(c.description || '')
    form.stack_rule = String(c.stack_rule || 'exclusive')
    form.target_type = String(c.target_type || 'all')
    form.target_value = String(c.target_value || '')
  } else {
    form.name = ''; form.type = 1; form.min_amount = ''; form.discount = ''
    form.total_count = ''; form.per_limit = '1'; form.status = 1
    form.start_at = ''; form.end_at = ''; form.description = ''
    form.stack_rule = 'exclusive'; form.target_type = 'all'; form.target_value = ''
  }
}, { immediate: true })

function submit() {
  if (!form.name.trim()) { uni.showToast({ title: '请输入名称', icon: 'none' }); return }
  emit('submit', {
    name: form.name.trim(),
    type: form.type,
    min_amount: Number(form.min_amount || 0),
    discount: Number(form.discount || 0),
    total_count: Number(form.total_count || 0),
    per_limit: Number(form.per_limit || 1),
    status: form.status,
    start_at: form.start_at,
    end_at: form.end_at,
    description: form.description,
    stack_rule: form.stack_rule,
    target_type: form.target_type,
    target_value: form.target_value,
  })
}
</script>

<template>
  <up-popup :show="show" mode="bottom" round="16" @close="$emit('close')">
    <view class="popup-body">
      <view class="popup-title">{{ coupon ? '编辑优惠券' : '新建优惠券' }}</view>
      <up-input v-model="form.name" placeholder="优惠券名称" clearable />
      <view class="mt" />
      <picker mode="selector" :range="[{ label: '满减', v: 1 }, { label: '折扣', v: 2 }, { label: '立减', v: 3 }]" range-key="label" @change="(e: any) => form.type = [1,2,3][e.detail.value]">
        <view class="picker">类型：{{ { 1: '满减', 2: '折扣', 3: '立减' }[form.type] || '满减' }}</view>
      </picker>
      <view class="mt" />
      <up-input v-model="form.min_amount" type="digit" placeholder="最低消费金额" />
      <view class="mt" />
      <up-input v-model="form.discount" type="digit" placeholder="优惠值（金额或折扣率）" />
      <view class="mt" />
      <up-input v-model="form.total_count" type="number" placeholder="发行总量（0=不限）" />
      <view class="mt" />
      <up-input v-model="form.per_limit" type="number" placeholder="每人限领" />
      <view class="mt" />
      <up-input v-model="form.start_at" placeholder="开始时间（ISO）" />
      <view class="mt" />
      <up-input v-model="form.end_at" placeholder="结束时间（ISO）" />
      <view class="mt" />
      <up-input v-model="form.description" placeholder="说明" />
      <view class="mt" />
      <picker mode="selector" :range="stackRuleOptions" range-key="label" @change="(e: any) => form.stack_rule = stackRuleOptions[e.detail.value].value">
        <view class="picker">叠加规则：{{ stackRuleOptions.find(o => o.value === form.stack_rule)?.label || '-' }}</view>
      </picker>
      <view class="mt" />
      <picker mode="selector" :range="targetTypeOptions" range-key="label" @change="(e: any) => form.target_type = targetTypeOptions[e.detail.value].value">
        <view class="picker">目标用户：{{ targetTypeOptions.find(o => o.value === form.target_type)?.label || '-' }}</view>
      </picker>
      <view class="mt" />
      <up-input v-if="form.target_type === 'vip_level'" v-model="form.target_value" placeholder="会员等级值" />
      <view class="row mt">
        <text>启用</text>
        <switch :checked="form.status === 1" @change="(e: any) => form.status = e.detail.value ? 1 : 0" />
      </view>
      <view class="mt-lg" />
      <up-button type="primary" @click="submit">保存</up-button>
    </view>
  </up-popup>
</template>

<style scoped>
.popup-body { padding: 24rpx; box-sizing: border-box; max-height: 80vh; overflow-y: auto; }
.popup-title { font-size: 30rpx; font-weight: 700; margin-bottom: 14rpx; }
.picker { min-height: 76rpx; border: 1px solid var(--eapp-border); border-radius: 12rpx; padding: 0 20rpx; display: flex; align-items: center; color: var(--eapp-text); }
.row { display: flex; align-items: center; justify-content: space-between; font-size: 26rpx; }
.mt { margin-top: 12rpx; }
.mt-lg { margin-top: 16rpx; }
</style>

<script setup lang="ts">
import { reactive, watch } from 'vue'
const props = defineProps<{ show: boolean; items: Array<{ id: number; title: string; price: number; qty: number }>; loading?: boolean }>()
const emit = defineEmits<{ (e: 'close'): void; (e: 'submit', payload: any): void }>()
const form = reactive<any>({ items: [], remark: '' })
watch(() => props.show, (v) => {
  if (v) {
    form.items = (props.items || []).map((it) => ({ item_id: it.id, price: Number(it.price || 0) }))
    form.remark = ''
  }
})
function onSubmit() {
  if (form.items.some((x: any) => Number(x.price) < 0)) { uni.showToast({ title: '价格不能为负', icon: 'none' }); return }
  emit('submit', { items: form.items, remark: form.remark.trim() })
}
</script>

<template>
  <up-popup :show="show" mode="bottom" round="16" @close="$emit('close')">
    <view class="popup">
      <view class="title">订单改价</view>
      <view v-for="(it, i) in props.items" :key="it.id" class="row">
        <text class="row-title">{{ it.title }}</text>
        <up-input v-model="form.items[i].price" type="digit" class="row-input" />
      </view>
      <up-input v-model="form.remark" placeholder="改价备注（必填）" class="mt" />
      <up-button type="primary" :loading="loading" class="mt-l" @click="onSubmit">保存改价</up-button>
    </view>
  </up-popup>
</template>

<style scoped>
.popup { padding: 24rpx; }
.title { font-size: 30rpx; font-weight: 700; margin-bottom: 14rpx; }
.row { display: flex; gap: 16rpx; align-items: center; padding: 10rpx 0; }
.row-title { flex: 1; font-size: 24rpx; }
.row-input { width: 220rpx; }
.mt { margin-top: 12rpx; }
.mt-l { margin-top: 20rpx; }
</style>

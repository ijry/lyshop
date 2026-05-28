<script setup lang="ts">
import { reactive, watch } from 'vue'
const props = defineProps<{ show: boolean; placeholder?: string; loading?: boolean }>()
const emit = defineEmits<{ (e: 'close'): void; (e: 'submit', content: string): void }>()
const form = reactive({ content: '' })
watch(() => props.show, (v) => { if (v) form.content = '' })
function onSubmit() {
  if (!form.content.trim()) { uni.showToast({ title: '请输入内容', icon: 'none' }); return }
  emit('submit', form.content.trim())
}
</script>

<template>
  <up-popup :show="show" mode="bottom" round="16" @close="$emit('close')">
    <view class="popup">
      <view class="title">添加备注</view>
      <up-textarea v-model="form.content" :placeholder="placeholder || '请输入备注内容'" />
      <up-button type="primary" :loading="loading" class="mt-l" @click="onSubmit">保存</up-button>
    </view>
  </up-popup>
</template>

<style scoped>
.popup { padding: 24rpx; }
.title { font-size: 30rpx; font-weight: 700; margin-bottom: 14rpx; }
.mt-l { margin-top: 20rpx; }
</style>

<template>
  <view class="min-h-screen bg-gray-50 pb-40rpx">
    <u-navbar title="订单评价" :placeholder="true" />

    <view class="p-20rpx">
      <view class="bg-white rounded-20rpx p-24rpx mb-20rpx" v-if="meta.order_no">
        <text class="text-28rpx font-600 text-gray-800 block mb-12rpx">订单 {{ meta.order_no }}</text>
        <view class="flex items-center justify-between text-24rpx text-gray-500">
          <text>物流评分</text>
          <up-rate v-model="logisticsScore" :count="5" size="20" active-color="#dc2626" />
        </view>
      </view>

      <view v-for="item in items" :key="item.order_item_id" class="bg-white rounded-20rpx p-24rpx mb-20rpx">
        <view class="flex gap-16rpx mb-16rpx">
          <image :src="item.product_cover" mode="aspectFill" style="width: 96rpx; height: 96rpx; border-radius: 12rpx;" />
          <view class="flex-1 min-w-0">
            <text class="text-28rpx font-500 text-gray-800 block truncate">{{ item.product_title }}</text>
            <text class="text-22rpx text-gray-400 block mt-6rpx">订单商品ID：{{ item.order_item_id }}</text>
          </view>
        </view>

        <view class="flex items-center justify-between mb-16rpx">
          <text class="text-24rpx text-gray-500">商品评分</text>
          <up-rate v-model="item.product_score" :count="5" size="20" active-color="#dc2626" />
        </view>

        <u-textarea v-model="item.content" placeholder="写点使用感受..." :auto-height="true" maxlength="500" />
      </view>

      <view v-if="canAppend" class="bg-white rounded-20rpx p-24rpx mb-20rpx">
        <text class="text-28rpx font-600 text-gray-800 block mb-16rpx">追加评论</text>
        <u-textarea v-model="appendContent" placeholder="追加评论内容" :auto-height="true" maxlength="500" />
      </view>

      <u-button type="primary" shape="circle" :loading="saving" text="提交评价" @click="submitReview" />
    </view>
  </view>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { get, post } from '@/utils/request'

const meta = ref<any>({})
const items = ref<any[]>([])
const logisticsScore = ref(5)
const appendContent = ref('')
const saving = ref(false)
const canAppend = ref(false)

async function loadMeta(id: number) {
  const data = await get<any>(`/api/v1/orders/${id}/review`)
  meta.value = data || {}
  items.value = Array.isArray(data?.options) ? data.options.map((item: any) => ({
    ...item,
    product_score: Number(item.product_score || 5),
    content: String(item.content || ''),
  })) : []
  logisticsScore.value = Number(data?.logistics_score || 5)
  canAppend.value = !!data?.can_append
}

async function submitReview() {
  if (saving.value) return
  saving.value = true
  try {
    const createItems = items.value.filter((item: any) => !item.has_review)
    const editItems = items.value.filter((item: any) => item.has_review)

    if (createItems.length) {
      await post(`/api/v1/orders/${meta.value.order_id}/review`, {
        mode: 'create',
        logistics_score: logisticsScore.value,
        items: createItems.map((item: any) => ({
          order_item_id: item.order_item_id,
          product_score: item.product_score,
          content: item.content,
        })),
        content: createItems[0]?.content || '',
      })
    }

    if (editItems.length) {
      await post(`/api/v1/orders/${meta.value.order_id}/review`, {
        mode: 'edit',
        logistics_score: logisticsScore.value,
        items: editItems.map((item: any) => ({
          order_item_id: item.order_item_id,
          product_score: item.product_score,
          content: item.content,
        })),
        content: editItems[0]?.content || '',
      })
    }

    if (canAppend.value && appendContent.value.trim()) {
      await post(`/api/v1/orders/${meta.value.order_id}/review`, {
        mode: 'append',
        items: editItems.map((item: any) => ({ order_item_id: item.order_item_id })),
        append_content: appendContent.value.trim(),
        content: appendContent.value.trim(),
      })
    }

    uni.showToast({ title: '提交成功', icon: 'success' })
    uni.navigateBack()
  } catch (error: any) {
    uni.showToast({ title: error?.message || '提交失败', icon: 'none' })
  } finally {
    saving.value = false
  }
}

onMounted(async () => {
  const pages = getCurrentPages()
  const query = (pages[pages.length - 1] as any).options
  await loadMeta(Number(query.id))
})
</script>

<template>
  <view class="min-h-screen bg-gray-50 pb-40rpx">
    <u-navbar :title="$t('afterSaleApply.title')" :placeholder="true" />

    <view v-if="detail.id" class="p-20rpx">
      <view class="bg-white rounded-20rpx p-24rpx mb-20rpx">
        <text class="text-28rpx font-600 text-gray-800 block">{{ $t('afterSaleApply.order') }} {{ detail.order_no }}</text>
        <text class="text-22rpx text-gray-400 block mt-8rpx">{{ $t('afterSaleApply.selectProducts') }}</text>
      </view>

      <view v-for="item in formItems" :key="item.order_item_id" class="bg-white rounded-20rpx p-24rpx mb-20rpx">
        <view class="flex gap-16rpx">
          <image :src="item.cover" mode="aspectFill" style="width: 120rpx; height: 120rpx; border-radius: 16rpx;" />
          <view class="flex-1 min-w-0">
            <text class="text-26rpx text-gray-800 block line-clamp-2">{{ item.title }}</text>
            <text class="text-22rpx text-gray-400 block mt-8rpx">{{ $t('afterSaleApply.refundableQty') }}{{ item.max_qty }}</text>
            <view class="flex items-center mt-16rpx">
              <text class="text-22rpx text-gray-500 mr-16rpx">{{ $t('afterSaleApply.applyQty') }}</text>
              <u-number-box v-model="item.qty" :min="0" :max="item.max_qty" integer />
            </view>
          </view>
        </view>
      </view>

      <view class="bg-white rounded-20rpx p-24rpx mb-20rpx">
        <text class="text-26rpx font-600 text-gray-800 block mb-16rpx">{{ $t('afterSaleApply.afterSaleType') }}</text>
        <view class="flex gap-16rpx mb-24rpx">
          <view
            class="px-24rpx py-14rpx rounded-16rpx border"
            :class="form.case_type === 'return' ? 'border-red-300 bg-red-50 text-red-600' : 'border-gray-200 text-gray-600'"
            @click="form.case_type = 'return'"
          >
            {{ $t('afterSaleApply.refund') }}
          </view>
          <view
            class="px-24rpx py-14rpx rounded-16rpx border"
            :class="form.case_type === 'exchange' ? 'border-red-300 bg-red-50 text-red-600' : 'border-gray-200 text-gray-600'"
            @click="form.case_type = 'exchange'"
          >
            {{ $t('afterSaleApply.exchange') }}
          </view>
        </view>

        <text class="text-26rpx font-600 text-gray-800 block mb-12rpx">{{ $t('afterSaleApply.reason') }}</text>
        <u-input v-model="form.reason" border="surround" :placeholder="$t('afterSaleApply.reasonPlaceholder')" />

        <text class="text-26rpx font-600 text-gray-800 block mb-12rpx mt-20rpx">{{ $t('afterSaleApply.description') }}</text>
        <u-textarea v-model="form.apply_content" :placeholder="$t('afterSaleApply.descriptionPlaceholder')" :auto-height="true" maxlength="500" />

        <text class="text-26rpx font-600 text-gray-800 block mb-12rpx mt-20rpx">{{ $t('afterSaleApply.evidenceImages') }}</text>
        <view class="flex flex-wrap gap-12rpx">
          <view
            v-for="(img, idx) in form.apply_images"
            :key="img + idx"
            class="relative"
            style="width: 150rpx; height: 150rpx;"
          >
            <image :src="img" mode="aspectFill" style="width: 150rpx; height: 150rpx; border-radius: 16rpx;" @click="previewImages(form.apply_images, idx)" />
            <view
              class="absolute right-8rpx top-8rpx w-32rpx h-32rpx rounded-full bg-black text-white flex items-center justify-center"
              style="opacity: 0.65;"
              @click.stop="removeImage(idx)"
            >
              ×
            </view>
          </view>
          <view
            v-if="form.apply_images.length < 9"
            class="w-150rpx h-150rpx rounded-16rpx border border-dashed border-gray-300 flex items-center justify-center text-24rpx text-gray-400 bg-gray-50"
            @click="chooseImages"
          >
            {{ $t('afterSaleApply.addImage') }}
          </view>
        </view>
      </view>

      <u-button type="primary" shape="circle" :loading="submitting" :text="submitting ? $t('afterSaleApply.submitting') : $t('afterSaleApply.submit')" @click="submit" />
    </view>
  </view>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { get, post, upload } from '@/utils/request'

const { t } = useI18n()

const detail = ref<any>({})
const submitting = ref(false)
const formItems = ref<any[]>([])
const orderID = ref(0)
const form = ref<any>({
  case_type: 'return',
  reason: '',
  apply_content: '',
  apply_images: [],
})

function toast(msg: string) {
  uni.showToast({ title: msg, icon: 'none' })
}

function readOrderID() {
  const pages = getCurrentPages()
  const query = (pages[pages.length - 1] as any).options
  return Number(query.id || 0)
}

function removeImage(index: number) {
  form.value.apply_images.splice(index, 1)
}

function previewImages(urls: string[], index: number) {
  if (!urls.length) return
  uni.previewImage({ urls, current: urls[index] || urls[0] })
}

function chooseFiles(limit: number): Promise<string[]> {
  return new Promise((resolve, reject) => {
    uni.chooseImage({
      count: limit,
      sizeType: ['compressed'],
      sourceType: ['album', 'camera'],
      success(res) {
        resolve(res.tempFilePaths || [])
      },
      fail(err) {
        reject(err)
      },
    })
  })
}

async function chooseImages() {
  const remain = Math.max(0, 9 - form.value.apply_images.length)
  if (!remain) {
    toast(t('afterSaleApply.maxImages'))
    return
  }
  try {
    const files = await chooseFiles(remain)
    for (const filePath of files) {
      const result: any = await upload('/api/v1/upload', filePath)
      if (result?.url) {
        form.value.apply_images.push(String(result.url))
      }
      if (form.value.apply_images.length >= 9) break
    }
  } catch {
    toast(t('afterSaleApply.imageUploadFailed'))
  }
}

function normalizeQty(raw: any, maxQty: number) {
  const num = Number(raw || 0)
  if (Number.isNaN(num) || num < 0) return 0
  return Math.min(Math.floor(num), maxQty)
}

async function submit() {
  if (submitting.value) return
  const reason = String(form.value.reason || '').trim()
  if (!reason) {
    toast(t('afterSaleApply.reasonRequired'))
    return
  }
  const items = formItems.value
    .map((item: any) => ({
      order_item_id: Number(item.order_item_id),
      qty: normalizeQty(item.qty, Number(item.max_qty || 0)),
    }))
    .filter((item: any) => item.qty > 0)
  if (!items.length) {
    toast(t('afterSaleApply.selectAtLeastOne'))
    return
  }
  submitting.value = true
  try {
    const result: any = await post(`/api/v1/orders/${orderID.value}/after-sales`, {
      case_type: form.value.case_type,
      reason,
      apply_content: String(form.value.apply_content || ''),
      apply_images: form.value.apply_images.slice(),
      items,
    })
    const caseID = Number(result?.id || 0)
    uni.showToast({ title: t('afterSaleApply.submitted'), icon: 'success' })
    setTimeout(() => {
      if (caseID > 0) {
        uni.redirectTo({ url: `/pages/order/after-sale-detail?id=${caseID}` })
      } else {
        uni.redirectTo({ url: `/pages/order/detail?id=${orderID.value}` })
      }
    }, 400)
  } catch (error: any) {
    toast(error?.message || t('afterSaleApply.submitFailed'))
  } finally {
    submitting.value = false
  }
}

onMounted(async () => {
  orderID.value = readOrderID()
  detail.value = await get<any>(`/api/v1/orders/${orderID.value}`)
  formItems.value = (detail.value?.items || []).map((item: any) => ({
    order_item_id: Number(item.id),
    title: String(item.title || ''),
    cover: String(item.cover || ''),
    max_qty: Math.max(1, Number(item.qty || 1)),
    qty: 0,
  }))
})
</script>

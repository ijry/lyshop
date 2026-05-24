<template>
  <view class="min-h-screen bg-gray-50 pb-40rpx">
    <u-navbar title="订单评价" :placeholder="true" />

    <view class="p-20rpx">
      <view class="bg-white rounded-20rpx p-24rpx mb-20rpx" v-if="meta.order_no">
        <text class="text-28rpx font-600 text-gray-800 block mb-12rpx">订单 {{ meta.order_no }}</text>
        <view class="flex items-center justify-between text-24rpx text-gray-500" v-if="viewMode === 'root'">
          <text>物流评分</text>
          <up-rate v-model="logisticsScore" :count="5" size="20" active-color="#dc2626" />
        </view>
        <text class="text-22rpx text-gray-400 block mt-12rpx" v-if="viewMode === 'root'">根评价提交后，才能继续追加评论。</text>
        <text class="text-22rpx text-gray-400 block mt-12rpx" v-else>你正在追加已发布评价。</text>
      </view>

      <view v-if="viewMode === 'root'" v-for="item in items" :key="item.order_item_id" class="bg-white rounded-20rpx p-24rpx mb-20rpx">
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

        <view class="mt-16rpx">
          <view class="flex items-center justify-between mb-12rpx">
            <text class="text-24rpx text-gray-500">图片</text>
            <text class="text-22rpx text-gray-400">{{ item.images.length }}/9</text>
          </view>
          <view class="flex flex-wrap gap-12rpx">
            <view
              v-for="(img, idx) in item.images"
              :key="img + idx"
              class="relative"
              style="width: 150rpx; height: 150rpx;"
            >
              <image
                :src="img"
                mode="aspectFill"
                style="width: 150rpx; height: 150rpx; border-radius: 16rpx;"
                @click="previewImages(item.images, idx)"
              />
              <view
                class="absolute right-8rpx top-8rpx w-32rpx h-32rpx rounded-full bg-black text-white flex items-center justify-center"
                style="opacity: 0.65;"
                @click.stop="removeImage(item.images, idx)"
              >
                ×
              </view>
            </view>
            <view
              v-if="item.images.length < 9"
              class="w-150rpx h-150rpx rounded-16rpx border border-dashed border-gray-300 flex items-center justify-center text-24rpx text-gray-400 bg-gray-50"
              @click="chooseImages(item.images)"
            >
              + 添加
            </view>
          </view>
        </view>
      </view>

      <view v-if="viewMode === 'append' && canAppend" class="bg-white rounded-20rpx p-24rpx mb-20rpx">
        <text class="text-28rpx font-600 text-gray-800 block mb-16rpx">追加评论</text>
        <u-textarea v-model="appendContent" placeholder="追加评论内容" :auto-height="true" maxlength="500" />

        <view class="mt-16rpx">
          <view class="flex items-center justify-between mb-12rpx">
            <text class="text-24rpx text-gray-500">追评图片</text>
            <text class="text-22rpx text-gray-400">{{ appendImages.length }}/9</text>
          </view>
          <view class="flex flex-wrap gap-12rpx">
            <view
              v-for="(img, idx) in appendImages"
              :key="img + idx"
              class="relative"
              style="width: 150rpx; height: 150rpx;"
            >
              <image
                :src="img"
                mode="aspectFill"
                style="width: 150rpx; height: 150rpx; border-radius: 16rpx;"
                @click="previewImages(appendImages, idx)"
              />
              <view
                class="absolute right-8rpx top-8rpx w-32rpx h-32rpx rounded-full bg-black text-white flex items-center justify-center"
                style="opacity: 0.65;"
                @click.stop="removeImage(appendImages, idx)"
              >
                ×
              </view>
            </view>
            <view
              v-if="appendImages.length < 9"
              class="w-150rpx h-150rpx rounded-16rpx border border-dashed border-gray-300 flex items-center justify-center text-24rpx text-gray-400 bg-gray-50"
              @click="chooseImages(appendImages)"
            >
              + 添加
            </view>
          </view>
        </view>
      </view>

      <view v-else class="bg-white rounded-20rpx p-24rpx mb-20rpx">
        <text class="text-28rpx font-600 text-gray-800 block mb-8rpx">追加评论</text>
        <text class="text-24rpx text-gray-400">完成根评价后，才能追加评论。</text>
      </view>

      <u-button
        v-if="viewMode === 'root'"
        type="primary"
        shape="circle"
        :loading="savingRoot"
        :text="rootButtonText"
        @click="submitRootReview"
        class="mb-16rpx"
      />
      <u-button
        v-if="viewMode === 'append' && canAppend"
        type="success"
        shape="circle"
        :loading="savingAppend"
        text="提交追加评价"
        @click="submitAppendReview"
      />
    </view>
  </view>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { get, post, upload } from '@/utils/request'

const meta = ref<any>({})
const items = ref<any[]>([])
const logisticsScore = ref(5)
const appendContent = ref('')
const appendImages = ref<string[]>([])
const savingRoot = ref(false)
const savingAppend = ref(false)
const canAppend = ref(false)
const orderID = ref(0)
const targetItemID = ref(0)
const viewMode = ref<'root' | 'append'>('root')

const rootButtonText = computed(() => {
  return items.value.some((item: any) => !item.has_review) ? '提交根评价' : '更新根评价'
})

function toast(message: string) {
  uni.showToast({ title: message, icon: 'none' })
}

function readQueryID() {
  const pages = getCurrentPages()
  const query = (pages[pages.length - 1] as any).options
  const mode = String(query.mode || '').trim()
  viewMode.value = mode === 'append' ? 'append' : 'root'
  targetItemID.value = Number(query.item_id || query.order_item_id || query.review_id || 0)
  return Number(query.id || 0)
}

async function loadMeta(id: number) {
  const data = await get<any>(`/api/v1/orders/${id}/review`)
  meta.value = data || {}
  const options = Array.isArray(data?.options)
    ? data.options.map((item: any) => ({
        ...item,
        product_score: Number(item.product_score || 5),
        content: String(item.content || ''),
        images: Array.isArray(item.images) ? item.images.map((img: any) => String(img || '')) : [],
      }))
    : []
  logisticsScore.value = Number(data?.logistics_score || 5)
  if (viewMode.value === 'append') {
    const target = options.find((item: any) => {
      if (targetItemID.value > 0) return Number(item.order_item_id || 0) === targetItemID.value && !!item.has_review
      return !!item.has_review
    })
    items.value = target ? [target] : []
    canAppend.value = !!target
    logisticsScore.value = Number(target?.logistics_score || logisticsScore.value || 5)
    return
  }
  items.value = options
  canAppend.value = !!data?.can_append
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

async function uploadImages(target: string[]) {
  const remain = Math.max(0, 9 - target.length)
  if (remain <= 0) {
    toast('最多上传 9 张图片')
    return
  }
  const files = await chooseFiles(remain)
  if (!files.length) return
  for (const filePath of files) {
    const result: any = await upload('/api/v1/upload', filePath)
    if (result?.url) {
      target.push(String(result.url))
    }
    if (target.length >= 9) break
  }
}

async function chooseImages(target: string[]) {
  try {
    await uploadImages(target)
  } catch {
    toast('图片上传失败')
  }
}

function removeImage(target: string[], index: number) {
  target.splice(index, 1)
}

function previewImages(urls: string[], index: number) {
  if (!urls.length) return
  uni.previewImage({ urls, current: urls[index] || urls[0] })
}

async function submitRootReview() {
  if (viewMode.value !== 'root') return
  if (savingRoot.value) return
  const createItems = items.value.filter((item: any) => !item.has_review)
  const editItems = items.value.filter((item: any) => item.has_review)
  if (!createItems.length && !editItems.length) {
    toast('暂无可提交的根评价')
    return
  }
  savingRoot.value = true
  try {
    if (createItems.length) {
      await post(`/api/v1/orders/${meta.value.order_id}/review`, {
        mode: 'create',
        logistics_score: logisticsScore.value,
        items: createItems.map((item: any) => ({
          order_item_id: item.order_item_id,
          product_score: item.product_score,
          content: item.content,
          images: item.images || [],
        })),
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
          images: item.images || [],
        })),
      })
    }
    await loadMeta(orderID.value)
    toast('根评价已提交')
  } catch (error: any) {
    toast(error?.message || '提交失败')
  } finally {
    savingRoot.value = false
  }
}

async function submitAppendReview() {
  if (viewMode.value !== 'append') return
  if (savingAppend.value) return
  if (!canAppend.value) {
    toast('请先完成根评价')
    return
  }
  const content = appendContent.value.trim()
  const itemsToAppend = items.value.length ? [items.value[0]] : []
  if (!content && appendImages.value.length === 0) {
    toast('请填写追评内容或上传图片')
    return
  }
  if (!itemsToAppend.length) {
    toast('没有可追加的评价')
    return
  }
  savingAppend.value = true
  try {
    await post(`/api/v1/orders/${meta.value.order_id}/review`, {
      mode: 'append',
      items: itemsToAppend.map((item: any) => ({ order_item_id: item.order_item_id })),
      append_content: content,
      append_images: appendImages.value.slice(),
    })
    appendContent.value = ''
    appendImages.value = []
    await loadMeta(orderID.value)
    toast('追评已提交')
  } catch (error: any) {
    toast(error?.message || '提交失败')
  } finally {
    savingAppend.value = false
  }
}

onMounted(async () => {
  orderID.value = readQueryID()
  await loadMeta(orderID.value)
})
</script>

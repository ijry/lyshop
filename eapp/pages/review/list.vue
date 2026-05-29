<script setup lang="ts">
import { onLoad, onShow, onPullDownRefresh, onReachBottom } from '@dcloudio/uni-app'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import ReviewCard from '@/components/biz/ReviewCard.vue'
import EmptyState from '@/components/biz/EmptyState.vue'
import { useReviewList } from '@/composables/useReviewList'
import { getReviewDetail, replyReview } from '@/api/review'

const { t } = useI18n()

const tabs = [
  { name: '全部', status: 'all' as const },
  { name: '待回复', status: 'pending' as const },
  { name: '已回复', status: 'replied' as const },
]
const current = ref(0)
const keyword = ref('')

const { list, loading, total, load, refresh, loadMore, applyFilter } = useReviewList()

function onTab(idx: number) {
  current.value = idx
  applyFilter({ reply_status: tabs[idx].status })
  load()
}

function onSearch() {
  applyFilter({ keyword: keyword.value.trim() })
  load()
}

/* ---------- detail popup ---------- */
const showDetail = ref(false)
const detailLoading = ref(false)
const detail = ref<any>(null)

async function openDetail(id: number) {
  showDetail.value = true
  detailLoading.value = true
  try {
    detail.value = await getReviewDetail(id)
  } finally {
    detailLoading.value = false
  }
}

function onPreviewImage(url: string, urls: string[]) {
  uni.previewImage({ current: url, urls })
}

/* ---------- reply popup ---------- */
const showReply = ref(false)
const replyContent = ref('')
const replyLoading = ref(false)
const replyTargetId = ref(0)

function openReply(id: number) {
  replyTargetId.value = id
  replyContent.value = ''
  showReply.value = true
}

async function submitReply() {
  const content = replyContent.value.trim()
  if (!content) {
    uni.showToast({ title: '请输入回复内容', icon: 'none' })
    return
  }
  replyLoading.value = true
  try {
    await replyReview(replyTargetId.value, content)
    uni.showToast({ title: t('review.replySuccess'), icon: 'success' })
    showReply.value = false
    showDetail.value = false
    await load()
  } finally {
    replyLoading.value = false
  }
}

onLoad(() => load())
onShow(() => load())
onPullDownRefresh(async () => { await refresh(); uni.stopPullDownRefresh() })
onReachBottom(() => loadMore())
</script>

<template>
  <view class="page">
    <!-- Tabs -->
    <up-tabs
      :list="tabs"
      :current="current"
      :scrollable="true"
      keyName="name"
      @click="(item) => onTab(item.index)"
      :activeStyle="{ color: '#fff', backgroundColor: 'var(--eapp-primary)', borderRadius: '999rpx', height: '56rpx', lineHeight: '56rpx', padding: '0 24rpx' }"
      :inactiveStyle="{ color: 'var(--eapp-text-muted)', backgroundColor: 'var(--eapp-bg)', borderRadius: '999rpx', height: '56rpx', lineHeight: '56rpx', padding: '0 24rpx' }"
      :itemStyle="{ padding: '0 4rpx', height: '80rpx' }"
      lineColor="transparent"
    />

    <!-- Search -->
    <view class="search-bar">
      <up-input
        v-model="keyword"
        placeholder="搜索评价内容"
        clearable
        @confirm="onSearch"
      />
    </view>

    <!-- List -->
    <view class="list">
      <EmptyState v-if="!loading && !list.length" :title="t('review.noReviews')" />
      <ReviewCard
        v-for="row in list"
        :key="row.id"
        :row="row"
        @click="openDetail(row.id)"
      />
      <view v-if="loading" class="loading-hint">{{ t('common.loading') }}</view>
      <view v-else-if="list.length >= total && list.length > 0" class="loading-hint">已加载全部</view>
    </view>

    <!-- Detail Popup -->
    <up-popup :show="showDetail" mode="bottom" round="16" @close="showDetail = false">
      <view class="popup-body detail-popup">
        <view v-if="detailLoading" class="loading-hint">{{ t('common.loading') }}</view>
        <template v-else-if="detail">
          <!-- Product info -->
          <view class="detail-product">
            <image
              v-if="detail.product?.cover"
              :src="detail.product.cover"
              mode="aspectFill"
              class="detail-cover"
            />
            <view class="detail-product-info">
              <text class="detail-product-title">{{ detail.product?.title || '-' }}</text>
              <text class="detail-scores">
                商品 {{ '★' }}{{ Number(detail.product_score || 0).toFixed(1) }}
                / 物流 {{ '★' }}{{ Number(detail.logistics_score || 0).toFixed(1) }}
              </text>
            </view>
          </view>

          <!-- Content -->
          <text class="detail-content">{{ detail.content || '无评价内容' }}</text>

          <!-- Images -->
          <view v-if="detail.images && detail.images.length" class="detail-images">
            <image
              v-for="(img, idx) in detail.images"
              :key="idx"
              :src="img"
              mode="aspectFill"
              class="detail-img"
              @click="onPreviewImage(img, detail.images)"
            />
          </view>

          <!-- Appends -->
          <view v-if="detail.appends && detail.appends.length" class="detail-section">
            <text class="detail-section-title">{{ t('review.appends') }}</text>
            <view v-for="ap in detail.appends" :key="ap.id" class="append-row">
              <text class="append-content">{{ ap.content }}</text>
              <view v-if="ap.images && ap.images.length" class="detail-images">
                <image
                  v-for="(img, idx) in ap.images"
                  :key="idx"
                  :src="img"
                  mode="aspectFill"
                  class="detail-img"
                  @click="onPreviewImage(img, ap.images)"
                />
              </view>
            </view>
          </view>

          <!-- Existing Reply -->
          <view v-if="detail.admin_reply" class="detail-section">
            <text class="detail-section-title">商家回复</text>
            <view class="reply-box">
              <text class="reply-text">{{ detail.admin_reply.content }}</text>
            </view>
          </view>

          <!-- Reply button -->
          <up-button type="primary" @click="openReply(detail.id)">
            {{ detail.admin_reply ? '修改回复' : t('review.reply') }}
          </up-button>
        </template>
      </view>
    </up-popup>

    <!-- Reply Popup -->
    <up-popup :show="showReply" mode="bottom" round="16" @close="showReply = false">
      <view class="popup-body">
        <view class="popup-title">{{ t('review.reply') }}</view>
        <up-input
          v-model="replyContent"
          type="textarea"
          :placeholder="t('review.replyPlaceholder')"
          :maxlength="500"
          height="200"
        />
        <view class="mt-16rpx" />
        <up-button type="primary" :loading="replyLoading" @click="submitReply">
          {{ t('common.confirm') }}
        </up-button>
      </view>
    </up-popup>
  </view>
</template>

<style scoped>
.page { min-height: 100vh; background: var(--eapp-bg); }
.search-bar { padding: 16rpx 20rpx; }
.list { padding: 0 20rpx 20rpx; display: grid; gap: 14rpx; }
.loading-hint { text-align: center; padding: 24rpx 0; font-size: 24rpx; color: var(--eapp-text-muted); }

/* Detail popup */
.popup-body { padding: 24rpx; box-sizing: border-box; max-height: 80vh; overflow-y: auto; }
.popup-title { font-size: 30rpx; font-weight: 700; margin-bottom: 14rpx; }
.detail-popup { padding-bottom: 40rpx; }
.detail-product { display: flex; align-items: center; gap: 12rpx; margin-bottom: 16rpx; }
.detail-cover { width: 96rpx; height: 96rpx; border-radius: 14rpx; flex-shrink: 0; }
.detail-product-info { flex: 1; min-width: 0; display: flex; flex-direction: column; gap: 6rpx; }
.detail-product-title { font-size: 28rpx; font-weight: 600; color: var(--eapp-text); }
.detail-scores { font-size: 24rpx; color: #f59e0b; }
.detail-content { font-size: 26rpx; color: var(--eapp-text); line-height: 1.7; display: block; margin-bottom: 12rpx; }
.detail-images { display: flex; gap: 10rpx; flex-wrap: wrap; margin-bottom: 12rpx; }
.detail-img { width: 140rpx; height: 140rpx; border-radius: 12rpx; }
.detail-section { margin-bottom: 16rpx; }
.detail-section-title { font-size: 26rpx; font-weight: 600; color: var(--eapp-text); margin-bottom: 8rpx; display: block; }
.append-row { padding: 12rpx; background: var(--eapp-bg); border-radius: 12rpx; margin-bottom: 8rpx; }
.append-content { font-size: 24rpx; color: var(--eapp-text-muted); display: block; }
.reply-box { background: var(--eapp-bg); border-radius: 12rpx; padding: 14rpx; margin-bottom: 8rpx; }
.reply-text { font-size: 24rpx; color: var(--eapp-text); }
</style>

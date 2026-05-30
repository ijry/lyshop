<template>
  <view class="page">
    <view v-if="product" class="content">
      <swiper v-if="images.length" class="swiper" indicator-dots>
        <swiper-item v-for="(img, idx) in images" :key="idx">
          <image :src="img" class="swiper-image" mode="aspectFill" />
        </swiper-item>
      </swiper>
      <view v-else class="cover-placeholder">
        <image v-if="product.cover" :src="product.cover" class="cover-image" mode="aspectFill" />
        <view v-else class="no-image">暂无图片</view>
      </view>

      <view class="info-card">
        <view class="product-title">{{ product.title }}</view>
        <view class="product-price">{{ product.points_price }} 积分</view>
        <view class="product-meta">
          <text>库存: {{ product.stock === 0 ? '无限' : product.stock }}</text>
          <text>已兑: {{ product.sold_count }}</text>
        </view>
        <view v-if="product.limit_per_user > 0" class="limit-tip">每人限兑 {{ product.limit_per_user }} 件</view>
        <view v-if="product.limit_per_day > 0" class="limit-tip">每日限兑 {{ product.limit_per_day }} 件</view>
      </view>

      <view class="desc-card">
        <view class="card-title">商品详情</view>
        <view class="desc-content">{{ product.description || '暂无描述' }}</view>
      </view>
    </view>

    <view class="footer">
      <view class="my-points">我的积分: {{ userPoints }}</view>
      <button class="exchange-btn" @click="exchange">立即兑换</button>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { get, post } from '@/utils/request'

const productId = ref(0)
const product = ref<any>(null)
const userPoints = ref(0)

const images = computed(() => {
  if (!product.value?.images) return []
  try {
    return JSON.parse(product.value.images)
  } catch {
    return []
  }
})

onLoad((options: any) => {
  productId.value = Number(options.id)
  loadProduct()
  loadUserPoints()
})

async function loadProduct() {
  try {
    const data: any = await get(`/api/v1/points/products/${productId.value}`)
    product.value = data
  } catch (e: any) {
    uni.showToast({ title: e.message || '加载失败', icon: 'none' })
  }
}

async function loadUserPoints() {
  try {
    const data: any = await get('/api/v1/points/balance')
    userPoints.value = data.points || 0
  } catch (e) {
    console.error(e)
  }
}

async function exchange() {
  if (!product.value) return

  if (userPoints.value < product.value.points_price) {
    uni.showToast({ title: '积分不足', icon: 'none' })
    return
  }

  uni.showModal({
    title: '确认兑换',
    content: `确定使用 ${product.value.points_price} 积分兑换该商品吗？`,
    success: async (res) => {
      if (res.confirm) {
        try {
          await post(`/api/v1/points/products/${productId.value}/exchange`, { qty: 1 })
          uni.showToast({ title: '兑换成功', icon: 'success' })
          setTimeout(() => {
            uni.navigateTo({ url: '/pages/points/my-exchanges' })
          }, 1500)
        } catch (e: any) {
          uni.showToast({ title: e.message || '兑换失败', icon: 'none' })
        }
      }
    }
  })
}
</script>

<style scoped>
.page { min-height: 100vh; background: #f5f5f5; padding-bottom: 120rpx; }
.swiper { width: 100%; height: 600rpx; }
.swiper-image { width: 100%; height: 100%; }
.cover-placeholder { width: 100%; height: 600rpx; background: #fff; }
.cover-image { width: 100%; height: 100%; }
.no-image { width: 100%; height: 100%; display: flex; align-items: center; justify-content: center; color: #999; }
.info-card { background: #fff; padding: 30rpx; margin-bottom: 20rpx; }
.product-title { font-size: 32rpx; font-weight: bold; color: #333; margin-bottom: 20rpx; }
.product-price { font-size: 40rpx; color: #667eea; font-weight: bold; margin-bottom: 20rpx; }
.product-meta { display: flex; gap: 40rpx; font-size: 24rpx; color: #999; margin-bottom: 10rpx; }
.limit-tip { font-size: 24rpx; color: #ff6b6b; margin-top: 10rpx; }
.desc-card { background: #fff; padding: 30rpx; }
.card-title { font-size: 28rpx; font-weight: bold; color: #333; margin-bottom: 20rpx; }
.desc-content { font-size: 26rpx; color: #666; line-height: 1.6; white-space: pre-wrap; }
.footer { position: fixed; bottom: 0; left: 0; right: 0; background: #fff; padding: 20rpx 30rpx; display: flex; align-items: center; justify-content: space-between; box-shadow: 0 -2rpx 10rpx rgba(0,0,0,0.05); }
.my-points { font-size: 26rpx; color: #666; }
.exchange-btn { flex: 1; margin-left: 30rpx; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: #fff; border: none; border-radius: 50rpx; height: 80rpx; line-height: 80rpx; font-size: 28rpx; }
</style>

<template>
  <view>
    <template v-for="comp in components" :key="comp.id">
      <view v-if="comp.type === 'banner'" class="mb-16rpx">
        <u-swiper
          :list="bannerList(comp)"
          keyName="url"
          :height="comp.props?.height || 340"
          :indicator="true"
          indicatorMode="dot"
          :circular="true"
          :autoplay="true"
          radius="0"
          @click="(index) => onBannerClick(comp, index)"
        />
      </view>

      <view v-else-if="comp.type === 'category_nav'" class="bg-white py-28rpx mb-16rpx">
        <scroll-view scroll-x :show-scrollbar="false" class="px-16rpx">
          <view style="display: flex; gap: 12px;">
            <view
              v-for="item in (comp.props?.items || [])"
              :key="item.title"
              @click="openLink(item.link)"
              style="display: flex; flex-direction: column; align-items: center; gap: 8px; flex-shrink: 0; width: 65px;"
            >
              <view
                v-if="item.icon"
                style="width: 48px; height: 48px; border-radius: 14px; overflow: hidden;"
              >
                <image :src="item.icon" mode="aspectFill" style="width: 48px; height: 48px;" />
              </view>
              <view
                v-else
                style="width: 48px; height: 48px; border-radius: 14px; display: flex; align-items: center; justify-content: center; background: #fef2f2;"
              >
                <text style="color: #dc2626; font-size: 13px; font-weight: 600;">{{ item.title?.slice(0, 2) }}</text>
              </view>
              <text style="font-size: 12px; color: #374151;">{{ item.title }}</text>
            </view>
          </view>
        </scroll-view>
      </view>

      <view v-else-if="comp.type === 'notice'" class="mx-20rpx mb-16rpx">
        <u-column-notice
          v-if="noticeTexts(comp).length"
          :text="noticeTexts(comp)"
          :mode="comp.props?.mode || 'link'"
          :color="comp.props?.color || '#f97316'"
          :bgColor="comp.props?.bgColor || '#fff7ed'"
          :duration="comp.props?.duration || 2500"
          :fontSize="comp.props?.fontSize || 14"
          @click="(index) => onNoticeClick(comp, index)"
        >
          <template #icon>
            <u-icon name="volume" :color="comp.props?.color || '#f97316'" size="18" />
          </template>
        </u-column-notice>
      </view>

      <view v-else-if="comp.type === 'product_grid'" class="px-16rpx mb-16rpx">
        <view
          v-if="comp.props?.title"
          style="display: flex; align-items: center; justify-content: space-between; padding: 16px 4px 12px;"
        >
          <text style="font-size: 16px; font-weight: 700; color: #111827;">{{ comp.props.title }}</text>
          <text style="font-size: 12px; color: #9ca3af;">查看全部 ></text>
        </view>
        <u-waterfall
          :key="`${comp.id}-${waterfallVersions[comp.id] || 0}`"
          :modelValue="gridProducts[comp.id] || []"
          :addTime="30"
        >
          <template #left="{ leftList }">
            <view
              v-for="product in leftList"
              :key="`left-${product.id}`"
              style="padding: 0 6rpx 12rpx;"
            >
              <view
                @click="openProduct(product.id)"
                style="background: #fff; border-radius: 12px; overflow: hidden; box-shadow: 0 1px 8px rgba(0,0,0,0.05);"
              >
                <image
                  :src="product.cover || ''"
                  mode="aspectFill"
                  style="width: 100%; height: 170px; display: block;"
                />
                <view style="padding: 10px 12px 14px;">
                  <text
                    style="font-size: 13px; color: #1f2937; font-weight: 500; display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden; line-height: 1.4;"
                  >
                    {{ product.title }}
                  </text>
                  <view style="display: flex; align-items: center; justify-content: space-between; margin-top: 8px;">
                    <text style="font-size: 16px; color: #dc2626; font-weight: 700;">¥{{ product.price }}</text>
                    <text style="font-size: 10px; color: #9ca3af;">{{ product.sales || 0 }}付款</text>
                  </view>
                </view>
              </view>
            </view>
          </template>
          <template #right="{ rightList }">
            <view
              v-for="product in rightList"
              :key="`right-${product.id}`"
              style="padding: 0 6rpx 12rpx;"
            >
              <view
                @click="openProduct(product.id)"
                style="background: #fff; border-radius: 12px; overflow: hidden; box-shadow: 0 1px 8px rgba(0,0,0,0.05);"
              >
                <image
                  :src="product.cover || ''"
                  mode="aspectFill"
                  style="width: 100%; height: 170px; display: block;"
                />
                <view style="padding: 10px 12px 14px;">
                  <text
                    style="font-size: 13px; color: #1f2937; font-weight: 500; display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden; line-height: 1.4;"
                  >
                    {{ product.title }}
                  </text>
                  <view style="display: flex; align-items: center; justify-content: space-between; margin-top: 8px;">
                    <text style="font-size: 16px; color: #dc2626; font-weight: 700;">¥{{ product.price }}</text>
                    <text style="font-size: 10px; color: #9ca3af;">{{ product.sales || 0 }}付款</text>
                  </view>
                </view>
              </view>
            </view>
          </template>
        </u-waterfall>
      </view>

      <view v-else-if="comp.type === 'image_ad'" class="px-20rpx mb-16rpx">
        <image
          v-if="comp.props?.url"
          :src="comp.props.url"
          mode="widthFix"
          style="width: 100%; border-radius: 12px; display: block;"
          @click="openLink(comp.props.link)"
        />
      </view>

      <!-- Grid 宫格 -->
      <view v-else-if="comp.type === 'grid'" class="bg-white mb-16rpx py-20rpx">
        <view :style="{ display: 'flex', flexWrap: 'wrap' }">
          <view v-for="item in (comp.props?.items || [])" :key="item.title"
            @click="openLink(item.link)"
            :style="{ width: (100 / (comp.props?.columns || 4)) + '%', display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '6px', padding: '10px 0', cursor: 'pointer' }">
            <view :style="{ width: '44px', height: '44px', borderRadius: '12px', background: item.bg || '#f5f5f5', display: 'flex', alignItems: 'center', justifyContent: 'center', fontSize: '20px' }">
              {{ item.icon }}
            </view>
            <text style="font-size: 12px; color: #333;">{{ item.title }}</text>
          </view>
        </view>
      </view>

      <view v-else-if="comp.type === 'rich_text'" class="px-30rpx py-16rpx">
        <rich-text :nodes="comp.props?.content || ''" />
      </view>

      <view
        v-else-if="comp.type === 'spacer'"
        :style="{ height: (comp.props?.height || 16) + 'rpx', background: comp.props?.background || '#f5f5f5' }"
      />

      <view v-else-if="comp.type === 'marketing_zone'" class="px-20rpx mb-16rpx">
        <view style="border-radius: 12px; padding: 16px 20px; background: linear-gradient(135deg, #dc2626 0%, #ea580c 100%);">
          <view style="display: flex; align-items: center; justify-content: space-between;">
            <view>
              <text style="color: #fff; font-size: 17px; font-weight: 700;">限时秒杀</text>
              <text style="color: rgba(255,255,255,0.8); font-size: 11px; margin-left: 8px;">抢购进行中</text>
            </view>
            <text style="color: rgba(255,255,255,0.8); font-size: 12px;">更多 ></text>
          </view>
        </view>
      </view>
    </template>
  </view>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { get } from '@/utils/request'

const props = defineProps<{ components: any[] }>()
const gridProducts = ref<Record<string, any[]>>({})
const waterfallVersions = ref<Record<string, number>>({})

function bannerList(comp: any) {
  const images = comp.props?.images || []
  return images
    .map((img: any) => {
      if (typeof img === 'string') return { url: img }
      if (!img?.url) return null
      return {
        url: img.url,
        title: img.title || '',
        link: img.link || '',
      }
    })
    .filter(Boolean)
}

function normalizeMarketingLink(link: string) {
  if (!link) return ''
  if (!link.startsWith('/pages/marketing/coupon')) return link
  if (link.includes('mode=')) return link
  return link.includes('?') ? `${link}&mode=claim` : `${link}?mode=claim`
}

function openLink(link: string) {
  const target = normalizeMarketingLink(link || '')
  if (!target) return
  const [path] = target.split('?')
  const tabPages = [
    '/pages/index/index',
    '/pages/product/list',
    '/pages/cart/index',
    '/pages/order/list',
    '/pages/user/index',
  ]
  if (tabPages.includes(path)) {
    uni.switchTab({ url: path })
    return
  }
  uni.navigateTo({ url: target })
}

function onBannerClick(comp: any, index: number) {
  const item = bannerList(comp)[index]
  openLink(item?.link || '')
}

function noticeTexts(comp: any) {
  const items = comp.props?.items
  if (Array.isArray(items) && items.length) {
    return items
      .map((item: any) => (typeof item === 'string' ? item : item?.text || ''))
      .filter(Boolean)
  }
  return comp.props?.text ? [comp.props.text] : []
}

function onNoticeClick(comp: any, index: number) {
  const items = comp.props?.items
  if (Array.isArray(items) && items.length) {
    const item = items[index]
    const link = typeof item === 'object' ? item?.link : ''
    openLink(link)
    return
  }
  const links = comp.props?.links
  if (Array.isArray(links) && links[index]) openLink(links[index])
}

function openProduct(id: string | number) {
  uni.navigateTo({ url: `/pages/product/detail?id=${id}` })
}

async function loadGridProducts() {
  const nextProducts: Record<string, any[]> = {}
  const nextVersions = { ...waterfallVersions.value }

  await Promise.all(
    props.components
      .filter((comp) => comp.type === 'product_grid')
      .map(async (comp) => {
        const params: any = { page: 1, size: comp.props?.limit || 10 }
        const data = await get<any>('/api/v1/products', params)
        nextProducts[comp.id] = data?.list || []
        nextVersions[comp.id] = (nextVersions[comp.id] || 0) + 1
      })
  )

  gridProducts.value = nextProducts
  waterfallVersions.value = nextVersions
}

onMounted(loadGridProducts)
watch(() => props.components, loadGridProducts, { deep: true })
</script>

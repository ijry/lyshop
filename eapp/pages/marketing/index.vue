<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import SectionHeader from '@/components/biz/SectionHeader.vue'
import EntryList, { type EntryItem } from '@/components/biz/EntryList.vue'
import AppIcon from '@/components/common/AppIcon.vue'

const { t } = useI18n()

const promotions: EntryItem[] = [
  { key: 'coupon', title: t('marketing.coupon'), desc: '满减 / 折扣券与定向发券', iconName: 'coupon', path: '/pages/marketing/coupon' },
  { key: 'seckill', title: t('marketing.seckill'), desc: '限时秒杀活动与活动商品', iconName: 'flash', path: '/pages/marketing/seckill' },
  { key: 'groupBuy', title: t('marketing.groupBuy'), desc: '拼团活动配置与成团管理', iconName: 'users', path: '/pages/marketing/group-buy' },
  { key: 'bargain', title: t('marketing.bargain'), desc: '砍价活动与底价规则', iconName: 'money', path: '/pages/marketing/bargain' },
]

const operations: EntryItem[] = [
  { key: 'vip', title: t('marketing.vip'), desc: '会员套餐 / 等级 / 专属价', iconName: 'crown', path: '/pages/vip/plan-list' },
  { key: 'distribution', title: t('marketing.distribution'), desc: '分销商、佣金与比例配置', iconName: 'share', path: '/pages/distribution/index' },
  { key: 'pointsMall', title: t('marketing.pointsMall'), desc: '积分商品与兑换记录', iconName: 'diamond', path: '/pages/points/index' },
  { key: 'messageBroadcast', title: '消息群发', desc: '按分组与优先级触达客户', iconName: 'megaphone', path: '/pages/me/messages' },
]

const tools: EntryItem[] = [
  { key: 'decor', title: t('marketing.decor'), desc: '移动端首页组件编排', iconName: 'decor', path: '/pages/decor/editor' },
  { key: 'specTemplate', title: t('marketing.specTemplate'), desc: '预定义属性组与 SKU 生成', iconName: 'clipboard', path: '/pages/product/spec-templates' },
]

function go(item: EntryItem) {
  if (item.path) uni.navigateTo({ url: item.path })
}
</script>

<template>
  <view class="page">
    <view class="intro">
      <view class="intro__icon">
        <AppIcon name="target" :size="40" color="primary" />
      </view>
      <view class="intro__text">
        <text class="intro__title">营销中心</text>
        <text class="intro__desc">统一管理活动、会员与客户触达渠道</text>
      </view>
    </view>

    <view class="group">
      <SectionHeader title="活动与优惠" />
      <EntryList :items="promotions" @click="go" />
    </view>

    <view class="group">
      <SectionHeader title="会员与运营" />
      <EntryList :items="operations" @click="go" />
    </view>

    <view class="group">
      <SectionHeader title="内容与工具" />
      <EntryList :items="tools" @click="go" />
    </view>
  </view>
</template>

<style scoped>
.page {
  min-height: 100vh;
  background: var(--eapp-bg);
  padding: var(--eapp-space-5);
  padding-bottom: calc(env(safe-area-inset-bottom) + 48rpx);
  box-sizing: border-box;
}

.intro {
  display: flex;
  align-items: center;
  gap: var(--eapp-space-3);
  padding: var(--eapp-space-4) var(--eapp-space-5);
  margin-bottom: var(--eapp-space-6);
  background: var(--eapp-card);
  border: 1rpx solid var(--eapp-border);
  border-radius: var(--eapp-radius-lg);
  box-shadow: var(--eapp-shadow);
  box-sizing: border-box;
}

.intro__icon {
  width: 72rpx;
  height: 72rpx;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--eapp-radius);
  background: var(--eapp-primary-soft);
  border: 1rpx solid var(--eapp-primary-border);
  box-sizing: border-box;
}

.intro__text {
  display: flex;
  flex-direction: column;
  gap: 4rpx;
  min-width: 0;
}

.intro__title {
  font-size: var(--eapp-fs-md);
  font-weight: 700;
  color: var(--eapp-text);
}

.intro__desc {
  font-size: var(--eapp-fs-sm);
  color: var(--eapp-text-muted);
}

.group {
  margin-bottom: var(--eapp-space-6);
}
</style>

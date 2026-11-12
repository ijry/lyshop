<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { hasPermission } from '@/utils/permission'
import { useAuthStore } from '@/stores/auth'
import PageHeader from '@/components/biz/PageHeader.vue'
import SectionHeader from '@/components/biz/SectionHeader.vue'
import EntryList, { type EntryItem } from '@/components/biz/EntryList.vue'
import AppIcon from '@/components/common/AppIcon.vue'

const { t } = useI18n()
const authStore = useAuthStore()

const TAB_PAGES = [
  '/pages/dashboard/index',
  '/pages/order/list',
  '/pages/me/im-sessions',
  '/pages/me/index',
]

const initial = computed(() => String(authStore.username || 'A').charAt(0).toUpperCase())

const shopMenus = computed<EntryItem[]>(() => {
  const raw = [
    { key: 'site', title: t('me.siteSettings'), desc: '店铺名称、Logo 与基础信息', iconName: 'settings', path: '/pages/me/site-settings', perm: 'system.site.update' },
    { key: 'admins', title: t('me.admins'), desc: '子账号与访问权限', iconName: 'users', path: '/pages/me/admins', perm: 'system.admin.list' },
    { key: 'roles', title: t('me.roles'), desc: '角色与细粒度权限配置', iconName: 'shield', path: '/pages/me/roles', perm: 'system.role.list' },
  ]
  if (!authStore.perms.length) return raw
  return raw.filter((item) => hasPermission(authStore.perms, item.perm))
})

const imMenus = computed<EntryItem[]>(() => {
  const raw = [
    { key: 'autoReply', title: t('me.autoReply'), desc: '快捷回复与自动应答规则', iconName: 'chat', path: '/pages/im/auto-replies', perm: 'im.auto_reply.list' },
    { key: 'messages', title: '消息群发', desc: '按分组与优先级推送消息', iconName: 'megaphone', path: '/pages/me/messages' },
    { key: 'imSessions', title: '客服会话', desc: '接待买家咨询与转接', iconName: 'headset', path: '/pages/me/im-sessions' },
  ]
  if (!authStore.perms.length) return raw
  return raw.filter((item) => !item.perm || hasPermission(authStore.perms, item.perm))
})

const toolMenus: EntryItem[] = [
  { key: 'analytics', title: '数据分析', desc: '销售 / 商品 / 客户 / 流量 / 转化', iconName: 'chart', path: '/pages/analytics/index' },
  { key: 'decor', title: '店铺装修', desc: '移动端首页组件编排', iconName: 'decor', path: '/pages/decor/editor' },
  { key: 'wms', title: '仓储管理', desc: '仓库、库存台账与出入库', iconName: 'warehouse', path: '/pages/wms/index' },
]

function go(path?: string) {
  if (!path) return
  if (TAB_PAGES.includes(path)) uni.switchTab({ url: path })
  else uni.navigateTo({ url: path })
}

function onEntry(item: EntryItem) {
  go(item.path)
}

function onLogout() {
  uni.showModal({
    title: '退出登录',
    content: '确认退出当前商家账号？',
    confirmText: '退出',
    confirmColor: '#b42318',
    success: (res) => {
      if (res.confirm) authStore.logout()
    },
  })
}
</script>

<template>
  <view class="page">
    <PageHeader title="我的" eyebrow="ACCOUNT & SETTINGS" />

    <view class="content">
      <!-- 账号卡 -->
      <view class="profile">
        <view class="profile__avatar">
          <text class="profile__avatar-text">{{ initial }}</text>
        </view>
        <view class="profile__info">
          <text class="profile__name">{{ authStore.username || 'admin' }}</text>
          <text class="profile__role">商家管理员 · 已登录</text>
        </view>
        <view class="profile__badge">
          <AppIcon name="shield" :size="26" color="white" />
          <text class="profile__badge-text">已认证</text>
        </view>
      </view>

      <view class="group">
        <SectionHeader title="店铺管理" />
        <EntryList :items="shopMenus" @click="onEntry" />
      </view>

      <view class="group">
        <SectionHeader title="客服与消息" />
        <EntryList :items="imMenus" @click="onEntry" />
      </view>

      <view class="group">
        <SectionHeader title="经营工具" />
        <EntryList :items="toolMenus" @click="onEntry" />
      </view>

      <view class="logout" hover-class="eapp-hover" @click="onLogout">
        <AppIcon name="logout" :size="32" color="danger" />
        <text class="logout__text">{{ t('me.logout') }}</text>
      </view>

      <text class="version">LYShop 商家工作台 · v2.1.0</text>
    </view>
  </view>
</template>

<style scoped>
.page {
  min-height: 100vh;
  background: var(--eapp-bg);
}

.content {
  padding: var(--eapp-space-5);
  padding-bottom: calc(env(safe-area-inset-bottom) + 48rpx);
  box-sizing: border-box;
}

.profile {
  display: flex;
  align-items: center;
  gap: var(--eapp-space-4);
  padding: var(--eapp-space-5);
  margin-bottom: var(--eapp-space-6);
  background: var(--eapp-primary-deep);
  border-radius: var(--eapp-radius-lg);
  box-sizing: border-box;
}

.profile__avatar {
  width: 96rpx;
  height: 96rpx;
  flex-shrink: 0;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.16);
  border: 1rpx solid rgba(255, 255, 255, 0.35);
  display: flex;
  align-items: center;
  justify-content: center;
  box-sizing: border-box;
}

.profile__avatar-text {
  font-size: 40rpx;
  font-weight: 700;
  color: #ffffff;
}

.profile__info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 6rpx;
}

.profile__name {
  font-size: var(--eapp-fs-lg);
  font-weight: 700;
  color: #ffffff;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.profile__role {
  font-size: var(--eapp-fs-sm);
  color: rgba(255, 255, 255, 0.75);
}

.profile__badge {
  display: flex;
  align-items: center;
  gap: 6rpx;
  padding: 6rpx 12rpx;
  border-radius: var(--eapp-radius-sm);
  background: rgba(255, 255, 255, 0.14);
  flex-shrink: 0;
}

.profile__badge-text {
  font-size: 20rpx;
  color: #ffffff;
}

.group {
  margin-bottom: var(--eapp-space-6);
}

.logout {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12rpx;
  height: 88rpx;
  margin-top: var(--eapp-space-6);
  background: var(--eapp-card);
  border: 1rpx solid var(--eapp-border);
  border-radius: var(--eapp-radius-lg);
  box-sizing: border-box;
}

.logout__text {
  font-size: var(--eapp-fs-md);
  font-weight: 600;
  color: var(--eapp-danger);
}

.version {
  display: block;
  text-align: center;
  margin-top: var(--eapp-space-5);
  font-size: var(--eapp-fs-xs);
  color: var(--eapp-text-faint);
  letter-spacing: 1rpx;
}
</style>

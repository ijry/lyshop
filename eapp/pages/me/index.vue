<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { hasPermission } from '@/utils/permission'
import { useAuthStore } from '@/stores/auth'

const { t } = useI18n()
const authStore = useAuthStore()

const menus = computed(() => {
  const raw = [
    { key: 'messages', title: t('me.messages'), path: '/pages/me/messages', perm: 'message.list' },
    { key: 'sessions', title: t('me.sessions'), path: '/pages/me/im-sessions', perm: 'im.session.list' },
    { key: 'autoReply', title: t('me.autoReply'), path: '/pages/im/auto-replies', perm: 'im.auto_reply.list' },
    { key: 'site', title: t('me.siteSettings'), path: '/pages/me/site-settings', perm: 'system.site.update' },
    { key: 'admins', title: t('me.admins'), path: '/pages/me/admins', perm: 'system.admin.list' },
    { key: 'roles', title: t('me.roles'), path: '/pages/me/roles', perm: 'system.role.list' },
  ]
  if (!authStore.perms.length) return raw
  return raw.filter((item) => hasPermission(authStore.perms, item.perm))
})

function go(path: string) {
  uni.navigateTo({ url: path })
}
</script>

<template>
  <view class="page">
    <view class="header">
      <view class="user">{{ authStore.username || 'admin' }}</view>
      <view class="desc">{{ t('me.title') }}</view>
    </view>

    <view class="menu-list">
      <view v-for="menu in menus" :key="menu.key" class="menu-item" @click="go(menu.path)">
        <text>{{ menu.title }}</text>
        <text class="arrow">></text>
      </view>
    </view>

    <view class="logout">
      <up-button type="error" plain @click="authStore.logout">{{ t('me.logout') }}</up-button>
    </view>
  </view>
</template>

<style scoped>
.page { min-height: 100vh; background: var(--eapp-bg); padding: 20rpx; box-sizing: border-box; }
.header { background: #fff; border: 1px solid var(--eapp-border); border-radius: 24rpx; padding: 24rpx; }
.user { font-size: 34rpx; font-weight: 700; }
.desc { margin-top: 8rpx; color: var(--eapp-text-muted); }
.menu-list { margin-top: 16rpx; display: grid; gap: 12rpx; }
.menu-item { background: #fff; border: 1px solid var(--eapp-border); border-radius: 20rpx; padding: 20rpx; display: flex; align-items: center; justify-content: space-between; }
.arrow { color: var(--eapp-text-muted); }
.logout { margin-top: 28rpx; }
</style>

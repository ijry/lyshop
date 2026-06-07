<template>
  <view class="min-h-screen" style="background: var(--app-page-bg);">
    <view style="background: linear-gradient(135deg, #dc2626 0%, #ef4444 100%); padding: 60px 24px 40px; border-radius: 0 0 24px 24px;">
      <view style="display: flex; align-items: center; gap: 16px;">
        <image :src="user.avatar || 'https://api.dicebear.com/7.x/adventurer/svg?seed=default'"
          mode="aspectFill"
          style="width: 64px; height: 64px; border-radius: 50%; border: 3px solid rgba(255,255,255,0.3);" />
        <view style="flex: 1;">
          <text style="color: #fff; font-size: 20px; font-weight: 700; display: block;">{{ user.nickname || $t('user.notLoggedIn') }}</text>
          <text style="color: rgba(255,255,255,0.7); font-size: 13px; margin-top: 4px; display: block;">{{ user.phone || $t('user.clickToLogin') }}</text>
        </view>
        <view @click="uni.navigateTo({url:'/pages/message/index'})"
          style="width: 36px; height: 36px; border-radius: 50%; background: rgba(255,255,255,0.2); display: flex; align-items: center; justify-content: center; position: relative;">
          <u-icon name="bell" size="18" color="#fff" />
          <view v-if="unreadTotal > 0"
            style="position: absolute; top: -2px; right: -2px; min-width: 14px; height: 14px; background: #facc15; border-radius: 7px; display: flex; align-items: center; justify-content: center;">
            <text style="color: #111; font-size: 9px; font-weight: 700;">{{ unreadTotal }}</text>
          </view>
        </view>
        <view @click="uni.navigateTo({url:'/pages/im/chat'})"
          style="width: 36px; height: 36px; border-radius: 50%; background: rgba(255,255,255,0.2); display: flex; align-items: center; justify-content: center;">
          <u-icon name="kefu-ermai" size="18" color="#fff" />
        </view>
      </view>
    </view>

    <view :style="{ margin: '-20px 16px 0', background: 'var(--app-card-bg)', borderRadius: '16px', padding: '20px', boxShadow: 'var(--app-shadow)' }">
      <view style="display: flex; align-items: center; justify-content: space-between; margin-bottom: 16px;">
        <text :style="{ fontSize: '15px', fontWeight: '700', color: 'var(--app-text-primary)' }">{{ $t('user.myOrders') }}</text>
        <text @click="uni.switchTab({url:'/pages/order/list'})"
          :style="{ fontSize: '12px', color: 'var(--app-text-placeholder)' }">{{ $t('user.allOrders') }}</text>
      </view>
      <view style="display: flex; justify-content: space-around;">
        <view v-for="item in orderEntries" :key="item.label"
          @click="uni.switchTab({url:'/pages/order/list'})"
          style="display: flex; flex-direction: column; align-items: center; gap: 6px;">
          <view style="position: relative;">
            <u-icon :name="item.icon" size="24" :color="iconColor" />
            <view v-if="item.badge"
              style="position: absolute; top: -4px; right: -8px; min-width: 16px; height: 16px; background: #dc2626; border-radius: 8px; display: flex; align-items: center; justify-content: center;">
              <text style="color: #fff; font-size: 10px;">{{ item.badge }}</text>
            </view>
          </view>
          <text :style="{ fontSize: '12px', color: 'var(--app-text-tertiary)' }">{{ item.label }}</text>
        </view>
      </view>
    </view>

    <view :style="{ margin: '12px 16px', background: 'var(--app-card-bg)', borderRadius: '16px', overflow: 'hidden', boxShadow: 'var(--app-shadow-sm)' }">
      <view style="display: flex; padding: 20px 0;">
        <view v-for="entry in quickEntries" :key="entry.label"
          @click="entry.action()"
          style="flex: 1; display: flex; flex-direction: column; align-items: center; gap: 8px;">
          <view :style="{ background: entry.bg, width: '44px', height: '44px', borderRadius: '12px', display: 'flex', alignItems: 'center', justifyContent: 'center' }">
            <u-icon :name="entry.icon" size="22" :color="entry.color" />
          </view>
          <text :style="{ fontSize: '12px', color: 'var(--app-text-secondary)' }">{{ entry.label }}</text>
        </view>
      </view>
    </view>

    <view :style="{ margin: '12px 16px', background: 'var(--app-card-bg)', borderRadius: '16px', overflow: 'hidden', boxShadow: 'var(--app-shadow-sm)' }">
      <!-- Dark mode toggle -->
      <view :style="{ display: 'flex', alignItems: 'center', padding: '14px 20px', borderBottom: '1px solid var(--app-divider-color)' }">
        <u-icon name="setting" size="20" :color="iconColor" />
        <text :style="{ flex: '1', marginLeft: '12px', fontSize: '14px', color: 'var(--app-text-secondary)' }">{{ $t('user.darkMode') }}</text>
        <view style="display: flex; align-items: center; gap: 8px;">
          <text :style="{ fontSize: '12px', color: 'var(--app-text-placeholder)' }">{{ themeModeLabel }}</text>
          <u-switch v-model="isDark" size="22" @change="onThemeToggle" />
        </view>
      </view>
      <!-- Menu cells -->
      <view v-for="(cell, i) in menuCells" :key="cell.label"
        @click="cell.action()"
        :style="{ display: 'flex', alignItems: 'center', padding: '14px 20px', borderBottom: i < menuCells.length - 1 ? '1px solid var(--app-divider-color)' : 'none' }">
        <u-icon :name="cell.icon" size="20" :color="iconColor" />
        <text :style="{ flex: '1', marginLeft: '12px', fontSize: '14px', color: 'var(--app-text-secondary)' }">{{ cell.label }}</text>
        <text v-if="cell.value" :style="{ fontSize: '13px', color: 'var(--app-text-placeholder)', marginRight: '4px' }">{{ cell.value }}</text>
        <u-icon name="arrow-right" size="14" :color="arrowColor" />
      </view>
    </view>

    <!-- Demo preset switcher (mock mode only) -->
    <view v-if="isMock" :style="{ margin: '12px 16px', background: 'var(--app-card-bg)', borderRadius: '16px', overflow: 'hidden', boxShadow: 'var(--app-shadow-sm)' }">
      <view @click="presetExpanded = !presetExpanded"
        :style="{ display: 'flex', alignItems: 'center', padding: '14px 20px' }">
        <u-icon name="grid" size="20" :color="iconColor" />
        <text :style="{ flex: '1', marginLeft: '12px', fontSize: '14px', color: 'var(--app-text-secondary)' }">{{ $t('user.switchPreset') }}</text>
        <text :style="{ fontSize: '12px', color: 'var(--app-text-placeholder)', marginRight: '4px' }">{{ presetCurrentName }}</text>
        <u-icon :name="presetExpanded ? 'arrow-up' : 'arrow-down'" size="14" :color="arrowColor" />
      </view>
      <view v-if="presetExpanded" :style="{ padding: '0 16px 16px' }">
        <view style="display: flex; flex-wrap: wrap; gap: 10px;">
          <view v-for="p in presetList" :key="p.key"
            @click="switchPreset(p.key)"
            :style="{
              padding: '8px 16px',
              borderRadius: '20px',
              fontSize: '13px',
              background: p.key === presetActiveKey ? '#1d4ed8' : 'var(--app-page-bg)',
              color: p.key === presetActiveKey ? '#fff' : 'var(--app-text-secondary)',
            }">
            {{ p.name }}
          </view>
        </view>
      </view>
    </view>

    <view :style="{ margin: '16px 16px 40px', background: 'var(--app-card-bg)', borderRadius: '16px', overflow: 'hidden' }">
      <view @click="logout"
        :style="{ textAlign: 'center', padding: '14px', fontSize: '14px', color: 'var(--app-accent)' }">
        {{ $t('user.logout') }}
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { get } from '@/utils/request'
import { useTheme } from '@/composables/useTheme'

const { t } = useI18n()
const { effectiveTheme, setTheme } = useTheme()

const isDark = ref(effectiveTheme.value === 'dark')
const themeModeLabel = computed(() => isDark.value ? t('user.enabled') : t('user.disabled'))
const iconColor = computed(() => effectiveTheme.value === 'dark' ? '#9ca3af' : '#666')
const arrowColor = computed(() => effectiveTheme.value === 'dark' ? '#6b7280' : '#ccc')

function onThemeToggle(val: boolean) {
  setTheme(val ? 'dark' : 'light')
}

const user = ref<any>({})
const unreadTotal = ref(0)

const orderEntries = computed(() => [
  { label: t('user.pendingPayment'), icon: 'rmb-circle', badge: 0 },
  { label: t('user.pendingShipment'), icon: 'gift', badge: 0 },
  { label: t('user.pendingReceipt'), icon: 'car', badge: 0 },
  { label: t('user.afterSale'), icon: 'reload', badge: 0 },
])

const quickEntries = computed(() => [
  { label: t('user.coupon'), icon: 'coupon', bg: '#fef2f2', color: '#dc2626',
    action: () => uni.navigateTo({ url: '/pages/marketing/coupon?mode=claim' }) },
  { label: t('user.myPoints'), icon: 'integral', bg: '#fff7ed', color: '#f97316',
    action: () => uni.navigateTo({ url: '/pages/user/points' }) },
  { label: t('user.shippingAddress'), icon: 'map', bg: '#eff6ff', color: '#3b82f6',
    action: () => uni.navigateTo({ url: '/pages/user/address' }) },
  { label: t('user.onlineService'), icon: 'kefu-ermai', bg: '#f0fdf4', color: '#22c55e',
    action: () => uni.navigateTo({ url: '/pages/im/chat' }) },
])

const menuCellValues = ref<{ points: string; unread: string }>({ points: '', unread: '' })
const menuCells = computed(() => [
  { label: t('user.dailyCheckin'), icon: 'calendar', value: '',
    action: () => uni.navigateTo({ url: '/pages/checkin/index' }) },
  { label: t('user.myCoupons'), icon: 'coupon', value: '',
    action: () => uni.navigateTo({ url: '/pages/marketing/coupon' }) },
  { label: t('user.myFavorites'), icon: 'heart', value: '',
    action: () => uni.navigateTo({ url: '/pages/user/favorites' }) },
  { label: t('user.vipCenter'), icon: 'level', value: '',
    action: () => uni.navigateTo({ url: '/pages/user/vip' }) },
  { label: t('user.distributionCenter'), icon: 'share', value: '',
    action: () => uni.navigateTo({ url: '/pages/distribution/index' }) },
  { label: t('user.accountSecurity'), icon: 'lock', value: '',
    action: () => uni.navigateTo({ url: '/pages/user/security' }) },
])

// Demo preset switching (mock mode only)
const isMock = import.meta.env.VITE_MOCK === 'true'
const presetExpanded = ref(false)
const presetList = ref<Array<{ key: string; name: string }>>([])
const presetActiveKey = ref('')
const presetCurrentName = ref('')

function switchPreset(key: string) {
  if (key === presetActiveKey.value) {
    presetExpanded.value = false
    return
  }
  const url = new URL(window.location.href)
  if (key === 'mall') {
    url.searchParams.delete('demo')
  } else {
    url.searchParams.set('demo', key)
  }
  window.location.href = url.toString()
}

function logout() {
  uni.removeStorageSync('user_token')
  uni.reLaunch({ url: '/pages/login/index' })
}

onMounted(async () => {
  // Load demo presets if in mock mode
  if (isMock) {
    try {
      const { listPresets, getPresetKey } = await import('@/mock/presets/index')
      presetList.value = listPresets()
      presetActiveKey.value = getPresetKey()
      presetCurrentName.value = presetList.value.find(p => p.key === presetActiveKey.value)?.name || presetActiveKey.value
    } catch {}
  }

  const data = await get<any>('/api/v1/user/profile')
  if (data) {
    user.value = data
    menuCellValues.value.points = `${data.points || 0} ${t('user.pointsUnit')}`
  }
  const unread = await get<any>('/api/v1/messages/unread')
  if (unread) {
    const total = Object.values(unread).reduce((s: number, v: any) => s + (v || 0), 0)
    unreadTotal.value = total as number
    if (total > 0) menuCellValues.value.unread = `${total} ${t('user.unreadCount')}`
  }
})
</script>

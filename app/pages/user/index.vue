<template>
  <view class="min-h-screen" style="background: #f5f5f5;">
    <view style="background: linear-gradient(135deg, #dc2626 0%, #ef4444 100%); padding: 60px 24px 40px; border-radius: 0 0 24px 24px;">
      <view style="display: flex; align-items: center; gap: 16px;">
        <image :src="user.avatar || 'https://api.dicebear.com/7.x/adventurer/svg?seed=default'"
          mode="aspectFill"
          style="width: 64px; height: 64px; border-radius: 50%; border: 3px solid rgba(255,255,255,0.3);" />
        <view style="flex: 1;">
          <text style="color: #fff; font-size: 20px; font-weight: 700; display: block;">{{ user.nickname || '未登录' }}</text>
          <text style="color: rgba(255,255,255,0.7); font-size: 13px; margin-top: 4px; display: block;">{{ user.phone || '点击登录' }}</text>
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

    <view style="margin: -20px 16px 0; background: #fff; border-radius: 16px; padding: 20px; box-shadow: 0 2px 12px rgba(0,0,0,0.06);">
      <view style="display: flex; align-items: center; justify-content: space-between; margin-bottom: 16px;">
        <text style="font-size: 15px; font-weight: 700; color: #111;">我的订单</text>
        <text @click="uni.switchTab({url:'/pages/order/list'})"
          style="font-size: 12px; color: #999;">全部订单 ></text>
      </view>
      <view style="display: flex; justify-content: space-around;">
        <view v-for="item in orderEntries" :key="item.label"
          @click="uni.switchTab({url:'/pages/order/list'})"
          style="display: flex; flex-direction: column; align-items: center; gap: 6px;">
          <view style="position: relative;">
            <u-icon :name="item.icon" size="24" color="#333" />
            <view v-if="item.badge"
              style="position: absolute; top: -4px; right: -8px; min-width: 16px; height: 16px; background: #dc2626; border-radius: 8px; display: flex; align-items: center; justify-content: center;">
              <text style="color: #fff; font-size: 10px;">{{ item.badge }}</text>
            </view>
          </view>
          <text style="font-size: 12px; color: #666;">{{ item.label }}</text>
        </view>
      </view>
    </view>

    <view style="margin: 12px 16px; background: #fff; border-radius: 16px; overflow: hidden; box-shadow: 0 2px 12px rgba(0,0,0,0.04);">
      <view style="display: flex; padding: 20px 0;">
        <view v-for="entry in quickEntries" :key="entry.label"
          @click="entry.action()"
          style="flex: 1; display: flex; flex-direction: column; align-items: center; gap: 8px;">
          <view :style="{ background: entry.bg, width: '44px', height: '44px', borderRadius: '12px', display: 'flex', alignItems: 'center', justifyContent: 'center' }">
            <u-icon :name="entry.icon" size="22" :color="entry.color" />
          </view>
          <text style="font-size: 12px; color: #333;">{{ entry.label }}</text>
        </view>
      </view>
    </view>

    <view style="margin: 12px 16px; background: #fff; border-radius: 16px; overflow: hidden; box-shadow: 0 2px 12px rgba(0,0,0,0.04);">
      <view v-for="(cell, i) in menuCells" :key="cell.label"
        @click="cell.action()"
        :style="{ display: 'flex', alignItems: 'center', padding: '14px 20px', borderBottom: i < menuCells.length - 1 ? '1px solid #f5f5f5' : 'none' }">
        <u-icon :name="cell.icon" size="20" color="#666" />
        <text style="flex: 1; margin-left: 12px; font-size: 14px; color: #333;">{{ cell.label }}</text>
        <text v-if="cell.value" style="font-size: 13px; color: #999; margin-right: 4px;">{{ cell.value }}</text>
        <u-icon name="arrow-right" size="14" color="#ccc" />
      </view>
    </view>

    <view style="margin: 16px 16px 40px; background: #fff; border-radius: 16px; overflow: hidden;">
      <view @click="logout"
        style="text-align: center; padding: 14px; font-size: 14px; color: #dc2626;">
        退出登录
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { get } from '@/utils/request'

const user = ref<any>({})
const unreadTotal = ref(0)

const orderEntries = [
  { label: '待付款', icon: 'rmb-circle', badge: 0 },
  { label: '待发货', icon: 'gift', badge: 0 },
  { label: '待收货', icon: 'car', badge: 0 },
  { label: '售后', icon: 'reload', badge: 0 },
]

const quickEntries = [
  { label: '优惠券', icon: 'coupon', bg: '#fef2f2', color: '#dc2626',
    action: () => uni.navigateTo({ url: '/pages/marketing/coupon?mode=claim' }) },
  { label: '我的积分', icon: 'integral', bg: '#fff7ed', color: '#f97316',
    action: () => uni.navigateTo({ url: '/pages/user/points' }) },
  { label: '收货地址', icon: 'map', bg: '#eff6ff', color: '#3b82f6',
    action: () => uni.navigateTo({ url: '/pages/user/address' }) },
  { label: '在线客服', icon: 'kefu-ermai', bg: '#f0fdf4', color: '#22c55e',
    action: () => uni.navigateTo({ url: '/pages/im/chat' }) },
]

const menuCells = ref([
  { label: '消息中心', icon: 'bell', value: '',
    action: () => uni.navigateTo({ url: '/pages/message/index' }) },
  { label: '每日签到', icon: 'calendar', value: '',
    action: () => uni.navigateTo({ url: '/pages/checkin/index' }) },
  { label: '收货地址', icon: 'map', value: '',
    action: () => uni.navigateTo({ url: '/pages/user/address' }) },
  { label: '我的优惠券', icon: 'coupon', value: '',
    action: () => uni.navigateTo({ url: '/pages/marketing/coupon' }) },
  { label: '我的收藏', icon: 'heart', value: '',
    action: () => uni.navigateTo({ url: '/pages/user/favorites' }) },
  { label: '我的积分', icon: 'integral', value: '',
    action: () => uni.navigateTo({ url: '/pages/user/points' }) },
  { label: '会员中心', icon: 'level', value: '',
    action: () => uni.navigateTo({ url: '/pages/user/vip' }) },
  { label: '账号与安全', icon: 'lock', value: '',
    action: () => uni.navigateTo({ url: '/pages/user/security' }) },
  { label: '联系客服', icon: 'kefu-ermai', value: '',
    action: () => uni.navigateTo({ url: '/pages/im/chat' }) },
])

function logout() {
  uni.removeStorageSync('user_token')
  uni.reLaunch({ url: '/pages/login/index' })
}

onMounted(async () => {
  const data = await get<any>('/api/v1/user/profile')
  if (data) {
    user.value = data
    menuCells.value[4].value = `${data.points || 0} 积分`
  }
  const unread = await get<any>('/api/v1/messages/unread')
  if (unread) {
    const total = Object.values(unread).reduce((s: number, v: any) => s + (v || 0), 0)
    unreadTotal.value = total as number
    if (total > 0) menuCells.value[0].value = `${total} 条未读`
  }
})
</script>

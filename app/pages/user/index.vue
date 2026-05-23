<template>
  <view class="min-h-screen" style="background: #f5f5f5;">
    <!-- Header with gradient -->
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
          <u-icon name="chat" size="18" color="#fff" />
        </view>
      </view>
    </view>

    <!-- Order status cards -->
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

    <!-- Quick entries -->
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

    <!-- Menu cells -->
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

    <!-- Logout + Delete Account -->
    <view style="margin: 16px 16px 40px; background: #fff; border-radius: 16px; overflow: hidden;">
      <view @click="logout"
        style="text-align: center; padding: 14px; font-size: 14px; color: #dc2626; border-bottom: 1px solid #f5f5f5;">
        退出登录
      </view>
      <view @click="showDeleteConfirm = true"
        style="text-align: center; padding: 14px; font-size: 13px; color: #999;">
        注销账号
      </view>
    </view>

    <!-- Delete account popup -->
    <u-popup :show="showDeleteConfirm" mode="center" round="20" @close="showDeleteConfirm=false">
      <view style="padding: 30px; width: 300px;">
        <text style="font-size: 17px; font-weight: 700; color: #111; display: block; text-align: center;">注销账号</text>
        <text style="font-size: 13px; color: #999; display: block; text-align: center; margin: 12px 0 20px; line-height: 1.5;">
          注销后账号数据将被永久删除且无法恢复，请谨慎操作。需要短信验证码确认身份。
        </text>
        <view style="margin-bottom: 12px;">
          <u-input v-model="deleteForm.phone" placeholder="手机号" type="number" :maxlength="11" border="surround" shape="circle" />
        </view>
        <view style="display: flex; gap: 10px; margin-bottom: 20px;">
          <view style="flex: 1;">
            <u-input v-model="deleteForm.code" placeholder="验证码" type="number" :maxlength="6" border="surround" shape="circle" />
          </view>
          <u-button size="small" :disabled="deleteCountdown > 0"
            :text="deleteCountdown > 0 ? `${deleteCountdown}s` : '获取验证码'"
            @click="sendDeleteCode" type="primary" plain shape="circle" />
        </view>
        <view style="display: flex; gap: 10px;">
          <u-button text="取消" @click="showDeleteConfirm=false" shape="circle" class="flex-1" />
          <u-button text="确认注销" type="error" @click="deleteAccount" shape="circle" class="flex-1"
            :custom-style="{background: '#dc2626', borderColor: '#dc2626'}" />
        </view>
      </view>
    </u-popup>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { get, post } from '@/utils/request'

const user = ref<any>({})
const unreadTotal = ref(0)

const orderEntries = [
  { label: '待付款', icon: 'wallet', badge: 0 },
  { label: '待发货', icon: 'gift', badge: 0 },
  { label: '待收货', icon: 'car', badge: 0 },
  { label: '售后', icon: 'reload', badge: 0 },
]

const quickEntries = [
  { label: '优惠券', icon: 'coupon', bg: '#fef2f2', color: '#dc2626',
    action: () => uni.navigateTo({ url: '/pages/marketing/coupon' }) },
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
  { label: '我的积分', icon: 'integral', value: '',
    action: () => uni.navigateTo({ url: '/pages/user/points' }) },
  { label: '联系客服', icon: 'kefu-ermai', value: '',
    action: () => uni.navigateTo({ url: '/pages/im/chat' }) },
])

function logout() {
  uni.removeStorageSync('user_token')
  uni.reLaunch({ url: '/pages/login/index' })
}

// Account deletion
const showDeleteConfirm = ref(false)
const deleteForm = ref({ phone: '', code: '' })
const deleteCountdown = ref(0)

async function sendDeleteCode() {
  if (!deleteForm.value.phone || deleteForm.value.phone.length !== 11) {
    uni.showToast({ title: '请输入手机号', icon: 'none' }); return
  }
  try {
    const data = await post<any>('/api/v1/auth/sms/send', { phone: deleteForm.value.phone })
    if (data?.dev_code) deleteForm.value.code = data.dev_code
  } catch {}
  deleteCountdown.value = 60
  const t = setInterval(() => { if (--deleteCountdown.value <= 0) clearInterval(t) }, 1000)
}

async function deleteAccount() {
  if (!deleteForm.value.code) {
    uni.showToast({ title: '请输入验证码', icon: 'none' }); return
  }
  try {
    await post('/api/v1/user/delete', deleteForm.value)
    uni.showToast({ title: '账号已注销', icon: 'success' })
    setTimeout(() => {
      uni.removeStorageSync('user_token')
      uni.reLaunch({ url: '/pages/login/index' })
    }, 1500)
  } catch (e: any) {
    uni.showToast({ title: e.message || '注销失败', icon: 'none' })
  }
}

onMounted(async () => {
  const data = await get<any>('/api/v1/user/profile')
  if (data) {
    user.value = data
    menuCells.value[4].value = `${data.points || 0} 积分`
  }
  // Unread message count
  const unread = await get<any>('/api/v1/messages/unread')
  if (unread) {
    const total = Object.values(unread).reduce((s: number, v: any) => s + (v || 0), 0)
    unreadTotal.value = total as number
    if (total > 0) menuCells.value[0].value = `${total} 条未读`
  }
})
</script>

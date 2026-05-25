<template>
  <view class="min-h-screen" style="background:#f5f5f5;padding:16px;">
    <view style="background:#fff;border-radius:14px;padding:16px;margin-bottom:12px;">
      <text style="font-size:16px;font-weight:700;color:#111;display:block;">会员状态</text>
      <text style="font-size:14px;color:#666;display:block;margin-top:10px;">
        {{ profile.is_vip ? `已开通 · ${profile.level_name}` : '未开通' }}
      </text>
      <text style="font-size:13px;color:#999;display:block;margin-top:6px;">成长值：{{ profile.growth_value || 0 }}</text>
      <text v-if="profile.expire_at" style="font-size:13px;color:#999;display:block;margin-top:6px;">到期时间：{{ profile.expire_at }}</text>
    </view>

    <view style="background:#fff;border-radius:14px;padding:16px;">
      <text style="font-size:16px;font-weight:700;color:#111;display:block;margin-bottom:12px;">本月可领券</text>
      <view v-for="item in coupons" :key="item.rule_id" style="display:flex;align-items:center;justify-content:space-between;padding:10px 0;border-bottom:1px solid #f1f1f1;">
        <view>
          <text style="font-size:14px;color:#111;display:block;">{{ item.name }}</text>
          <text style="font-size:12px;color:#666;display:block;margin-top:4px;">{{ item.coupon_name }}（{{ item.claimed }}/{{ item.monthly_limit }}）</text>
        </view>
        <button
          :disabled="loadingMap[item.rule_id] || item.claimed >= item.monthly_limit"
          @click="claim(item.rule_id)"
          style="margin:0;font-size:12px;line-height:30px;height:30px;padding:0 14px;background:#dc2626;color:#fff;border-radius:15px;border:none;"
        >
          {{ item.claimed >= item.monthly_limit ? '已领完' : '领取' }}
        </button>
      </view>
      <view v-if="!coupons.length" style="text-align:center;color:#999;font-size:13px;padding:20px 0;">暂无可领取优惠券</view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { get, post } from '@/utils/request'

const profile = ref<any>({})
const coupons = ref<any[]>([])
const loadingMap = ref<Record<number, boolean>>({})

async function load() {
  profile.value = (await get<any>('/api/v1/vip/profile')) || {}
  coupons.value = (await get<any[]>('/api/v1/vip/coupons/monthly')) || []
}

async function claim(ruleID: number) {
  if (loadingMap.value[ruleID]) return
  loadingMap.value[ruleID] = true
  try {
    await post(`/api/v1/vip/coupons/monthly/${ruleID}/claim`)
    await load()
    uni.showToast({ title: '领取成功', icon: 'success' })
  } finally {
    loadingMap.value[ruleID] = false
  }
}

onMounted(load)
</script>

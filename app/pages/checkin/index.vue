<template>
  <view style="min-height: 100vh; background: #f5f5f5;">
    <u-navbar title="每日签到" :placeholder="true" />

    <!-- Points card -->
    <view style="margin: 16px; background: linear-gradient(135deg, #dc2626, #f97316); border-radius: 20px; padding: 24px; color: #fff; position: relative; overflow: hidden;">
      <view style="position: absolute; right: -20px; top: -20px; width: 120px; height: 120px; border-radius: 60px; background: rgba(255,255,255,0.08);" />
      <view style="display: flex; align-items: center; justify-content: space-between;">
        <view>
          <text style="font-size: 13px; opacity: 0.8;">已连续签到</text>
          <view style="display: flex; align-items: baseline; gap: 4px; margin-top: 4px;">
            <text style="font-size: 42px; font-weight: 800;">{{ status.consecutive_days || 0 }}</text>
            <text style="font-size: 14px; opacity: 0.7;">天</text>
          </view>
          <text style="font-size: 12px; opacity: 0.6; margin-top: 4px; display: block;">本月签到 {{ status.month_count || 0 }} 天，获得 {{ status.month_points || 0 }} 积分</text>
        </view>
        <view @click="doCheckin"
          :style="{ width: '80px', height: '80px', borderRadius: '50%', display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', background: status.checked_today ? 'rgba(255,255,255,0.15)' : '#fff', cursor: 'pointer' }">
          <text :style="{ fontSize: '13px', fontWeight: '700', color: status.checked_today ? 'rgba(255,255,255,0.6)' : '#dc2626' }">
            {{ status.checked_today ? '已签到' : '签到' }}
          </text>
          <text v-if="!status.checked_today" style="font-size: 10px; color: #f97316; margin-top: 2px;">+{{ nextPoints }}积分</text>
        </view>
      </view>
    </view>

    <!-- 7-day rule preview -->
    <view style="margin: 0 16px 16px; background: #fff; border-radius: 16px; padding: 20px;">
      <text style="font-size: 15px; font-weight: 700; color: #111; display: block; margin-bottom: 16px;">签到奖励规则</text>
      <view style="display: flex; justify-content: space-between;">
        <view v-for="d in 7" :key="d" style="display: flex; flex-direction: column; align-items: center; gap: 6px;">
          <view :style="{
            width: '36px', height: '36px', borderRadius: '50%', display: 'flex', alignItems: 'center', justifyContent: 'center',
            background: (status.consecutive_days || 0) >= d ? '#dc2626' : '#f5f5f5',
            color: (status.consecutive_days || 0) >= d ? '#fff' : '#999',
            fontSize: '11px', fontWeight: '600'
          }">
            {{ getPointsForDay(d) }}
          </view>
          <text style="font-size: 10px; color: #999;">第{{ d }}天</text>
        </view>
      </view>
    </view>

    <!-- Calendar dots -->
    <view style="margin: 0 16px; background: #fff; border-radius: 16px; padding: 20px;">
      <text style="font-size: 15px; font-weight: 700; color: #111; display: block; margin-bottom: 12px;">本月签到日历</text>
      <view style="display: flex; flex-wrap: wrap; gap: 6px;">
        <view v-for="d in daysInMonth" :key="d"
          :style="{
            width: '36px', height: '36px', borderRadius: '8px', display: 'flex', alignItems: 'center', justifyContent: 'center',
            fontSize: '12px', fontWeight: '500',
            background: isCheckedDate(d) ? '#fef2f2' : '#fafafa',
            color: isCheckedDate(d) ? '#dc2626' : '#999',
            border: isToday(d) ? '2px solid #dc2626' : '1px solid transparent'
          }">
          {{ d }}
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { get, post } from '@/utils/request'

const status = ref<any>({})
const rules = ref<any[]>([])

const now = new Date()
const daysInMonth = new Date(now.getFullYear(), now.getMonth() + 1, 0).getDate()
const todayStr = now.toISOString().slice(0, 10)

const nextPoints = computed(() => {
  const consecutive = (status.value.consecutive_days || 0) + (status.value.checked_today ? 0 : 1)
  return getPointsForDay(consecutive)
})

function getPointsForDay(day: number): number {
  if (!rules.value.length) return 10
  let result = rules.value.find(r => r.day === 0)?.points || 10
  for (const r of rules.value) {
    if (r.day > 0 && r.day === day) return r.points
    if (r.day > 0 && r.day <= day) result = r.points
  }
  return result
}

function isCheckedDate(day: number) {
  const dateStr = `${now.getFullYear()}-${String(now.getMonth()+1).padStart(2,'0')}-${String(day).padStart(2,'0')}`
  return (status.value.month_dates || []).includes(dateStr)
}

function isToday(day: number) {
  return day === now.getDate()
}

async function doCheckin() {
  if (status.value.checked_today) {
    uni.showToast({ title: '今日已签到', icon: 'none' })
    return
  }
  try {
    const data = await post<any>('/api/v1/checkin')
    uni.showToast({ title: `签到成功！+${data.points}积分`, icon: 'success' })
    await loadStatus()
  } catch (e: any) {
    uni.showToast({ title: e.message || '签到失败', icon: 'none' })
  }
}

async function loadStatus() {
  status.value = (await get<any>('/api/v1/checkin/status')) || {}
  rules.value = (await get<any[]>('/api/v1/checkin/rules')) || []
}

onMounted(loadStatus)
</script>

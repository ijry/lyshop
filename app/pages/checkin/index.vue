<template>
  <view style="min-height: 100vh; background: linear-gradient(180deg, #dc2626 0%, #f5f5f5 35%);">
    <u-navbar title="每日签到" :placeholder="true" :bg-color="'transparent'" title-color="#fff" />

    <!-- Hero card -->
    <view style="margin: 0 16px 16px; padding: 28px 24px; position: relative; overflow: hidden;">
      <!-- Decorative circles -->
      <view style="position: absolute; right: -30px; top: -30px; width: 140px; height: 140px; border-radius: 70px; background: rgba(255,255,255,0.06);" />
      <view style="position: absolute; right: 40px; bottom: -50px; width: 100px; height: 100px; border-radius: 50px; background: rgba(255,255,255,0.04);" />

      <view style="display: flex; align-items: center; justify-content: space-between; position: relative; z-index: 1;">
        <view>
          <text style="font-size: 14px; color: rgba(255,255,255,0.75);">已连续签到</text>
          <view style="display: flex; align-items: baseline; gap: 4px; margin-top: 6px;">
            <text style="font-size: 52px; font-weight: 800; color: #fff; line-height: 1;">{{ status.consecutive_days || 0 }}</text>
            <text style="font-size: 16px; color: rgba(255,255,255,0.6);">天</text>
          </view>
          <view style="display: flex; align-items: center; gap: 12px; margin-top: 10px;">
            <view style="display: flex; align-items: center; gap: 4px; background: rgba(255,255,255,0.12); border-radius: 20px; padding: 3px 10px;">
              <text style="font-size: 11px; color: rgba(255,255,255,0.8);">本月 {{ status.month_count || 0 }} 天</text>
            </view>
            <view style="display: flex; align-items: center; gap: 4px; background: rgba(255,255,255,0.12); border-radius: 20px; padding: 3px 10px;">
              <text style="font-size: 11px; color: rgba(255,255,255,0.8);">获得 {{ status.month_points || 0 }} 积分</text>
            </view>
          </view>
        </view>

        <!-- Checkin button -->
        <view @click="doCheckin"
          :style="{
            width: '88px', height: '88px', borderRadius: '50%',
            display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center',
            background: status.checked_today ? 'rgba(255,255,255,0.12)' : '#fff',
            boxShadow: status.checked_today ? 'none' : '0 4px 20px rgba(0,0,0,0.15)',
            transition: 'all 0.3s'
          }">
          <text :style="{ fontSize: '14px', fontWeight: '700', color: status.checked_today ? 'rgba(255,255,255,0.5)' : '#dc2626' }">
            {{ status.checked_today ? '已签' : '签到' }}
          </text>
          <text v-if="!status.checked_today" style="font-size: 10px; color: #f97316; margin-top: 2px; font-weight: 600;">+{{ nextPoints }}</text>
        </view>
      </view>
    </view>

    <!-- 7-day consecutive reward track -->
    <view style="margin: 0 16px 12px; background: #fff; border-radius: 20px; padding: 20px 16px; box-shadow: 0 2px 12px rgba(0,0,0,0.04);">
      <view style="display: flex; align-items: center; justify-content: space-between; margin-bottom: 18px;">
        <text style="font-size: 15px; font-weight: 700; color: #111;">连续签到奖励</text>
        <text style="font-size: 11px; color: #999;">连续7天额外奖励50积分</text>
      </view>

      <!-- Track line with dots -->
      <view style="position: relative; padding: 0 8px;">
        <!-- Background line -->
        <view style="position: absolute; top: 18px; left: 24px; right: 24px; height: 3px; background: #f0f0f0; border-radius: 2px;" />
        <!-- Active line -->
        <view :style="{ position: 'absolute', top: '18px', left: '24px', width: Math.min(((status.consecutive_days || 0) / 7) * 100, 100) + '%', maxWidth: 'calc(100% - 48px)', height: '3px', background: 'linear-gradient(90deg, #dc2626, #f97316)', borderRadius: '2px', transition: 'width 0.5s' }" />

        <view style="display: flex; justify-content: space-between; position: relative;">
          <view v-for="d in 7" :key="d" style="display: flex; flex-direction: column; align-items: center; gap: 8px; width: 40px;">
            <!-- Dot/Circle -->
            <view :style="{
              width: d === 7 ? '38px' : '32px',
              height: d === 7 ? '38px' : '32px',
              borderRadius: '50%',
              display: 'flex', alignItems: 'center', justifyContent: 'center',
              background: (status.consecutive_days || 0) >= d ? (d === 7 ? 'linear-gradient(135deg, #dc2626, #f97316)' : '#dc2626') : '#fff',
              border: (status.consecutive_days || 0) >= d ? 'none' : '2px solid #e5e5e5',
              boxShadow: (status.consecutive_days || 0) >= d ? '0 2px 8px rgba(220,38,38,0.3)' : 'none',
              transition: 'all 0.3s'
            }">
              <text v-if="(status.consecutive_days || 0) >= d" style="color: #fff; font-size: 12px;">✓</text>
              <text v-else style="color: #ccc; font-size: 10px; font-weight: 600;">{{ d }}</text>
            </view>
            <!-- Label -->
            <view style="text-align: center;">
              <text :style="{ fontSize: '11px', fontWeight: d === 7 ? '700' : '500', color: (status.consecutive_days || 0) >= d ? '#dc2626' : '#999' }">
                +{{ getPointsForDay(d) }}
              </text>
              <text style="display: block; font-size: 9px; color: #ccc; margin-top: 1px;">第{{ d }}天</text>
            </view>
          </view>
        </view>
      </view>
    </view>

    <!-- Calendar -->
    <view style="margin: 0 16px 16px; background: #fff; border-radius: 20px; padding: 20px; box-shadow: 0 2px 12px rgba(0,0,0,0.04);">
      <view style="display: flex; align-items: center; justify-content: space-between; margin-bottom: 16px;">
        <text style="font-size: 15px; font-weight: 700; color: #111;">{{ now.getFullYear() }}年{{ now.getMonth()+1 }}月</text>
        <text style="font-size: 12px; color: #dc2626; font-weight: 600;">{{ status.month_count || 0 }}/{{ daysInMonth }} 天</text>
      </view>
      <!-- Weekday headers -->
      <view style="display: flex; margin-bottom: 8px;">
        <view v-for="w in ['日','一','二','三','四','五','六']" :key="w"
          style="flex: 1; text-align: center; font-size: 11px; color: #bbb; font-weight: 500;">
          {{ w }}
        </view>
      </view>
      <!-- Calendar grid -->
      <view style="display: flex; flex-wrap: wrap;">
        <!-- Empty cells for offset -->
        <view v-for="_ in firstDayOffset" :key="'e'+_" style="width: 14.28%; aspect-ratio: 1;" />
        <!-- Day cells -->
        <view v-for="d in daysInMonth" :key="d" style="width: 14.28%; display: flex; align-items: center; justify-content: center; padding: 3px 0;">
          <view :style="{
            width: '34px', height: '34px', borderRadius: '10px',
            display: 'flex', alignItems: 'center', justifyContent: 'center',
            fontSize: '13px', fontWeight: isToday(d) ? '700' : '500',
            background: isCheckedDate(d) ? '#fef2f2' : 'transparent',
            color: isCheckedDate(d) ? '#dc2626' : (d > now.getDate() ? '#ddd' : '#666'),
            border: isToday(d) ? '2px solid #dc2626' : 'none',
            position: 'relative'
          }">
            {{ d }}
            <view v-if="isCheckedDate(d)" style="position: absolute; bottom: 2px; width: 4px; height: 4px; border-radius: 2px; background: #dc2626;" />
          </view>
        </view>
      </view>
    </view>

    <!-- Rules card -->
    <view style="margin: 0 16px 30px; background: #fff; border-radius: 20px; padding: 20px; box-shadow: 0 2px 12px rgba(0,0,0,0.04);">
      <text style="font-size: 15px; font-weight: 700; color: #111; display: block; margin-bottom: 12px;">签到规则</text>
      <view style="display: flex; flex-direction: column; gap: 8px;">
        <view style="display: flex; align-items: center; gap: 8px;">
          <view style="width: 6px; height: 6px; border-radius: 3px; background: #dc2626;" />
          <text style="font-size: 13px; color: #666; line-height: 1.5;">每日签到可获得基础积分奖励</text>
        </view>
        <view style="display: flex; align-items: center; gap: 8px;">
          <view style="width: 6px; height: 6px; border-radius: 3px; background: #f97316;" />
          <text style="font-size: 13px; color: #666; line-height: 1.5;">连续签到天数越多，奖励积分越高</text>
        </view>
        <view style="display: flex; align-items: center; gap: 8px;">
          <view style="width: 6px; height: 6px; border-radius: 3px; background: #eab308;" />
          <text style="font-size: 13px; color: #666; line-height: 1.5;">中断签到后连续天数将重新计算</text>
        </view>
        <view style="display: flex; align-items: center; gap: 8px;">
          <view style="width: 6px; height: 6px; border-radius: 3px; background: #22c55e;" />
          <text style="font-size: 13px; color: #666; line-height: 1.5;">100积分 = ¥1.00，下单时可抵扣</text>
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
const firstDayOffset = new Date(now.getFullYear(), now.getMonth(), 1).getDay()

const nextPoints = computed(() => {
  const consecutive = (status.value.consecutive_days || 0) + (status.value.checked_today ? 0 : 1)
  return getPointsForDay(consecutive)
})

function getPointsForDay(day: number): number {
  if (!rules.value.length) return 10
  let result = rules.value.find((r: any) => r.day === 0)?.points || 10
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

function isToday(day: number) { return day === now.getDate() }

async function doCheckin() {
  if (status.value.checked_today) {
    uni.showToast({ title: '今日已签到', icon: 'none' }); return
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

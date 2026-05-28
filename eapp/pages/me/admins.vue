<script setup lang="ts">
import { onShow } from '@dcloudio/uni-app'
import { ref } from 'vue'
import { getAdmins, getRoles } from '@/api/system'

const list = ref<any[]>([])
const roleMap = ref<Record<number, string>>({})

async function loadData() {
  const [admins, roles] = await Promise.all([getAdmins(), getRoles()])
  list.value = Array.isArray(admins) ? admins : []
  const map: Record<number, string> = {}
  ;(Array.isArray(roles) ? roles : []).forEach((role: any) => { map[Number(role.id)] = String(role.name || '-') })
  roleMap.value = map
}

onShow(loadData)
</script>

<template>
  <view class="page">
    <view v-if="!list.length" class="empty">暂无管理员</view>
    <view v-for="item in list" :key="item.id" class="card">
      <view class="title">{{ item.username || '-' }}</view>
      <view class="desc">角色：{{ roleMap[Number(item.role_id)] || '-' }}</view>
    </view>
  </view>
</template>

<style scoped>
.page { min-height: 100vh; background: var(--eapp-bg); padding: 20rpx; box-sizing: border-box; display: grid; gap: 12rpx; }
.card { background: #fff; border: 1px solid var(--eapp-border); border-radius: 20rpx; padding: 20rpx; }
.title { font-size: 28rpx; font-weight: 600; }
.desc { margin-top: 8rpx; color: var(--eapp-text-muted); }
.empty { text-align: center; color: var(--eapp-text-muted); padding: 80rpx 0; }
</style>

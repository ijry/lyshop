<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { getCategoriesTree } from '@/api/category'

const props = defineProps<{ show: boolean; value: number }>()
const emit = defineEmits<{ (e: 'close'): void; (e: 'pick', payload: { id: number; path_name: string }): void }>()

const tree = ref<any[]>([])
const expanded = ref<Record<number, boolean>>({})

async function load() { tree.value = (await getCategoriesTree()) || [] }
function toggle(id: number) { expanded.value[id] = !expanded.value[id] }
function pick(node: any, path: string[]) { emit('pick', { id: Number(node.id), path_name: path.concat(node.name).join(' / ') }) }

onMounted(load)
</script>

<template>
  <up-popup :show="show" mode="right" round="0" @close="$emit('close')">
    <view class="drawer">
      <view class="head">
        <text class="title">选择分类</text>
        <text class="close" @click="$emit('close')">关闭</text>
      </view>
      <scroll-view scroll-y class="body">
        <view v-for="root in tree" :key="root.id">
          <view :class="['node lvl1', value === root.id ? 'active' : '']">
            <text class="caret" @click="toggle(root.id)">{{ expanded[root.id] ? '▾' : '▸' }}</text>
            <text class="name" @click="pick(root, [])">{{ root.name }}</text>
          </view>
          <view v-if="expanded[root.id]" class="sub">
            <view v-for="mid in (root.children || [])" :key="mid.id">
              <view :class="['node lvl2', value === mid.id ? 'active' : '']">
                <text class="caret" @click="toggle(mid.id)">{{ expanded[mid.id] ? '▾' : '▸' }}</text>
                <text class="name" @click="pick(mid, [root.name])">{{ mid.name }}</text>
              </view>
              <view v-if="expanded[mid.id]" class="sub">
                <view v-for="leaf in (mid.children || [])" :key="leaf.id" :class="['node lvl3', value === leaf.id ? 'active' : '']">
                  <text class="name" @click="pick(leaf, [root.name, mid.name])">{{ leaf.name }}</text>
                </view>
              </view>
            </view>
          </view>
        </view>
      </scroll-view>
    </view>
  </up-popup>
</template>

<style scoped>
.drawer { width: 600rpx; height: 100vh; background: var(--eapp-bg); display: flex; flex-direction: column; }
.head { display: flex; align-items: center; justify-content: space-between; padding: 24rpx; padding-top: calc(24rpx + env(safe-area-inset-top)); background: var(--eapp-card); }
.title { font-size: 30rpx; font-weight: 700; }
.close { color: var(--eapp-text-muted); font-size: 26rpx; }
.body { flex: 1; padding: 12rpx; }
.node { display: flex; align-items: center; padding: 12rpx 14rpx; border-radius: 10rpx; gap: 10rpx; }
.node.active { background: var(--eapp-primary-soft); color: var(--eapp-primary); }
.lvl2 { padding-left: 32rpx; }
.lvl3 { padding-left: 56rpx; }
.caret { width: 32rpx; color: var(--eapp-text-muted); font-size: 22rpx; }
.name { font-size: 26rpx; flex: 1; }
.sub { padding-left: 4rpx; }
</style>

<script setup lang="ts">
import { onShow } from '@dcloudio/uni-app'
import { ref } from 'vue'
import { createCategory, deleteCategory, getCategoriesTree, updateCategory } from '@/api/category'

const tree = ref<any[]>([])
const expanded = ref<Record<number, boolean>>({})

async function load() { tree.value = (await getCategoriesTree()) || [] }

function toggle(id: number) { expanded.value[id] = !expanded.value[id] }

function addUnder(parent: any) {
  uni.showModal({ title: '新增分类', editable: true, placeholderText: '请输入名称', success: async (m) => {
    if (!m.confirm || !m.content) return
    await createCategory({ name: m.content, parent_id: Number(parent?.id || 0) })
    await load()
  } })
}

function rename(node: any) {
  uni.showModal({ title: '重命名', editable: true, content: node.name, success: async (m) => {
    if (!m.confirm || !m.content) return
    await updateCategory(Number(node.id), { name: m.content })
    await load()
  } })
}

function remove(node: any) {
  uni.showModal({ title: '删除', content: `确认删除分类「${node.name}」？`, success: async (m) => {
    if (!m.confirm) return
    await deleteCategory(Number(node.id))
    await load()
  } })
}

onShow(load)
</script>

<template>
  <view class="page">
    <up-button type="primary" plain class="add-top" @click="addUnder({ id: 0 })">+ 新增根分类</up-button>
    <view v-for="root in tree" :key="root.id" class="node-block">
      <view class="node">
        <text class="caret" @click="toggle(root.id)">{{ expanded[root.id] ? '▾' : '▸' }}</text>
        <text class="name">{{ root.name }}</text>
        <view class="ops">
          <text class="op" @click="addUnder(root)">+ 子</text>
          <text class="op" @click="rename(root)">改名</text>
          <text class="op danger" @click="remove(root)">删</text>
        </view>
      </view>
      <view v-if="expanded[root.id]" class="children">
        <view v-for="mid in (root.children || [])" :key="mid.id" class="node-block">
          <view class="node lvl2">
            <text class="caret" @click="toggle(mid.id)">{{ expanded[mid.id] ? '▾' : '▸' }}</text>
            <text class="name">{{ mid.name }}</text>
            <view class="ops">
              <text class="op" @click="addUnder(mid)">+ 子</text>
              <text class="op" @click="rename(mid)">改名</text>
              <text class="op danger" @click="remove(mid)">删</text>
            </view>
          </view>
          <view v-if="expanded[mid.id]" class="children">
            <view v-for="leaf in (mid.children || [])" :key="leaf.id" class="node lvl3">
              <text class="name">{{ leaf.name }}</text>
              <view class="ops">
                <text class="op" @click="rename(leaf)">改名</text>
                <text class="op danger" @click="remove(leaf)">删</text>
              </view>
            </view>
          </view>
        </view>
      </view>
    </view>
  </view>
</template>

<style scoped>
.page { min-height: 100vh; padding: 20rpx; background: var(--eapp-bg); }
.add-top { margin-bottom: 14rpx; }
.node-block { background: var(--eapp-card); border: 1px solid var(--eapp-border); border-radius: 14rpx; padding: 6rpx 10rpx; margin-bottom: 8rpx; }
.node { display: flex; align-items: center; padding: 12rpx 8rpx; gap: 10rpx; }
.node.lvl2 { padding-left: 30rpx; }
.node.lvl3 { padding-left: 60rpx; }
.caret { width: 28rpx; color: var(--eapp-text-muted); }
.name { flex: 1; font-size: 26rpx; }
.ops { display: flex; gap: 12rpx; }
.op { font-size: 24rpx; color: var(--eapp-primary); }
.op.danger { color: var(--eapp-danger); }
.children { padding: 4rpx 8rpx; }
</style>

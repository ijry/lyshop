<script setup lang="ts">
import { onLoad } from '@dcloudio/uni-app'
import { ref } from 'vue'
import { useDecorEditor } from '@/composables/useDecorEditor'
import DecorPropEditor from '@/components/decor/DecorPropEditor.vue'

const editor = useDecorEditor()
const showCompLib = ref(false)
const showPropEditor = ref(false)
const showVariantPicker = ref(false)

async function bootstrap() {
  await editor.loadVariants()
  if (editor.variants.value.length) {
    await editor.selectVariant(editor.currentVariantKey.value)
  }
}

function onVariantChange(e: any) {
  const v = e.value?.[0]
  if (v) editor.selectVariant(String(v.variant_key || 'default'))
}

function onAction(key: string) {
  if (key === 'copy') {
    uni.showModal({
      title: '复制副本',
      editable: true,
      placeholderText: '新副本 key',
      success: async (res) => {
        if (!res.confirm || !res.content) return
        await editor.copyDecorVariant({ source_key: editor.currentVariantKey.value, new_key: res.content.trim(), name: res.content.trim() })
        await editor.loadVariants()
      },
    })
  } else if (key === 'rename') {
    uni.showModal({
      title: '重命名',
      editable: true,
      placeholderText: '新名称',
      success: async (res) => {
        if (!res.confirm || !res.content) return
        await editor.renameDecorVariant(editor.currentVariantKey.value, res.content.trim())
        await editor.loadVariants()
      },
    })
  } else if (key === 'delete') {
    uni.showModal({
      title: '确认删除',
      content: '删除后不可恢复',
      success: async (res) => {
        if (!res.confirm) return
        await editor.deleteDecorVariant(editor.currentVariantKey.value)
        await editor.loadVariants()
        if (editor.variants.value.length) {
          await editor.selectVariant(String(editor.variants.value[0].variant_key))
        }
      },
    })
  }
}

function showActionSheet() {
  uni.showActionSheet({
    itemList: ['复制副本', '重命名', '删除'],
    success: (res) => {
      const keys = ['copy', 'rename', 'delete']
      onAction(keys[res.tapIndex])
    },
  })
}

function onAddComp(type: string) {
  editor.appendComp(type)
  showCompLib.value = false
}

function onTapComp(index: number) {
  editor.selectComp(index)
  showPropEditor.value = true
}

function onPropsUpdated(props: any) {
  const idx = editor.selectedIndex.value
  if (idx >= 0 && idx < editor.components.value.length) {
    editor.components.value[idx].props = { ...props }
  }
}

onLoad(bootstrap)
</script>

<template>
  <view class="page">
    <view class="toolbar">
      <view class="picker" @click="showVariantPicker = true">{{ editor.variants.value.find((v: any) => v.variant_key === editor.currentVariantKey.value)?.variant_name || '选择副本' }}</view>
      <up-picker :show="showVariantPicker" :columns="[editor.variants.value]" keyName="variant_name" @confirm="(e) => { onVariantChange(e); showVariantPicker = false }" @cancel="showVariantPicker = false" @close="showVariantPicker = false" />
      <up-button size="mini" plain @click="showActionSheet">操作</up-button>
    </view>

    <view class="comp-list">
      <view v-if="!editor.components.value.length" class="empty">暂无组件，点击下方添加</view>
      <view v-for="(comp, idx) in editor.components.value" :key="comp.id" :class="['comp-card', editor.selectedIndex.value === idx ? 'selected' : '']" @click="onTapComp(idx)">
        <view class="comp-info">
          <text class="comp-icon">{{ editor.componentLib.find((c: any) => c.type === comp.type)?.icon || '?' }}</text>
          <text class="comp-title">{{ editor.componentLib.find((c: any) => c.type === comp.type)?.title || comp.type }}</text>
        </view>
        <view class="comp-btns" @click.stop>
          <text class="btn" @click="editor.moveUp(idx)">↑</text>
          <text class="btn" @click="editor.moveDown(idx)">↓</text>
          <text class="btn del" @click="editor.removeComp(idx)">✕</text>
        </view>
      </view>
    </view>

    <view class="bottom-bar">
      <up-button size="small" @click="showCompLib = true">+ 添加组件</up-button>
      <view class="bar-right">
        <up-button size="small" type="primary" :loading="editor.saving.value" @click="editor.save">保存</up-button>
        <up-button size="small" type="success" @click="editor.publish">发布</up-button>
      </view>
    </view>

    <up-popup :show="showCompLib" mode="bottom" round="16" @close="showCompLib = false">
      <view class="popup-body">
        <view class="popup-title">添加组件</view>
        <view class="lib-grid">
          <view v-for="def in editor.componentLib" :key="def.type" class="lib-item" @click="onAddComp(def.type)">
            <text class="lib-icon">{{ def.icon }}</text>
            <text class="lib-label">{{ def.title }}</text>
          </view>
        </view>
      </view>
    </up-popup>

    <DecorPropEditor
      :show="showPropEditor"
      :component="editor.selectedIndex.value >= 0 ? editor.components.value[editor.selectedIndex.value] : null"
      @close="showPropEditor = false"
      @update="onPropsUpdated"
    />
  </view>
</template>

<style scoped>
.page { min-height: 100vh; background: var(--eapp-bg); padding: 20rpx; padding-bottom: 160rpx; box-sizing: border-box; }
.toolbar { display: flex; gap: 12rpx; align-items: center; margin-bottom: 16rpx; }
.picker { min-height: 76rpx; border: 1px solid var(--eapp-border); border-radius: 12rpx; padding: 0 20rpx; display: flex; align-items: center; color: var(--eapp-text); flex: 1; }
.comp-list { display: grid; gap: 12rpx; }
.comp-card { background: #fff; border: 1px solid var(--eapp-border); border-radius: 18rpx; padding: 18rpx; display: flex; align-items: center; justify-content: space-between; }
.comp-card.selected { border-color: var(--eapp-primary); background: #eff6ff; }
.comp-info { display: flex; align-items: center; gap: 12rpx; }
.comp-icon { font-size: 32rpx; }
.comp-title { font-size: 28rpx; font-weight: 600; }
.comp-btns { display: flex; gap: 16rpx; }
.btn { font-size: 30rpx; color: var(--eapp-text-muted); padding: 4rpx 8rpx; }
.btn.del { color: var(--eapp-danger); }
.bottom-bar { position: fixed; bottom: 0; left: 0; right: 0; background: #fff; border-top: 1px solid var(--eapp-border); padding: 16rpx 20rpx; display: flex; align-items: center; justify-content: space-between; z-index: 30; }
.bar-right { display: flex; gap: 10rpx; }
.empty { text-align: center; color: var(--eapp-text-muted); padding: 80rpx 0; }
.popup-body { padding: 24rpx; box-sizing: border-box; }
.popup-title { font-size: 30rpx; font-weight: 700; margin-bottom: 14rpx; }
.lib-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 16rpx; }
.lib-item { display: flex; flex-direction: column; align-items: center; gap: 8rpx; padding: 20rpx 0; background: var(--eapp-bg); border-radius: 14rpx; }
.lib-icon { font-size: 40rpx; }
.lib-label { font-size: 22rpx; color: var(--eapp-text-muted); }
</style>

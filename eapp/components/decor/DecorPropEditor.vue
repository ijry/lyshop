<script setup lang="ts">
import { reactive, watch } from 'vue'

const props = defineProps<{ show: boolean; component: any | null }>()
const emit = defineEmits<{ (e: 'close'): void; (e: 'update', props: any): void }>()

const form = reactive<any>({})

watch(() => props.component, (comp) => {
  if (comp?.props) Object.assign(form, JSON.parse(JSON.stringify(comp.props)))
}, { immediate: true, deep: true })

function save() {
  emit('update', { ...form })
  emit('close')
}
</script>

<template>
  <up-popup :show="show" mode="bottom" round="16" @close="$emit('close')">
    <view class="popup-body" v-if="component">
      <view class="popup-title">编辑属性 - {{ component.type }}</view>

      <!-- banner -->
      <template v-if="component.type === 'banner'">
        <up-input v-model="form.height" type="number" placeholder="高度" />
        <view class="mt"><text class="label">图片URL（逗号分隔）</text></view>
        <up-input v-model="form._images_str" placeholder="https://...，https://..." @blur="form.images = (form._images_str || '').split(',').filter(Boolean).map((u: string) => ({ url: u.trim() }))" />
      </template>

      <!-- category_nav -->
      <template v-if="component.type === 'category_nav'">
        <up-input v-model="form.style" placeholder="样式（grid / scroll）" />
        <view class="mt"><text class="label">分类ID（逗号分隔）</text></view>
        <up-input v-model="form._cat_ids_str" placeholder="1,2,3" @blur="form.category_ids = (form._cat_ids_str || '').split(',').map(Number).filter(Boolean)" />
      </template>

      <!-- product_grid -->
      <template v-if="component.type === 'product_grid'">
        <up-input v-model="form.source" placeholder="数据源（hot / new / category）" />
        <view class="mt" />
        <up-input v-model="form.limit" type="number" placeholder="数量" />
        <view class="mt" />
        <up-input v-model="form.columns" type="number" placeholder="列数" />
      </template>

      <!-- notice -->
      <template v-if="component.type === 'notice'">
        <up-input v-model="form.text" placeholder="公告内容" />
        <view class="mt" />
        <up-input v-model="form.color" placeholder="颜色 #f59e0b" />
      </template>

      <!-- image_ad -->
      <template v-if="component.type === 'image_ad'">
        <up-input v-model="form.image_url" placeholder="图片 URL" />
        <view class="mt" />
        <up-input v-model="form.link" placeholder="跳转链接" />
      </template>

      <!-- rich_text -->
      <template v-if="component.type === 'rich_text'">
        <up-input v-model="form.content" placeholder="HTML 内容" />
      </template>

      <!-- marketing_zone -->
      <template v-if="component.type === 'marketing_zone'">
        <up-input v-model="form.title" placeholder="标题" />
        <view class="mt" />
        <up-input v-model="form.type" placeholder="类型（seckill / group-buy / bargain）" />
      </template>

      <!-- spacer -->
      <template v-if="component.type === 'spacer'">
        <up-input v-model="form.height" type="number" placeholder="高度（rpx）" />
      </template>

      <view class="mt-lg" />
      <up-button type="primary" @click="save">确定</up-button>
    </view>
  </up-popup>
</template>

<style scoped>
.popup-body { padding: 24rpx; box-sizing: border-box; }
.popup-title { font-size: 30rpx; font-weight: 700; margin-bottom: 14rpx; }
.label { font-size: 24rpx; color: var(--eapp-text-muted); }
.mt { margin-top: 12rpx; }
.mt-lg { margin-top: 16rpx; }
</style>

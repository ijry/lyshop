<template>
  <div class="space-y-4">
    <ImageUpload :modelValue="modelValue.url" @update:modelValue="update('url', $event)" label="广告图片" />
    <LinkPicker :modelValue="modelValue.link" @update:modelValue="update('link', $event)" label="跳转链接" />
  </div>
</template>

<script setup lang="ts">
import type { ImageAdProps } from '@/types/decor'
import ImageUpload from '../widgets/ImageUpload.vue'
import LinkPicker from '../widgets/LinkPicker.vue'

const props = defineProps<{ modelValue: ImageAdProps }>()
const emit = defineEmits<{ 'update:modelValue': [value: ImageAdProps] }>()

function update<K extends keyof ImageAdProps>(key: K, value: ImageAdProps[K]) {
  emit('update:modelValue', { ...props.modelValue, [key]: value })
}
</script>

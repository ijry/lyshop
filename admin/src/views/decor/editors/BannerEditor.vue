<template>
  <div class="space-y-4">
    <!-- Height -->
    <div>
      <label class="block text-xs text-slate-500 mb-1.5">轮播高度 (px)</label>
      <div class="flex items-center gap-3">
        <input type="range" :value="modelValue.height" @input="updateField('height', Number(($event.target as HTMLInputElement).value))"
          min="150" max="600" step="10" class="flex-1 accent-blue-600" />
        <input type="number" :value="modelValue.height" @change="updateField('height', Number(($event.target as HTMLInputElement).value))"
          min="150" max="600" class="w-16 border border-slate-200 rounded-lg px-2 py-1.5 text-xs text-center focus:outline-none focus:border-blue-400" />
      </div>
    </div>

    <!-- Image list -->
    <div>
      <div class="flex items-center justify-between mb-2">
        <label class="text-xs text-slate-500">轮播图片 ({{ modelValue.images.length }})</label>
        <button @click="addImage" class="text-xs text-blue-600 hover:underline">+ 添加</button>
      </div>

      <div class="space-y-2">
        <div v-for="(img, idx) in modelValue.images" :key="idx"
          class="border border-slate-200 rounded-lg p-3 space-y-2.5">
          <div class="flex items-center justify-between">
            <span class="text-xs text-slate-400">图片 {{ idx + 1 }}</span>
            <div class="flex gap-0.5">
              <button @click="moveImage(idx, -1)" :disabled="idx === 0"
                class="px-1.5 py-1 text-xs text-slate-400 hover:text-slate-600 disabled:opacity-30">↑</button>
              <button @click="moveImage(idx, 1)" :disabled="idx === modelValue.images.length - 1"
                class="px-1.5 py-1 text-xs text-slate-400 hover:text-slate-600 disabled:opacity-30">↓</button>
              <button @click="removeImage(idx)"
                class="px-1.5 py-1 text-xs text-red-400 hover:text-red-600">×</button>
            </div>
          </div>

          <!-- Image preview + upload -->
          <div class="flex items-center gap-3">
            <div v-if="img.url"
              class="w-20 h-12 rounded-lg border border-slate-200 overflow-hidden bg-slate-50 shrink-0">
              <img :src="img.url" class="w-full h-full object-cover" />
            </div>
            <div v-else
              class="w-20 h-12 rounded-lg border-2 border-dashed border-slate-200 flex items-center justify-center text-slate-300 text-xs shrink-0">
              暂无
            </div>
            <div class="flex flex-col gap-1">
              <button @click="triggerUpload(idx)" :disabled="uploading === idx"
                class="px-2.5 py-1 text-xs bg-slate-100 rounded-lg hover:bg-slate-200 transition disabled:opacity-50">
                {{ uploading === idx ? '上传中...' : '上传' }}
              </button>
            </div>
          </div>

          <!-- Link picker -->
          <LinkPicker :modelValue="img.link" @update:modelValue="updateImage(idx, 'link', $event)" label="跳转链接" />
        </div>
      </div>

      <div v-if="!modelValue.images.length"
        class="text-center py-4 text-xs text-slate-300 border-2 border-dashed border-slate-200 rounded-lg">
        暂无轮播图，点击上方添加
      </div>
    </div>

    <input ref="fileInput" type="file" accept="image/*" class="hidden" @change="onFileChange" />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { BannerProps } from '@/types/decor'
import { uploadFile } from '@/api/plugins'
import LinkPicker from '../widgets/LinkPicker.vue'

const props = defineProps<{ modelValue: BannerProps }>()
const emit = defineEmits<{ 'update:modelValue': [value: BannerProps] }>()

const fileInput = ref<HTMLInputElement>()
const uploading = ref<number | null>(null)
let uploadTargetIdx = -1

function updateField<K extends keyof BannerProps>(key: K, value: BannerProps[K]) {
  emit('update:modelValue', { ...props.modelValue, [key]: value })
}

function addImage() {
  emit('update:modelValue', {
    ...props.modelValue,
    images: [...props.modelValue.images, { url: '', link: '' }],
  })
}

function removeImage(idx: number) {
  const images = [...props.modelValue.images]
  images.splice(idx, 1)
  emit('update:modelValue', { ...props.modelValue, images })
}

function moveImage(idx: number, delta: number) {
  const target = idx + delta
  if (target < 0 || target >= props.modelValue.images.length) return
  const images = [...props.modelValue.images]
  ;[images[idx], images[target]] = [images[target], images[idx]]
  emit('update:modelValue', { ...props.modelValue, images })
}

function updateImage(idx: number, field: 'url' | 'link', value: string) {
  const images = props.modelValue.images.map((img, i) =>
    i === idx ? { ...img, [field]: value } : img
  )
  emit('update:modelValue', { ...props.modelValue, images })
}

function triggerUpload(idx: number) {
  uploadTargetIdx = idx
  fileInput.value?.click()
}

async function onFileChange(e: Event) {
  const file = (e.target as HTMLInputElement).files?.[0]
  if (!file || uploadTargetIdx < 0) return
  const idx = uploadTargetIdx
  uploading.value = idx
  try {
    const result: any = await uploadFile(file)
    if (result?.url) {
      updateImage(idx, 'url', result.url)
    }
  } finally {
    uploading.value = null
    if (fileInput.value) fileInput.value.value = ''
  }
}
</script>

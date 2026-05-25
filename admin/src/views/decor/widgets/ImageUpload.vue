<template>
  <div class="space-y-2">
    <label v-if="label" class="block text-xs text-slate-500">{{ label }}</label>
    <div class="flex items-center gap-3">
      <div v-if="modelValue"
        class="w-16 h-16 rounded-lg border border-slate-200 overflow-hidden bg-slate-50 shrink-0">
        <img :src="modelValue" class="w-full h-full object-cover" />
      </div>
      <div v-else
        class="w-16 h-16 rounded-lg border-2 border-dashed border-slate-200 flex items-center justify-center text-slate-300 text-xs shrink-0">
        暂无
      </div>
      <div class="flex flex-col gap-1.5">
        <button @click="triggerUpload" :disabled="uploading"
          class="px-3 py-1.5 text-xs bg-slate-100 rounded-lg hover:bg-slate-200 transition disabled:opacity-50">
          {{ uploading ? '上传中...' : '上传图片' }}
        </button>
        <button v-if="modelValue" @click="emit('update:modelValue', '')"
          class="px-3 py-1.5 text-xs text-red-500 bg-red-50 rounded-lg hover:bg-red-100 transition">
          清除
        </button>
      </div>
    </div>
    <input ref="fileInput" type="file" accept="image/*" class="hidden" @change="onFileChange" />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { uploadFile } from '@/api/plugins'

defineProps<{
  modelValue: string
  label?: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const fileInput = ref<HTMLInputElement>()
const uploading = ref(false)

function triggerUpload() {
  fileInput.value?.click()
}

async function onFileChange(e: Event) {
  const file = (e.target as HTMLInputElement).files?.[0]
  if (!file) return
  uploading.value = true
  try {
    const result: any = await uploadFile(file)
    if (result?.url) {
      emit('update:modelValue', result.url)
    }
  } finally {
    uploading.value = false
    if (fileInput.value) fileInput.value.value = ''
  }
}
</script>

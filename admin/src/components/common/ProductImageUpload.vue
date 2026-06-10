<template>
  <div class="product-image-upload">
    <div class="image-preview" :class="{ 'image-preview--empty': !modelValue }">
      <img v-if="modelValue" :src="modelValue" :alt="previewAlt" />
      <span v-else>{{ emptyText }}</span>
    </div>

    <div class="image-actions">
      <div class="action-row">
        <el-upload
          :show-file-list="false"
          :http-request="handleUpload"
          :before-upload="beforeUpload"
          accept="image/*"
        >
          <button type="button" class="upload-button" :disabled="uploading">
            {{ uploading ? uploadingText : uploadText }}
          </button>
        </el-upload>
        <button
          v-if="modelValue"
          type="button"
          class="clear-button"
          :disabled="uploading"
          @click="emit('update:modelValue', '')"
        >
          {{ clearText }}
        </button>
      </div>

      <input
        :value="modelValue"
        class="url-input"
        :placeholder="placeholder"
        @input="emit('update:modelValue', ($event.target as HTMLInputElement).value)"
      />
      <p v-if="errorMessage" class="error-text">{{ errorMessage }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { UploadRawFile, UploadRequestOptions } from 'element-plus'
import { ElUpload } from 'element-plus'
import 'element-plus/es/components/upload/style/css'
import { uploadFile } from '@/api/plugins'

const props = withDefaults(defineProps<{
  modelValue: string
  placeholder?: string
  emptyText?: string
  uploadText?: string
  uploadingText?: string
  clearText?: string
  previewAlt?: string
  invalidFileText?: string
  uploadResultMissingText?: string
  uploadFailedText?: string
}>(), {
  placeholder: 'https://...',
  emptyText: '暂无图片',
  uploadText: '上传图片',
  uploadingText: '上传中...',
  clearText: '清除',
  previewAlt: '图片预览',
  invalidFileText: '请选择图片文件',
  uploadResultMissingText: '上传结果缺少图片地址',
  uploadFailedText: '图片上传失败',
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const uploading = ref(false)
const errorMessage = ref('')

function beforeUpload(file: UploadRawFile) {
  errorMessage.value = ''
  if (!file.type.startsWith('image/')) {
    errorMessage.value = props.invalidFileText
    return false
  }
  return true
}

async function handleUpload(options: UploadRequestOptions) {
  const file = options.file
  uploading.value = true
  errorMessage.value = ''
  try {
    const result: any = await uploadFile(file)
    if (!result?.url) {
      throw new Error(props.uploadResultMissingText)
    }
    emit('update:modelValue', result.url)
    options.onSuccess?.(result)
  } catch (error: any) {
    const message = error?.message || props.uploadFailedText
    errorMessage.value = message
    options.onError?.(error)
  } finally {
    uploading.value = false
  }
}

</script>

<style scoped>
.product-image-upload {
  display: grid;
  grid-template-columns: 104px minmax(0, 1fr);
  gap: 12px;
  align-items: start;
}

.image-preview {
  width: 104px;
  height: 78px;
  overflow: hidden;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  background: #f8fafc;
}

.image-preview img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

.image-preview--empty {
  display: flex;
  align-items: center;
  justify-content: center;
  border-style: dashed;
  color: #94a3b8;
  font-size: 12px;
}

.image-actions {
  min-width: 0;
}

.action-row {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 8px;
}

.upload-button,
.clear-button {
  min-height: 32px;
  border-radius: 8px;
  padding: 0 12px;
  font-size: 12px;
  cursor: pointer;
  transition: background-color 0.2s, color 0.2s, border-color 0.2s;
}

.upload-button {
  border: 1px solid #bfdbfe;
  background: #eff6ff;
  color: #2563eb;
}

.upload-button:hover:not(:disabled) {
  border-color: #93c5fd;
  background: #dbeafe;
}

.clear-button {
  border: 1px solid #fecaca;
  background: #fef2f2;
  color: #dc2626;
}

.clear-button:hover:not(:disabled) {
  background: #fee2e2;
}

.upload-button:disabled,
.clear-button:disabled {
  cursor: not-allowed;
  opacity: 0.55;
}

.url-input {
  width: 100%;
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  padding: 9px 12px;
  font-size: 13px;
  color: #334155;
  outline: none;
}

.url-input:focus {
  border-color: #60a5fa;
}

.error-text {
  margin-top: 6px;
  color: #dc2626;
  font-size: 12px;
}

@media (max-width: 640px) {
  .product-image-upload {
    grid-template-columns: 1fr;
  }

  .image-preview {
    width: 100%;
    height: 160px;
  }
}
</style>

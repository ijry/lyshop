<template>
  <div class="space-y-2">
    <label class="block text-xs text-slate-500">{{ $t('decor.richText.content') }}</label>
    <div v-if="editorReady" class="overflow-hidden rounded-lg border border-slate-200 bg-white">
      <component
        :is="toolbarComponent"
        class="border-b border-slate-200"
        :editor="editorRef"
        :default-config="toolbarConfig"
        mode="default"
      />
      <component
        :is="editorComponent"
        class="min-h-[320px] text-sm"
        :default-config="editorConfig"
        mode="default"
        :model-value="html"
        @on-created="handleCreated"
        @on-change="handleChange"
      />
    </div>
    <div v-else class="flex min-h-[320px] items-center justify-center rounded-lg border border-slate-200 bg-slate-50 text-sm text-slate-400">编辑器加载中...</div>
    <p class="text-xs text-slate-400">{{ $t('decor.richText.placeholder') }}</p>
  </div>
</template>

<script setup lang="ts">
import { computed, markRaw, onBeforeUnmount, onMounted, shallowRef } from 'vue'
import { useI18n } from 'vue-i18n'
import type { Component } from 'vue'
import type { IEditorConfig, IToolbarConfig, IDomEditor } from '@wangeditor-next/editor'
import type { RichTextProps } from '@/types/decor'
import { uploadFile } from '@/api/plugins'

const props = defineProps<{ modelValue: RichTextProps }>()
const emit = defineEmits<{ 'update:modelValue': [value: RichTextProps] }>()

const { t } = useI18n()
const editorRef = shallowRef<IDomEditor>()
const editorComponent = shallowRef<Component>()
const toolbarComponent = shallowRef<Component>()
const editorReady = computed(() => !!editorComponent.value && !!toolbarComponent.value)

const html = computed(() => props.modelValue?.content || '')

const toolbarConfig: Partial<IToolbarConfig> = {
  toolbarKeys: [
    'headerSelect',
    'bold',
    'italic',
    'underline',
    'through',
    '|',
    'color',
    'bgColor',
    '|',
    'bulletedList',
    'numberedList',
    '|',
    'insertLink',
    'uploadImage',
    'insertTable',
    'codeBlock',
    'blockquote',
    '|',
    'undo',
    'redo',
  ],
}

const editorConfig: Partial<IEditorConfig> = {
  placeholder: t('decor.richText.placeholder'),
  MENU_CONF: {
    uploadImage: {
      async customUpload(file: File, insertFn: (url: string, alt?: string, href?: string) => void) {
        const result: any = await uploadFile(file)
        const url = String(result?.url || '')
        if (!url) return
        insertFn(url, file.name, url)
      },
    },
  },
}

function handleCreated(editor: IDomEditor) {
  editorRef.value = editor
}

function handleChange(editor: IDomEditor) {
  emit('update:modelValue', { content: editor.getHtml() })
}

onMounted(async () => {
  const [{ Editor, Toolbar }] = await Promise.all([
    import('@wangeditor-next/editor-for-vue'),
    import('@wangeditor-next/editor/dist/css/style.css'),
  ])
  editorComponent.value = markRaw(Editor)
  toolbarComponent.value = markRaw(Toolbar)
})

onBeforeUnmount(() => {
  editorRef.value?.destroy()
  editorRef.value = undefined
})
</script>

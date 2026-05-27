<template>
  <div>
    <div v-if="loading" class="max-w-7xl mx-auto px-6 py-20 text-center text-gray-400">加载中...</div>
    <DecorRenderer v-else :components="components" />
    <PresetSwitcher />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { get } from '@/api/request'
import DecorRenderer from '@/components/decor/DecorRenderer.vue'
import PresetSwitcher from '@/components/PresetSwitcher.vue'

const components = ref<any[]>([])
const loading = ref(true)

onMounted(async () => {
  try {
    const data = await get<any>('/api/v1/pc/decor')
    components.value = data?.components || []
  } catch { /* fallback empty */ }
  loading.value = false
})
</script>

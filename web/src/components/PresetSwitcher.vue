<template>
  <div v-if="show" class="fixed left-0 right-0 bottom-0 z-[999]" style="pointer-events: none;">
    <div class="flex justify-center pb-5" style="pointer-events: auto;">
      <div v-if="!expanded" @click="expanded = true"
        class="px-3 py-1.5 rounded-full text-xs shadow-lg cursor-pointer select-none"
        style="background: rgba(0,0,0,0.7); color: #fff; backdrop-filter: blur(8px);">
        {{ currentName }} ▾
      </div>
      <div v-else class="mx-4 rounded-xl shadow-2xl overflow-hidden max-w-xs w-full"
        style="background: rgba(255,255,255,0.96); backdrop-filter: blur(12px);">
        <div class="px-3 pt-2.5 pb-1.5 flex items-center justify-between">
          <span class="text-xs text-gray-400">{{ $t('presetSwitcher.switchIndustry') }}</span>
          <span class="text-xs text-gray-400 px-1 cursor-pointer" @click="expanded = false">✕</span>
        </div>
        <div class="flex flex-wrap px-2 pb-2.5 gap-1.5">
          <span v-for="p in presetList" :key="p.key" @click="switchPreset(p.key)"
            class="px-3 py-1.5 rounded-full text-xs cursor-pointer select-none transition-colors"
            :class="p.key === activeKey ? 'bg-blue-600 text-white' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'">
            {{ p.name }}
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

const show = import.meta.env.VITE_MOCK === 'true'
const expanded = ref(false)
const presetList = ref<Array<{ key: string; name: string }>>([])
const activeKey = ref('')
const currentName = ref('')

onMounted(async () => {
  if (!show) return
  const { listPresets, getPresetKey } = await import('../../../app/mock/presets/index')
  presetList.value = listPresets()
  activeKey.value = getPresetKey()
  currentName.value = presetList.value.find(p => p.key === activeKey.value)?.name || activeKey.value
})

function switchPreset(key: string) {
  if (key === activeKey.value) {
    expanded.value = false
    return
  }
  const url = new URL(window.location.href)
  if (key === 'mall') {
    url.searchParams.delete('demo')
  } else {
    url.searchParams.set('demo', key)
  }
  window.location.href = url.toString()
}
</script>

<template>
  <div class="space-y-4">
    <div>
      <label class="block text-xs text-slate-500 mb-1.5">背景模式</label>
      <select
        :value="modelValue.background.mode"
        class="w-full border border-slate-200 rounded-lg px-2.5 py-1.5 text-xs focus:outline-none focus:border-blue-400"
        @change="onModeChange(($event.target as HTMLSelectElement).value)"
      >
        <option value="solid">纯色</option>
        <option value="gradient">渐变</option>
        <option value="image">背景图</option>
      </select>
    </div>

    <div v-if="modelValue.background.mode === 'solid'" class="space-y-2">
      <ColorInput
        :modelValue="modelValue.background.solidColor || '#f8fafc'"
        label="背景色"
        @update:modelValue="updateBackground({ solidColor: $event })"
      />
    </div>

    <div v-else-if="modelValue.background.mode === 'gradient'" class="space-y-2">
      <div>
        <label class="block text-xs text-slate-500 mb-1.5">渐变角度</label>
        <input
          type="number"
          :value="modelValue.background.gradient?.angle ?? 135"
          class="w-full border border-slate-200 rounded-lg px-2.5 py-1.5 text-xs focus:outline-none focus:border-blue-400"
          @input="updateGradient('angle', Number(($event.target as HTMLInputElement).value || 0))"
        />
      </div>
      <div class="grid grid-cols-2 gap-3">
        <ColorInput
          :modelValue="gradientStop(0).color"
          label="起始色"
          @update:modelValue="updateStop(0, 'color', $event)"
        />
        <ColorInput
          :modelValue="gradientStop(1).color"
          label="结束色"
          @update:modelValue="updateStop(1, 'color', $event)"
        />
      </div>
      <div class="grid grid-cols-2 gap-3">
        <div>
          <label class="block text-xs text-slate-500 mb-1.5">起始点(%)</label>
          <input
            type="number"
            min="0"
            max="100"
            :value="gradientStop(0).position"
            class="w-full border border-slate-200 rounded-lg px-2.5 py-1.5 text-xs focus:outline-none focus:border-blue-400"
            @input="updateStop(0, 'position', clampPos(($event.target as HTMLInputElement).value))"
          />
        </div>
        <div>
          <label class="block text-xs text-slate-500 mb-1.5">结束点(%)</label>
          <input
            type="number"
            min="0"
            max="100"
            :value="gradientStop(1).position"
            class="w-full border border-slate-200 rounded-lg px-2.5 py-1.5 text-xs focus:outline-none focus:border-blue-400"
            @input="updateStop(1, 'position', clampPos(($event.target as HTMLInputElement).value))"
          />
        </div>
      </div>
    </div>

    <div v-else class="space-y-2">
      <ImageUpload
        :modelValue="modelValue.background.image?.url || ''"
        label="背景图"
        @update:modelValue="updateImage('url', $event)"
      />
      <div class="grid grid-cols-2 gap-3">
        <div>
          <label class="block text-xs text-slate-500 mb-1.5">尺寸</label>
          <select
            :value="modelValue.background.image?.size || 'cover'"
            class="w-full border border-slate-200 rounded-lg px-2.5 py-1.5 text-xs focus:outline-none focus:border-blue-400"
            @change="updateImage('size', ($event.target as HTMLSelectElement).value)"
          >
            <option value="cover">cover</option>
            <option value="contain">contain</option>
            <option value="auto">auto</option>
            <option value="custom">custom</option>
          </select>
        </div>
        <div>
          <label class="block text-xs text-slate-500 mb-1.5">定位</label>
          <input
            :value="modelValue.background.image?.position || 'center top'"
            class="w-full border border-slate-200 rounded-lg px-2.5 py-1.5 text-xs focus:outline-none focus:border-blue-400"
            @input="updateImage('position', ($event.target as HTMLInputElement).value)"
          />
        </div>
      </div>
      <div v-if="(modelValue.background.image?.size || 'cover') === 'custom'">
        <label class="block text-xs text-slate-500 mb-1.5">自定义尺寸</label>
        <input
          :value="modelValue.background.image?.customSize || '100% auto'"
          class="w-full border border-slate-200 rounded-lg px-2.5 py-1.5 text-xs focus:outline-none focus:border-blue-400"
          @input="updateImage('customSize', ($event.target as HTMLInputElement).value)"
        />
      </div>
      <div class="grid grid-cols-2 gap-3">
        <div>
          <label class="block text-xs text-slate-500 mb-1.5">平铺</label>
          <select
            :value="modelValue.background.image?.repeat || 'no-repeat'"
            class="w-full border border-slate-200 rounded-lg px-2.5 py-1.5 text-xs focus:outline-none focus:border-blue-400"
            @change="updateImage('repeat', ($event.target as HTMLSelectElement).value)"
          >
            <option value="no-repeat">no-repeat</option>
            <option value="repeat">repeat</option>
            <option value="repeat-x">repeat-x</option>
            <option value="repeat-y">repeat-y</option>
          </select>
        </div>
        <div>
          <label class="block text-xs text-slate-500 mb-1.5">滚动</label>
          <select
            :value="modelValue.background.image?.attachment || 'scroll'"
            class="w-full border border-slate-200 rounded-lg px-2.5 py-1.5 text-xs focus:outline-none focus:border-blue-400"
            @change="updateImage('attachment', ($event.target as HTMLSelectElement).value)"
          >
            <option value="scroll">scroll</option>
            <option value="fixed">fixed</option>
          </select>
        </div>
      </div>
    </div>

    <div class="border-t border-slate-100 pt-3 space-y-2">
      <label class="flex items-center justify-between text-xs text-slate-500">
        <span>遮罩</span>
        <input
          type="checkbox"
          :checked="modelValue.background.overlay?.enabled || false"
          @change="updateOverlay('enabled', ($event.target as HTMLInputElement).checked)"
        />
      </label>
      <div v-if="modelValue.background.overlay?.enabled" class="grid grid-cols-2 gap-3">
        <ColorInput
          :modelValue="modelValue.background.overlay?.color || '#000000'"
          label="遮罩色"
          @update:modelValue="updateOverlay('color', $event)"
        />
        <div>
          <label class="block text-xs text-slate-500 mb-1.5">透明度(0~1)</label>
          <input
            type="number"
            min="0"
            max="1"
            step="0.05"
            :value="modelValue.background.overlay?.opacity ?? 0.2"
            class="w-full border border-slate-200 rounded-lg px-2.5 py-1.5 text-xs focus:outline-none focus:border-blue-400"
            @input="updateOverlay('opacity', clampOpacity(($event.target as HTMLInputElement).value))"
          />
        </div>
      </div>
    </div>

    <div class="border-t border-slate-100 pt-3 space-y-2">
      <p class="text-xs text-slate-500">内容布局</p>
      <div class="grid grid-cols-3 gap-3">
        <div>
          <label class="block text-xs text-slate-500 mb-1.5">最大宽度</label>
          <input
            type="number"
            :value="modelValue.content.maxWidth"
            class="w-full border border-slate-200 rounded-lg px-2.5 py-1.5 text-xs focus:outline-none focus:border-blue-400"
            @input="updateContent('maxWidth', Number(($event.target as HTMLInputElement).value || 0))"
          />
        </div>
        <div>
          <label class="block text-xs text-slate-500 mb-1.5">左右留白</label>
          <input
            type="number"
            :value="modelValue.content.gutterX"
            class="w-full border border-slate-200 rounded-lg px-2.5 py-1.5 text-xs focus:outline-none focus:border-blue-400"
            @input="updateContent('gutterX', Number(($event.target as HTMLInputElement).value || 0))"
          />
        </div>
        <div>
          <label class="block text-xs text-slate-500 mb-1.5">模块间距</label>
          <input
            type="number"
            :value="modelValue.content.sectionGap"
            class="w-full border border-slate-200 rounded-lg px-2.5 py-1.5 text-xs focus:outline-none focus:border-blue-400"
            @input="updateContent('sectionGap', Number(($event.target as HTMLInputElement).value || 0))"
          />
        </div>
      </div>
    </div>

    <div class="border-t border-slate-100 pt-3 space-y-2">
      <p class="text-xs text-slate-500">默认外观</p>
      <div class="grid grid-cols-2 gap-3">
        <div>
          <label class="block text-xs text-slate-500 mb-1.5">默认圆角</label>
          <input
            type="number"
            :value="modelValue.surface.radius"
            class="w-full border border-slate-200 rounded-lg px-2.5 py-1.5 text-xs focus:outline-none focus:border-blue-400"
            @input="updateSurface('radius', Number(($event.target as HTMLInputElement).value || 0))"
          />
        </div>
        <div>
          <label class="block text-xs text-slate-500 mb-1.5">默认阴影</label>
          <select
            :value="modelValue.surface.shadow"
            class="w-full border border-slate-200 rounded-lg px-2.5 py-1.5 text-xs focus:outline-none focus:border-blue-400"
            @change="updateSurface('shadow', ($event.target as HTMLSelectElement).value)"
          >
            <option value="none">none</option>
            <option value="sm">sm</option>
            <option value="md">md</option>
            <option value="lg">lg</option>
          </select>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import ColorInput from './ColorInput.vue'
import ImageUpload from './ImageUpload.vue'

const props = defineProps<{ modelValue: any }>()
const emit = defineEmits<{ 'update:modelValue': [value: any] }>()

function update(next: any) {
  emit('update:modelValue', next)
}

function gradientStop(index: number) {
  const stops = Array.isArray(props.modelValue.background.gradient?.stops)
    ? props.modelValue.background.gradient.stops
    : []
  return stops[index] || { color: index === 0 ? '#f8fafc' : '#eef2ff', position: index === 0 ? 0 : 100 }
}

function clampPos(value: string) {
  const n = Number(value || 0)
  return Math.min(100, Math.max(0, n))
}

function clampOpacity(value: string) {
  const n = Number(value || 0)
  return Math.min(1, Math.max(0, n))
}

function onModeChange(mode: string) {
  updateBackground({ mode })
}

function updateBackground(partial: Record<string, any>) {
  update({
    ...props.modelValue,
    background: {
      ...props.modelValue.background,
      ...partial,
    },
  })
}

function updateGradient(key: string, value: any) {
  const gradient = {
    angle: 135,
    stops: [gradientStop(0), gradientStop(1)],
    ...(props.modelValue.background.gradient || {}),
    [key]: value,
  }
  updateBackground({ gradient })
}

function updateStop(index: number, key: string, value: any) {
  const stops = [gradientStop(0), gradientStop(1)]
  stops[index] = { ...stops[index], [key]: value }
  updateGradient('stops', stops)
}

function updateImage(key: string, value: any) {
  const image = {
    url: '',
    size: 'cover',
    customSize: '100% auto',
    position: 'center top',
    repeat: 'no-repeat',
    attachment: 'scroll',
    ...(props.modelValue.background.image || {}),
    [key]: value,
  }
  updateBackground({ image })
}

function updateOverlay(key: string, value: any) {
  const overlay = {
    enabled: false,
    color: '#000000',
    opacity: 0.2,
    ...(props.modelValue.background.overlay || {}),
    [key]: value,
  }
  updateBackground({ overlay })
}

function updateContent(key: string, value: number) {
  update({
    ...props.modelValue,
    content: {
      ...props.modelValue.content,
      [key]: Math.max(0, Number(value || 0)),
    },
  })
}

function updateSurface(key: string, value: any) {
  update({
    ...props.modelValue,
    surface: {
      ...props.modelValue.surface,
      [key]: key === 'radius' ? Math.max(0, Number(value || 0)) : value,
    },
  })
}
</script>

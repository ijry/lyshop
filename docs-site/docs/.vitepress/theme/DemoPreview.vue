<template>
  <Teleport to="body">
    <!-- Floating trigger button (minimized state) -->
    <div
      v-show="!expanded"
      class="demo-fab"
      @click="expanded = true"
    >
      <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <rect x="5" y="2" width="14" height="20" rx="2" ry="2"/>
        <line x1="12" y1="18" x2="12" y2="18"/>
      </svg>
      <span class="demo-fab-text">演示</span>
    </div>

    <!-- Expanded phone preview -->
    <div
      v-show="expanded"
      class="demo-phone"
      :style="phoneStyle"
    >
      <!-- Title bar (drag handle) -->
      <div
        class="demo-phone-bar"
        @mousedown.prevent="startDrag"
        @touchstart.prevent="startDrag"
      >
        <span class="demo-phone-bar-title">LYShop 演示</span>
        <div class="demo-phone-bar-actions">
          <button class="demo-phone-btn" @click="refreshDemo" title="刷新">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M23 4v6h-6"/><path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10"/></svg>
          </button>
          <button class="demo-phone-btn" @click="expanded = false" title="最小化">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="5" y1="12" x2="19" y2="12"/></svg>
          </button>
        </div>
      </div>

      <!-- iPhone-style notch -->
      <div class="demo-phone-notch">
        <div class="demo-phone-notch-inner"></div>
      </div>

      <!-- iframe -->
      <div class="demo-phone-screen">
        <iframe
          ref="iframeRef"
          :src="demoUrl"
          class="demo-phone-iframe"
          frameborder="0"
          allow="clipboard-write"
        ></iframe>
      </div>

      <!-- Home indicator -->
      <div class="demo-phone-home">
        <div class="demo-phone-home-bar"></div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useData } from 'vitepress'

const { site } = useData()
const expanded = ref(false)
const iframeRef = ref<HTMLIFrameElement>()

// Position state
const posX = ref(0)
const posY = ref(0)
const positioned = ref(false)

const demoUrl = computed(() => {
  const base = site.value.base || '/'
  return `${base}demo/index.html`
})

const phoneStyle = computed(() => {
  if (!positioned.value) return {}
  return {
    position: 'fixed',
    left: `${posX.value}px`,
    top: `${posY.value}px`,
    right: 'auto',
    bottom: 'auto',
  }
})

// Drag logic
let dragging = false
let startX = 0
let startY = 0
let origX = 0
let origY = 0

function startDrag(e: MouseEvent | TouchEvent) {
  dragging = true
  const point = 'touches' in e ? e.touches[0] : e
  const el = (e.currentTarget as HTMLElement).parentElement!
  const rect = el.getBoundingClientRect()

  if (!positioned.value) {
    posX.value = rect.left
    posY.value = rect.top
    positioned.value = true
  }

  startX = point.clientX
  startY = point.clientY
  origX = posX.value
  origY = posY.value

  document.addEventListener('mousemove', onDrag)
  document.addEventListener('mouseup', stopDrag)
  document.addEventListener('touchmove', onDrag, { passive: false })
  document.addEventListener('touchend', stopDrag)
}

function onDrag(e: MouseEvent | TouchEvent) {
  if (!dragging) return
  e.preventDefault()
  const point = 'touches' in e ? e.touches[0] : e
  posX.value = origX + (point.clientX - startX)
  posY.value = origY + (point.clientY - startY)

  // Clamp to viewport
  posX.value = Math.max(0, Math.min(posX.value, window.innerWidth - 390))
  posY.value = Math.max(0, Math.min(posY.value, window.innerHeight - 100))
}

function stopDrag() {
  dragging = false
  document.removeEventListener('mousemove', onDrag)
  document.removeEventListener('mouseup', stopDrag)
  document.removeEventListener('touchmove', onDrag)
  document.removeEventListener('touchend', stopDrag)
}

function refreshDemo() {
  if (iframeRef.value) {
    iframeRef.value.src = iframeRef.value.src
  }
}

onUnmounted(() => {
  stopDrag()
})
</script>

<style scoped>
/* Floating action button */
.demo-fab {
  position: fixed;
  right: 24px;
  bottom: 24px;
  z-index: 999;
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 12px 18px;
  background: linear-gradient(135deg, #dc2626 0%, #ef4444 100%);
  color: #fff;
  border-radius: 999px;
  cursor: pointer;
  box-shadow: 0 4px 20px rgba(220, 38, 38, 0.4);
  transition: transform 0.2s, box-shadow 0.2s;
  user-select: none;
  font-size: 14px;
  font-weight: 500;
}
.demo-fab:hover {
  transform: scale(1.05);
  box-shadow: 0 6px 28px rgba(220, 38, 38, 0.5);
}
.demo-fab-text {
  line-height: 1;
}

/* Phone frame */
.demo-phone {
  position: fixed;
  right: 24px;
  bottom: 24px;
  z-index: 1000;
  width: 375px;
  height: 750px;
  background: #000;
  border-radius: 40px;
  box-shadow: 0 8px 40px rgba(0, 0, 0, 0.3), 0 0 0 2px rgba(255,255,255,0.1) inset;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  animation: demo-phone-in 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}

@keyframes demo-phone-in {
  from { opacity: 0; transform: scale(0.85) translateY(20px); }
  to   { opacity: 1; transform: scale(1) translateY(0); }
}

/* Drag bar */
.demo-phone-bar {
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 16px;
  background: #111;
  cursor: grab;
  user-select: none;
}
.demo-phone-bar:active {
  cursor: grabbing;
}
.demo-phone-bar-title {
  font-size: 12px;
  color: #999;
  font-weight: 500;
}
.demo-phone-bar-actions {
  display: flex;
  gap: 8px;
}
.demo-phone-btn {
  background: none;
  border: none;
  color: #888;
  cursor: pointer;
  padding: 4px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  transition: background 0.15s, color 0.15s;
}
.demo-phone-btn:hover {
  background: rgba(255,255,255,0.1);
  color: #fff;
}

/* Notch */
.demo-phone-notch {
  height: 28px;
  background: #000;
  display: flex;
  align-items: center;
  justify-content: center;
}
.demo-phone-notch-inner {
  width: 120px;
  height: 24px;
  background: #111;
  border-radius: 0 0 16px 16px;
}

/* Screen */
.demo-phone-screen {
  flex: 1;
  min-height: 0;
  background: #fff;
  margin: 0 8px;
  border-radius: 4px;
  overflow: hidden;
}
.demo-phone-iframe {
  width: 100%;
  height: 100%;
  border: none;
}

/* Home indicator */
.demo-phone-home {
  height: 28px;
  background: #000;
  display: flex;
  align-items: center;
  justify-content: center;
}
.demo-phone-home-bar {
  width: 134px;
  height: 5px;
  background: #333;
  border-radius: 999px;
}

/* Mobile: smaller phone */
@media (max-width: 640px) {
  .demo-fab {
    right: 12px;
    bottom: 12px;
    padding: 10px 14px;
    font-size: 12px;
  }
  .demo-phone {
    right: 8px;
    bottom: 8px;
    width: 320px;
    height: 620px;
    border-radius: 32px;
  }
  .demo-phone-screen {
    margin: 0 6px;
  }
}
</style>

import { ref, computed } from 'vue'

export type ThemeMode = 'light' | 'dark' | 'auto'

const STORAGE_KEY = 'app_theme_mode'

const themeMode = ref<ThemeMode>((uni.getStorageSync(STORAGE_KEY) as ThemeMode) || 'auto')

function getSystemTheme(): 'light' | 'dark' {
  try {
    const info = uni.getSystemInfoSync()
    if ((info as any).theme === 'dark') return 'dark'
  } catch {}
  if (typeof window !== 'undefined' && window.matchMedia?.('(prefers-color-scheme: dark)').matches) {
    return 'dark'
  }
  return 'light'
}

const effectiveTheme = computed<'light' | 'dark'>(() => {
  if (themeMode.value === 'auto') return getSystemTheme()
  return themeMode.value
})

function applyTheme() {
  const theme = effectiveTheme.value
  if (typeof document !== 'undefined') {
    document.documentElement.setAttribute('data-up-theme', theme)
  }
  // uni-app page element
  try {
    const pages = document.querySelectorAll('page')
    pages.forEach(p => p.setAttribute('data-up-theme', theme))
  } catch {}
}

function setTheme(mode: ThemeMode) {
  themeMode.value = mode
  uni.setStorageSync(STORAGE_KEY, mode)
  applyTheme()
}

function toggleTheme() {
  const next = effectiveTheme.value === 'light' ? 'dark' : 'light'
  setTheme(next)
}

function initTheme() {
  applyTheme()
  // Listen for system theme changes (auto mode)
  if (typeof window !== 'undefined') {
    window.matchMedia?.('(prefers-color-scheme: dark)')
      .addEventListener('change', () => {
        if (themeMode.value === 'auto') applyTheme()
      })
  }
  try {
    uni.onThemeChange?.((res: any) => {
      if (themeMode.value === 'auto') applyTheme()
    })
  } catch {}
}

export function useTheme() {
  return {
    themeMode,
    effectiveTheme,
    setTheme,
    toggleTheme,
    applyTheme,
    initTheme,
  }
}

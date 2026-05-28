import { createSSRApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import uviewPlus from 'uview-plus'
import lyCharts from '@/uni_modules/ly-charts'
import { i18n } from './locales'
import 'uno.css'

export function createApp() {
  const app = createSSRApp(App)
  app.use(createPinia())
  app.use(uviewPlus)
  app.use(lyCharts, () => ({}))
  app.use(i18n)
  return { app }
}

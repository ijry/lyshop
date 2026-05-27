import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import { i18n } from './locales'
import { setNotifyHandler } from './utils/notify'
import { pushToast } from './utils/toast'
import 'virtual:uno.css'
import './style.css'

// 可按需替换为第三方消息组件，例如：
// setNotifyHandler(({ message, level }) => ElMessage({ message, type: level === 'error' ? 'error' : 'success' }))
setNotifyHandler((payload) => {
  pushToast(payload)
})

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.use(i18n)
app.mount('#app')

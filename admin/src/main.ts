import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import App from './App.vue'
import router from './router'
import { i18n } from './locales'
import { setNotifyHandler } from './utils/notify'
import { pushToast } from './utils/toast'
import 'virtual:uno.css'
import './style.css'
// 组件样式放在基础 reset 之后，避免按钮背景等被重置覆盖。
import 'element-plus/dist/index.css'

// 可按需替换为第三方消息组件，例如：
// setNotifyHandler(({ message, level }) => ElMessage({ message, type: level === 'error' ? 'error' : 'success' }))
setNotifyHandler((payload) => {
  pushToast(payload)
})

const app = createApp(App)
// 注册模板组件、v-loading 和服务组件，确保表格作用域插槽由组件传入行数据。
app.use(ElementPlus, { locale: zhCn })
app.use(createPinia())
app.use(router)
app.use(i18n)
app.mount('#app')

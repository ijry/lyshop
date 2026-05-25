import { createI18n } from 'vue-i18n'
import zhCN from './zh-CN'
import en from './en'

export const i18n = createI18n({
  legacy: false,
  locale: localStorage.getItem('locale') || 'zh-CN',
  fallbackLocale: 'zh-CN',
  messages: { 'zh-CN': zhCN, en },
})

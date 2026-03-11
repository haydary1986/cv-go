import { createI18n } from 'vue-i18n'
import ar from './ar'
import en from './en'

const i18n = createI18n({
  legacy: false,
  locale: localStorage.getItem('locale') || 'ar',
  fallbackLocale: 'en',
  messages: { ar, en }
})

export default i18n

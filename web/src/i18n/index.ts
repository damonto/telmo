import { createI18n } from 'vue-i18n'

import en from './locales/en'
import zh from './locales/zh'

const supportedLocales = ['en', 'zh'] as const
export type AppLocale = (typeof supportedLocales)[number]

const pickLocale = (languages: readonly string[]): AppLocale => {
  for (const language of languages) {
    const normalized = language.toLowerCase()
    if (normalized.startsWith('en')) return 'en'
    if (normalized.startsWith('zh')) return 'zh'
  }
  return 'en'
}

const detectLocale = (): AppLocale => {
  if (typeof navigator === 'undefined') {
    return 'en'
  }

  return pickLocale(navigator.languages ?? [navigator.language])
}

const i18n = createI18n({
  legacy: false,
  globalInjection: true,
  locale: detectLocale(),
  fallbackLocale: 'en',
  messages: { en, zh },
})

export default i18n

import { createI18n } from 'vue-i18n';

import enUS from '../locales/en.json';

export type MessageSchema = typeof enUS

export const i18n = createI18n<[MessageSchema]>({
  legacy: false,
  locale: localStorage.getItem('locale') ? localStorage.getItem('locale')! : undefined,
  fallbackLocale: 'en-US',
  messages: {
    'en-US': enUS,
  },
});
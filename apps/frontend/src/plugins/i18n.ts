import { createI18n } from 'vue-i18n';

import en from '../locales/en.json';
import ru from '../locales/ru.json';

import { localeStore } from '@/stores/locale';

export type MessageSchema = typeof en

export const i18n = createI18n<[MessageSchema], 'en' | 'ru'>({
  legacy: false,
  locale: localeStore.get(),
  fallbackLocale: 'en',
  messages: {
    'en': en,
    'ru': ru,
  },
});

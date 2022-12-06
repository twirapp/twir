import { createI18n } from 'vue-i18n';

import gb from '../locales/gb.json';

import { localeStore } from '@/stores/locale';

export type MessageSchema = typeof gb

export const i18n = createI18n<[MessageSchema], 'gb'>({
  legacy: false,
  locale: localeStore.get(),
  fallbackLocale: 'gb',
  messages: {
    'gb': gb,
    // 'ru': ru,
  },
});
import messages from '@intlify/unplugin-vue-i18n/messages';
import { createI18n } from 'vue-i18n';

import en from './locales/en.json';

export const i18n = createI18n({
  messages,
	locale: 'en',
});

type Lang = typeof en

declare module 'vue-i18n' {
  export interface DefineLocaleMessage extends Lang {
  }
}

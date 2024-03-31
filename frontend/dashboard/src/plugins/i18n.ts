import messages from '@intlify/unplugin-vue-i18n/messages';
import { useLocalStorage } from '@vueuse/core';
import { createI18n } from 'vue-i18n';

import type en from '../locales/en.json';

const locale = useLocalStorage('twirLocale', 'en');


export const i18n = createI18n({
	messages,
	locale: locale.value,
	fallbackLocale: 'en',
});

type Lang = typeof en

declare module 'vue-i18n' {
	export interface DefineLocaleMessage extends Lang {
	}
}

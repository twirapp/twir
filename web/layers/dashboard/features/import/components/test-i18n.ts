import { createI18n } from 'vue-i18n'

import en from '../../../../../i18n/locales/en.json'

export function createImportsTestI18n() {
	return createI18n({
		legacy: false,
		globalInjection: true,
		locale: 'en',
		messages: { de: en, en, es: en, ja: en, pt: en, ru: en, sk: en, uk: en },
	})
}

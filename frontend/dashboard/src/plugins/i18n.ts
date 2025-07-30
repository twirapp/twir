import { useLocalStorage } from '@vueuse/core'
import { isRef, nextTick, watch } from 'vue'
import { createI18n } from 'vue-i18n'

import type en from '@/locales/en.json'
import type { Composer, I18n, I18nMode, Locale, VueI18n } from 'vue-i18n'

function isComposer(instance: VueI18n | Composer, mode: I18nMode): instance is Composer {
	return mode === 'composition' && isRef(instance.locale)
}

export function setLocale(i18n: I18n, locale: Locale): void {
	if (isComposer(i18n.global, i18n.mode)) {
		i18n.global.locale.value = locale
	} else {
		i18n.global.locale = locale
	}
}

export function getLocale(i18n: I18n): string {
	if (isComposer(i18n.global, i18n.mode)) {
		return i18n.global.locale.value
	} else {
		return i18n.global.locale
	}
}

const getResourceMessages = (r: any) => r.default || r

export async function loadLocaleMessages(i18n: I18n, locale: Locale) {
	// load locale messages
	const messages = await import(`../locales/${locale}.json`).then(getResourceMessages)

	// set locale and locale message
	i18n.global.setLocaleMessage(locale, messages)

	return nextTick()
}

const locale = useLocalStorage('twirLocale', 'en')
watch(locale, (newLocale) => {
	loadLocaleMessages(i18n, newLocale)
})

export const AVAILABLE_LOCALES = [
	{
		code: 'en',
		name: 'English',
	},
	{
		code: 'ru',
		name: 'Russian',
	},
	{
		code: 'uk',
		name: 'Українська',
	},
	{
		code: 'de',
		name: 'Deutsch',
	},
	{
		code: 'ja',
		name: '日本語',
	},
	{
		code: 'sk',
		name: 'Slovenčina',
	},
	{
		code: 'es',
		name: 'Español',
	},
	{
		code: 'pt',
		name: 'Português',
	},
]

function setupI18n(): I18n {
	const i18n = createI18n({
		locale: locale.value,
		availableLocales: AVAILABLE_LOCALES.map((locale) => locale.code),
		fallbackLocale: 'en',
	}) as I18n

	setLocale(i18n, locale.value)

	loadLocaleMessages(i18n, 'en')
	if (locale.value !== 'en') {
		loadLocaleMessages(i18n, locale.value)
	}

	return i18n
}

export const i18n = setupI18n()

type Lang = typeof en

declare module 'vue-i18n' {
	export interface DefineLocaleMessage extends Lang {}
}

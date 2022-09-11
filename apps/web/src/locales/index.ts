import type { WritableComputedRef } from 'vue';
import { createI18n } from 'vue-i18n';

import type ILandingLocale from '@/locales/landing/interface.js';

export type Locale = 'en' | 'ru';

export const locales: Locale[] = ['en', 'ru'];
export const defaultLocale: Locale = 'en';

type Languages = { name: string; locale: Locale }[];

export const languages: Languages = [
  { name: 'English', locale: 'en' },
  { name: 'Русский', locale: 'ru' },
];

export interface LocaleTypes {
  landing: ILandingLocale;
  app: ReturnType<() => typeof import('@/locales/app/en.json')>;
}

export type LocaleType = keyof LocaleTypes;

export async function loadLocaleMessages<L extends LocaleType>(
  localeType: L,
  locale: Locale,
): Promise<LocaleTypes[L]> {
  return (await import(`../locales/${localeType}/${locale}.ts`)).default;
}

export async function setupI18n(locale: Locale, localeType: LocaleType) {
  const i18n = createI18n({ locale, legacy: false });

  (i18n.global.locale as WritableComputedRef<Locale>).value = locale;
  const messages = await loadLocaleMessages(localeType, locale);
  i18n.global.setLocaleMessage<any>(locale, messages);

  return i18n;
}

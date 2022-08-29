import { nextTick, WritableComputedRef } from 'vue';
import { createI18n, I18n } from 'vue-i18n';

import type { Locale } from '@/types/locale';

export function setupI18n(
  options: { locale: Locale; legacy: boolean } = { locale: 'en', legacy: false },
) {
  const i18n = createI18n(options);
  setI18nLanguage(i18n, options.locale);
  return i18n;
}

export function setI18nLanguage(i18n: I18n, locale: Locale) {
  (i18n.global.locale as WritableComputedRef<Locale>).value = locale;
}

export async function loadLocaleMessages(
  i18n: I18n,
  localeType: 'landing' | 'app',
  locale: Locale,
) {
  const messages = await import(`../locales/${localeType}/${locale}.json`);

  i18n.global.setLocaleMessage(locale, messages.default);
  (i18n.global.locale as WritableComputedRef<Locale>).value = locale;

  return nextTick();
}

import type { WritableComputedRef } from 'vue';
import { type Locale, createI18n } from 'vue-i18n';

import type { LocaleType, LocaleTypes } from '@/locales';

export async function loadLocaleMessages<L extends LocaleType>(
  localeType: L,
  locale: Locale,
): Promise<LocaleTypes[L]> {
  return (await import(`../../locales/${localeType}/${locale}.ts`)).default;
}

export async function setupI18n(locale: Locale, localeType: LocaleType) {
  const i18n = createI18n({ locale, legacy: false });

  (i18n.global.locale as WritableComputedRef<Locale>).value = locale;
  const messages = await loadLocaleMessages(localeType, locale);
  i18n.global.setLocaleMessage<any>(locale, messages);

  return i18n;
}

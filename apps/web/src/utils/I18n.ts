import type { WritableComputedRef } from 'vue';
import { createI18n } from 'vue-i18n';

import type { Locale } from '@/types/locale';
import { loadLocaleMessages, LocaleType } from '@/utils/locales.js';

export async function setupI18n(locale: Locale = 'en', localeType: LocaleType) {
  const i18n = createI18n({ locale, legacy: false });

  (i18n.global.locale as WritableComputedRef<Locale>).value = locale;
  const messages = await loadLocaleMessages(localeType, locale);
  i18n.global.setLocaleMessage<any>(locale, messages);

  return i18n;
}

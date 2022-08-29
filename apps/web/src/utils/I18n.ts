import type { WritableComputedRef } from 'vue';
import { createI18n } from 'vue-i18n';

import type { Locale } from '@/types/locale';

export function setupI18n(locale: Locale = 'en') {
  const i18n = createI18n({ locale, legacy: false });
  (i18n.global.locale as WritableComputedRef<Locale>).value = locale;
  return i18n;
}

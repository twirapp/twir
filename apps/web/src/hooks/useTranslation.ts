import { useI18n, type ComposerTranslation, type Locale } from 'vue-i18n';

import type { LocaleType, LocaleTypes } from '@/locales';

export default function<L extends LocaleType>() {
  const { t } = useI18n();

  return t as ComposerTranslation<{ en: LocaleTypes[L] }, Locale>;
}

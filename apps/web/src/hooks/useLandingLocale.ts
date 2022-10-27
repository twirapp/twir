import { navigate } from 'vite-plugin-ssr/client/router';
import { useI18n } from 'vue-i18n';

import { loadLocaleMessages, type Locale } from '@/locales';

/**
 * @returns function to set landing locale
 */
export default function (): (locale: Locale) => Promise<void> {
  const { setLocaleMessage, locale: i18nLocale } = useI18n();

  return async (locale: Locale) => {
    const messages = await loadLocaleMessages('landing', locale);

    setLocaleMessage<any>(locale, messages);
    i18nLocale.value = locale;

    navigate(`/${locale}`, { keepScrollPosition: true });
  };
}

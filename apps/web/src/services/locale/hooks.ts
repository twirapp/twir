import { navigate } from 'vite-plugin-ssr/client/router';
import { LocaleMessages, useI18n } from 'vue-i18n';

import { loadLocaleMessages } from './i18n.js';

import type { Locale, LocaleTypes } from '@/locales';

/**
 * @returns function to set landing locale
 */
export function useLandingLocale() {
  const { setLocaleMessage, locale: i18nLocale } = useI18n();

  return {
    setLandingLocale: async (locale: Locale) => {
      const messages = await loadLocaleMessages('landing', locale);

      setLocaleMessage<any>(locale, messages);
      i18nLocale.value = locale;

      navigate(`/${locale}`, { keepScrollPosition: true });
      document.documentElement.setAttribute('lang', locale);
    },
  };
}

export function useTranslation<L extends keyof LocaleTypes>() {
  return useI18n<LocaleMessages<{ en: LocaleTypes[L] }>>();
}

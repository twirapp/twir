import { getCookie, setCookie } from 'cookies-next';
import { RU, US } from 'country-flag-icons/react/3x2';
import { useRouter } from 'next/router';
import { useCallback, useEffect, useState } from 'react';

import i18nconfig from '../../next-i18next.config';

// put in constants.ts
const DEFAULT_LOCALE = 'en';
const LOCALE_COOKIE_KEY = 'locale';
const SUPPORTED_LOCALES = i18nconfig.i18n.locales;
const ICON_SIZE = 14;

export const LOCALES = new Map([
  ['en', { name: 'English', icon: <US style={{ height: ICON_SIZE }} /> }],
  ['ru', { name: 'Russian', icon: <RU style={{ height: ICON_SIZE }} /> }],
]);

for (const locale of SUPPORTED_LOCALES) {
  if (!LOCALES.has(locale)) {
    throw new Error(`Locale "${locale}" is not implemented.`);
  }
}

export const useLocale = () => {
  const router = useRouter();
  const [locale, setLocale] = useState(() => getCookie(LOCALE_COOKIE_KEY) as string);

  const toggleLocale = useCallback(
    (newLocale = DEFAULT_LOCALE) => {
      if (!LOCALES.has(newLocale)) return;
      setLocale(newLocale);
      setCookie(LOCALE_COOKIE_KEY, newLocale);
    },
    [locale],
  );

  const isSupportedLocale = () => {
    return LOCALES.has(locale);
  };

  useEffect(() => {
    if (isSupportedLocale()) {
      const { pathname, asPath, query } = router;
      if (query.code || query.token) return;
      router.push({ pathname, query }, asPath, { locale });
    } else {
      toggleLocale();
    }
  }, [locale]);

  return {
    locale,
    toggleLocale,
    isSupportedLocale,
  };
};

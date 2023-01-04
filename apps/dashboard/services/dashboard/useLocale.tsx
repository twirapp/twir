import { useLocalStorage } from '@mantine/hooks';
import { RU, US } from 'country-flag-icons/react/3x2';
import { useCallback } from 'react';

import i18nconfig from '../../next-i18next.config';

const DEFAULT_LOCALE = 'en';
const LOCALE_STORAGE_KEY = 'locale';
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
  const [locale, setLocale] = useLocalStorage({
    key: LOCALE_STORAGE_KEY,
    defaultValue: DEFAULT_LOCALE,
  });

  const toggleLocale = useCallback(
    (newLocale = DEFAULT_LOCALE) => {
      if (!LOCALES.has(newLocale)) return;
      setLocale(newLocale);
    },
    [locale],
  );

  const isSupportedLocale = useCallback(() => {
    return LOCALES.has(locale);
  }, [locale]);

  return {
    locale,
    toggleLocale,
    isSupportedLocale,
  };
};

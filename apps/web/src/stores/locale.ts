import { persistentAtom } from '@nanostores/persistent';
import { usePreferredLanguages } from '@vueuse/core';

import { defaultLocale, Locale, locales } from '@/locales';

const getUserPreferedLocale = (): Locale => {
  let { value: userLangs } = usePreferredLanguages();

  // normaliza lang format from 'ru-RU' to just 'ru'
  userLangs = userLangs.map((lang) => (lang.length > 2 ? lang.slice(0, 2) : lang));

  let preferedLocale: Locale = defaultLocale;

  for (let i = 0; i < userLangs.length; i++) {
    if (locales.includes(userLangs[i] as Locale)) {
      preferedLocale = userLangs[i] as Locale;
      break;
    }
  }

  return preferedLocale;
};

export const localeStore = persistentAtom<Locale>('locale', getUserPreferedLocale());

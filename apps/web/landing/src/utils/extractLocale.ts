import type { Locale } from '@/types/locale';
import { defaultLocale, locales } from '@/utils/locales';

export function extractLocale(url: string) {
  const urlPaths = url.split('/');

  let locale;
  let urlWithoutLocale;

  const firstPath = urlPaths[1];

  if (locales.filter((locale) => locale !== defaultLocale).includes(firstPath as Locale)) {
    locale = firstPath;
    urlWithoutLocale = '/' + urlPaths.slice(2).join('/');
  } else {
    locale = defaultLocale;
    urlWithoutLocale = url;
  }

  console.log(locale, urlWithoutLocale);

  return { locale, urlWithoutLocale };
}

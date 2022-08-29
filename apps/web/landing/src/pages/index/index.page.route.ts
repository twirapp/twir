import { resolveRoute } from 'vite-plugin-ssr/routing';

import type { Locale } from '@/types/locale.js';
import type { PageContext } from '@/types/pageContext.js';
import { defaultLocale, locales } from '@/utils/locales';

export default (pageContext: PageContext) => {
  if (pageContext.urlPathname === '/')
    return { routeParams: { locale: defaultLocale }, match: true };

  const result = resolveRoute('/@locale', pageContext.urlPathname);
  if (!result.match) return false;

  const locale = result.routeParams.locale;

  if (!locales.includes(locale as Locale)) {
    return false;
  }

  return result;
};

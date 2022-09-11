import { resolveRoute } from 'vite-plugin-ssr/routing';

import { defaultLocale, locales, type Locale } from '@/locales';
import type { PageContext } from '@/types/pageContext.js';

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

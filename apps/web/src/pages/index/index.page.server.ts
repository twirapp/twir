import { renderToNodeStream } from '@vue/server-renderer';

import { createApp } from './app.js';

import { seoLocales } from '@/data/seo.js';
import { defaultLocale, locales, type Locale } from '@/locales';
import { setupI18n } from '@/services/locale';
import { htmlLayout } from '@/utils/htmlLayout.js';
import type { PageContext } from '@/utils/pageContext.js';
import type { PassToClient, PrerenderFn } from '@/utils/vitePluginSSR.js';

export async function render(pageContext: PageContext) {
  const locale = (pageContext.routeParams as { locale: Locale }).locale;
  pageContext.locale = locale;

  const app = createApp(pageContext);

  const i18n = await setupI18n(locale, 'landing');
  app.use(i18n);

  const seoInfo = seoLocales[locale];

  const documentHtml = htmlLayout(seoInfo, pageContext, renderToNodeStream(app));

  return {
    documentHtml,
    pageContext: {
      enableEagerStreaming: true,
      locale,
    },
  };
}

export const passToClient: PassToClient = ['pageProps', 'locale'];

export const prerender: PrerenderFn = () => [
  {
    pageContext: {
      locale: defaultLocale,
    },
    url: '/',
  },
  ...locales.map((locale) => ({
    url: `/${locale}`,
    pageContext: {
      locale,
    },
  })),
];

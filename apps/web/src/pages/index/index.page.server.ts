import { renderToNodeStream } from '@vue/server-renderer';

import { seoLocales } from '@/data/seo.js';
import { defaultLocale, locales, setupI18n, type Locale } from '@/locales';
import { createApp } from '@/pages/index/app';
import type { PageContext } from '@/types/pageContext.js';
import type { PassToClient, PrerenderFn } from '@/types/vitePluginSSR.js';
import { htmlLayout } from '@/utils/htmlLayout.js';

export async function render(pageContext: PageContext) {
  const locale = (pageContext.routeParams as { locale: Locale }).locale;
  pageContext.locale = locale;

  const app = createApp(pageContext);

  const i18n = await setupI18n(locale, 'landing');
  app.use(i18n);

  const seoInfo = seoLocales[locale];

  const documentHtml = htmlLayout({
    title: seoInfo.title,
    description: seoInfo.description,
    keywords: seoInfo.keywords,
    locale,
    urlCanonical: pageContext.urlParsed.origin || undefined,
    urlOriginal: pageContext.urlOriginal,
    content: renderToNodeStream(app),
  });

  return {
    documentHtml,
    pageContext: {
      enableEagerStreaming: true,
      locale,
    },
  };
}

export const passToClient: PassToClient = ['pageProps', 'documentProps', 'locale'];

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

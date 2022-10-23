import { seoLocales } from '@/data/seo.js';
import { defaultLocale, Locale } from '@/locales/index.js';
import type { PageContext } from '@/types/pageContext.js';
import type { PassToClient } from '@/types/vitePluginSSR.js';
import { htmlLayout } from '@/utils/htmlLayout.js';

export async function render(pageContext: PageContext) {
  const locale: Locale = defaultLocale;

  const seoInfo = seoLocales[locale];

  const documentHtml = htmlLayout({
    title: seoInfo.title,
    description: seoInfo.description,
    keywords: seoInfo.keywords,
    locale,
    urlCanonical: pageContext.urlParsed.origin || undefined,
    urlOriginal: pageContext.urlOriginal,
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

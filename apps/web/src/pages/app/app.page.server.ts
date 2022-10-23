import { seoLocales } from '@/data/seo.js';
import type { Locale } from '@/locales/index.js';
import { defaultLocale } from '@/locales/index.js';
import type { PageContext } from '@/types/pageContext.js';
import { htmlLayout } from '@/utils/htmlLayout.js';

export async function render(pageContext: PageContext) {
  const locale: Locale = defaultLocale;

  const seoInfo = seoLocales[locale];

  const documentHtml = htmlLayout({
    title: 'App',
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

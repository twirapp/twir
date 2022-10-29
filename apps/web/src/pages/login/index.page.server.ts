import { seoLocales } from '@/data/seo.js';
import { defaultLocale, Locale } from '@/locales/index.js';
import { htmlLayout } from '@/utils/htmlLayout.js';
import type { PageContext } from '@/utils/pageContext.js';
import type { PassToClient } from '@/utils/vitePluginSSR.js';

export async function render(pageContext: PageContext) {
  const locale: Locale = defaultLocale;

  const seoInfo = seoLocales[locale];
  const documentHtml = htmlLayout(seoInfo, pageContext);

  return {
    documentHtml,
    pageContext: {
      enableEagerStreaming: true,
      locale,
    },
  };
}

export const passToClient: PassToClient = ['pageProps', 'documentProps', 'locale'];

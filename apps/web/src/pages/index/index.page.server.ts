import { renderToNodeStream } from '@vue/server-renderer';
import { escapeInject } from 'vite-plugin-ssr';

import { createApp } from '@/pages/index/app';
import type { Locale } from '@/types/locale.js';
import type { PageContext } from '@/types/pageContext.js';
import { getPageTitle } from '@/utils/getPageTitle.js';
import { setupI18n } from '@/utils/I18n.js';
import { defaultLocale, loadLocaleMessages, locales } from '@/utils/locales.js';

export const passToClient = ['pageProps', 'documentProps', 'locale'];

export async function render(pageContext: PageContext) {
  const locale = (pageContext.routeParams as { locale: Locale }).locale;
  pageContext.locale = locale;

  const app = createApp(pageContext);

  const i18n = setupI18n(locale);
  const messages = await loadLocaleMessages('landing', locale);
  i18n.global.setLocaleMessage<any>(locale, messages);
  app.use(i18n);

  const stream = renderToNodeStream(app);

  const title = getPageTitle(pageContext);

  const documentHtml = escapeInject`<!DOCTYPE html>
    <html>
      <head>
        <meta charset="utf-8" />
        <title>${title}</title>
      </head>
      <body>
        <div id="app">${stream}</div>
      </body>
    </html>`;

  return {
    documentHtml,
    pageContext: {
      enableEagerStreaming: true,
      locale,
    },
  };
}

export function prerender() {
  return [
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
}

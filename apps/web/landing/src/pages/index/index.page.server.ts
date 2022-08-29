import { renderToNodeStream } from '@vue/server-renderer';
import { escapeInject } from 'vite-plugin-ssr';

import type { Locale } from '../../types/locale.js';
// import { defaultLocale, locales } from '../../utils/locales.js';

import { createApp } from '@/pages/index/app';
import type { PageContext } from '@/types/pageContext.js';
import { loadLocaleMessages, setupI18n } from '@/utils/createI18n.js';

export const passToClient = ['pageProps', 'documentProps', 'locale'];

export async function render(pageContext: PageContext) {
  const app = createApp(pageContext);

  const locale = (pageContext.routeParams as { locale: Locale }).locale;

  const i18n = setupI18n();
  loadLocaleMessages(i18n, 'landing', locale);
  app.use(i18n);

  console.log(i18n.global.t('hello'));

  const stream = renderToNodeStream(app);

  const documentHtml = escapeInject`<!DOCTYPE html>
    <html>
      <head>
        <meta charset="utf-8">
        <title>Tsuwari</title>
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

// export function onBeforePrerender(globalContext: { prerenderPageContexts: PageContext[] }) {
//   const prerenderPageContexts: PageContext[] = [];

//   globalContext.prerenderPageContexts.forEach((pageContext) => {
//     prerenderPageContexts.push({
//       ...pageContext,
//       locale: defaultLocale,
//     });
//     locales
//       .filter((locale) => locale !== defaultLocale)
//       .forEach((locale) => {
//         prerenderPageContexts.push({
//           ...pageContext,
//           urlOriginal: `/${locale}${pageContext.urlOriginal}`,
//           locale,
//         });
//       });
//   });

//   return {
//     globalContext: {
//       prerenderPageContexts,
//     },
//   };
// }

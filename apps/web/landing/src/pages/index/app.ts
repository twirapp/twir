import { createSSRApp, reactive } from 'vue';

import { setPageContext } from '../../hooks/usePageContext.js';

import type { PageContext } from '@/types/pageContext.js';

export { createApp };

function createApp(pageContext: PageContext) {
  const { Page, pageProps } = pageContext;

  const app = createSSRApp(Page, pageProps);

  const pageContextReactive = reactive(pageContext);
  setPageContext(app, pageContextReactive);

  return app;
}

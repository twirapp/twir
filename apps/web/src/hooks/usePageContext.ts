import { App, inject } from 'vue';

import type { PageContext } from '../types/pageContext.js';

const PAGE_CONTEXT_INJECTION_KEY = Symbol();

export function usePageContext(): PageContext {
  return inject(PAGE_CONTEXT_INJECTION_KEY) as PageContext;
}

export function setPageContext(app: App, pageContext: PageContext) {
  app.provide(PAGE_CONTEXT_INJECTION_KEY, pageContext);
}

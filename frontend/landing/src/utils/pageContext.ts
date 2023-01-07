import type { PageContextBuiltIn } from 'vite-plugin-ssr';
import { App, Component, inject, InjectionKey } from 'vue';

import type { Locale } from '@/locales';

export interface PageContext extends PageContextBuiltIn {
  Page: Component;
  pageProps: Record<string, unknown>;
  exports: {
    documentProps?: {
      title: string;
    };
  };
  locale: Locale;
}

const PAGE_CONTEXT_INJECTION_KEY = Symbol('PageContext') as InjectionKey<PageContext>;

export function usePageContext() {
  const pageContext = inject(PAGE_CONTEXT_INJECTION_KEY);
  if (!pageContext) {
    throw new Error('Cannot inject PageContext');
  }

  return pageContext;
}

export function setPageContext(app: App, pageContext: PageContext) {
  app.provide(PAGE_CONTEXT_INJECTION_KEY, pageContext);
}

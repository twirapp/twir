import type { PageContextBuiltIn } from 'vite-plugin-ssr';
import type { Component } from 'vue';

import type { Locale } from '@/locales';

export interface PageContext extends PageContextBuiltIn {
  Page: Component;
  pageProps: Record<string, unknown>;
  exports: {
    documentProps?: {
      title: string;
    };
  };
  documentProps?: {
    title: string;
  };
  locale: Locale;
}

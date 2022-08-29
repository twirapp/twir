import type { PageContextBuiltIn } from 'vite-plugin-ssr';

import type { Locale } from '@/types/locale';

export interface PageContext extends PageContextBuiltIn {
  pageProps: Record<string, unknown>;
  documentProps?: {
    title: string;
  };
  locale: Locale;
}

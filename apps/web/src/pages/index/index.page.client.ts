import { seoLocales } from '@/data/seo.js';
import type { Locale } from '@/locales/index.js';
import { setupI18n } from '@/locales/index.js';
import { createApp } from '@/pages/index/app';
import '@/styles/tailwind.base.css';
import type { PageContext } from '@/types/pageContext';

export const clientRouting = true;
export const prefetchStaticAssets = { when: 'VIEWPORT' };

let app: ReturnType<typeof createApp>;
export async function render(pageContext: PageContext) {
  if (!app) {
    app = createApp(pageContext);

    const i18n = await setupI18n(pageContext.locale, 'landing');
    app.use(i18n);

    app.mount('#app');
    return;
  }

  const locale = (pageContext.routeParams as { locale: Locale }).locale;
  pageContext.locale = locale;

  document.title = seoLocales[locale].title;
  app.changePage(pageContext);
}

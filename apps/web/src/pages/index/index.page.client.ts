import { landingPage } from '@/data/seo.js';
import { setupI18n, type Locale } from '@/locales/index.js';
import { createApp } from '@/pages/index/app';
import '@/styles/tailwind.base.css';
import type { PageContext } from '@/types/pageContext';

export const clientRouting = true;
export const prefetchStaticAssets = { when: 'VIEWPORT' };
export { render };

let app: ReturnType<typeof createApp>;
async function render(pageContext: PageContext) {
  if (!app) {
    app = createApp(pageContext);

    const i18n = await setupI18n(pageContext.locale, 'landing');
    app.use(i18n);

    app.mount('#app');
    return;
  }

  const locale = (pageContext.routeParams as { locale: Locale }).locale;
  pageContext.locale = locale;

  document.title = landingPage[locale].title;
  app.changePage(pageContext);
}

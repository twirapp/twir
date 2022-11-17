import { seoLocales } from '@/data/seo.js';
import type { Locale } from '@/locales/index.js';
import { createApp } from '@/pages/index/app';
import { setupI18n } from '@/services/locale/i18n.js';
import '@/styles/tailwind.base.css';
import type { PageContext } from '@/utils/pageContext.js';

export const clientRouting = true;
export const prefetchStaticAssets = { when: 'VIEWPORT' };

let app: ReturnType<typeof createApp>;
export async function render(pageContext: PageContext) {
  if (!app) {
    app = createApp(pageContext);

    const i18n = await setupI18n(pageContext.locale, 'landing');
    app.use(i18n).mount('#app');

    return;
  }

  const locale = (pageContext.routeParams as { locale: Locale }).locale;
  pageContext.locale = locale;

  document.title = seoLocales[locale].title;
  app.changePage(pageContext);
}

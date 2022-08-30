import { createApp } from '@/pages/index/app';
import '@/styles/tailwind.base.css';
import type { Locale } from '@/types/locale.js';
import type { PageContext } from '@/types/pageContext';
import { setupI18n } from '@/utils/locales.js';

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

  app.changePage(pageContext);
}

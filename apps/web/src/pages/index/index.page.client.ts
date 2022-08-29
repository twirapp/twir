import { createApp } from '@/pages/index/app';
import '@/styles/tailwind.base.css';
import type { PageContext } from '@/types/pageContext';
import { getPageTitle } from '@/utils/getPageTitle.js';
import { setupI18n } from '@/utils/I18n.js';

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
  }

  app.changePage(pageContext);

  document.title = getPageTitle(pageContext);
}

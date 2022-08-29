import { createApp } from '@/pages/index/app';
import '@/styles/tailwind.base.css';
import type { PageContext } from '@/types/pageContext';
import { loadLocaleMessages, setupI18n } from '@/utils/createI18n.js';

export const clientRouting = true;
export const prefetchStaticAssets = { when: 'VIEWPORT' };
export { render };

let app: ReturnType<typeof createApp>;
async function render(pageContext: PageContext) {
  app = createApp(pageContext);

  const i18n = setupI18n();
  loadLocaleMessages(i18n, 'landing', pageContext.locale);
  app.use(i18n);

  app.mount('#app');

  i18n.global.setLocaleMessage;
}

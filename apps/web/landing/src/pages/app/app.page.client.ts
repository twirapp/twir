import { createApp } from 'vue';

import { createAppRouter } from '@/pages/app/router';
import '@/styles/tailwind.base.css';
import type { PageContext } from '@/types/pageContext';

// export const clientRouting = true;
// export const prefetchStaticAssets = { when: 'VIEWPORT' };

export async function render(pageContext: PageContext) {
  const app = createApp(pageContext.Page);
  // const i18n = setupI18n();
  // loadLocaleMessages(i18n, pageContext.locale);
  // app.use(i18n);

  console.log('client');

  const router = createAppRouter();
  app.use(router);
  await router.isReady();

  app.mount('#app');
}

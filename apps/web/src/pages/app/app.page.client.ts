import { createApp } from 'vue';

import { createAppRouter } from '@/pages/app/router';
import '@/styles/tailwind.base.css';
import type { PageContext } from '@/types/pageContext';
import { loadLocaleMessages, setupI18n } from '@/utils/I18n.js';

export async function render(pageContext: PageContext) {
  const app = createApp(pageContext.Page);

  // TODO save locale to localstorage and manage it from landing and app

  const i18n = setupI18n();
  loadLocaleMessages(i18n, 'app', 'en');
  app.use(i18n);

  console.log('client');

  const router = createAppRouter();
  app.use(router);
  await router.isReady();

  app.mount('#app');
}

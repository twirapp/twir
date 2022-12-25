import { VueQueryPlugin } from '@tanstack/vue-query';
import { createApp } from 'vue';

import { createAppRouter } from '@/pages/app/router';
import '@/styles/tailwind.base.css';
import { setupI18n } from '@/services/locale';
import type { PageContext } from '@/utils/pageContext.js';

export async function render(pageContext: PageContext) {
  const app = createApp(pageContext.Page);

  // TODO save locale to localstorage and manage it from landing and app

  const i18n = await setupI18n('en', 'app');
  app.use(i18n);

  app.use(VueQueryPlugin);

  const router = createAppRouter();
  app.use(router);
  await router.isReady();

  app.mount('#app');
}

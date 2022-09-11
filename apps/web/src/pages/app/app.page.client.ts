import { createApp } from 'vue';

import { setupI18n } from '@/locales';
import { createAppRouter } from '@/pages/app/router';
import '@/styles/tailwind.base.css';
import type { PageContext } from '@/types/pageContext';

export async function render(pageContext: PageContext) {
  const app = createApp(pageContext.Page);

  // TODO save locale to localstorage and manage it from landing and app

  const i18n = await setupI18n('en', 'app');
  app.use(i18n);

  const router = createAppRouter();
  app.use(router);
  await router.isReady();

  app.mount('#app');
}

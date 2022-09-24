import { createApp } from 'vue';

import { setupI18n } from '@/locales';
import '@/styles/tailwind.base.css';
import type { PageContext } from '@/types/pageContext';

export async function render(pageContext: PageContext) {
  const app = createApp(pageContext.Page);

  const i18n = await setupI18n('en', 'app');
  app.use(i18n);

  app.mount('#app');
}

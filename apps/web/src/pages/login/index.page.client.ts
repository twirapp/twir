import { VueQueryPlugin } from '@tanstack/vue-query';
import { createApp } from 'vue';

import { defaultLocale, setupI18n } from '@/locales';
import '@/styles/tailwind.base.css';
import type { PageContext } from '@/types/pageContext';

export async function render(pageContext: PageContext) {
  const app = createApp(pageContext.Page);

  const i18n = await setupI18n(defaultLocale, 'app');

  app.use(i18n).use(VueQueryPlugin).mount('#app');
}

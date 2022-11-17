import { VueQueryPlugin } from '@tanstack/vue-query';
import { createApp } from 'vue';

import { defaultLocale } from '@/locales';
import { setupI18n } from '@/services/locale';
import '@/styles/tailwind.base.css';
import type { PageContext } from '@/utils/pageContext.js';

export async function render(pageContext: PageContext) {
  const app = createApp(pageContext.Page);

  const i18n = await setupI18n(defaultLocale, 'app');

  app.use(i18n).use(VueQueryPlugin).mount('#app');
}

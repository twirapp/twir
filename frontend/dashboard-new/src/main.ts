import { VueQueryPlugin } from '@tanstack/vue-query';
import { createApp } from 'vue';

import { getProfile } from './api/index.js';
import App from './App.vue';
import { router } from './router.js';

await getProfile();
createApp(App)
	.use(router)
	.use(VueQueryPlugin)
	.mount('#app');

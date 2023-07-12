import { VueQueryPlugin } from '@tanstack/vue-query';
import { createApp } from 'vue';

import { router } from './router.js';

import App from '@/App.vue';

createApp(App)
	.use(router)
	.use(VueQueryPlugin)
	.mount('#app');

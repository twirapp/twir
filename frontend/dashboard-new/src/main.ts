import { VueQueryPlugin } from '@tanstack/vue-query';
import { createApp } from 'vue';

import { router } from './router.js';

createApp({
	template: '<router-view />',
})
	.use(router)
	.use(VueQueryPlugin)
	.mount('#app');

import { createPinia } from 'pinia';
import { createApp } from 'vue';

import MainApp from './app.vue';
import { routes } from './routes.js';
import './style.css';

const pinia = createPinia();
createApp(MainApp).use(pinia).use(routes).mount('#app');

// refresh the page when new version comes
document.body.addEventListener('plugin_web_update_notice', () => {
	window.location.reload();
});

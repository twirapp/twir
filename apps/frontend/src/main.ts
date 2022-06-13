import { createApp } from 'vue';
import 'tw-elements';

import App from './App.vue';
import { i18n } from './plugins/i18n';
import { router } from './plugins/router';
import './main.css';
import './plugins/socket';

const app = createApp(App).use(i18n).use(router);

app.mount('#app');


import { createApp } from 'vue';
import 'tw-elements';

import App from './App.vue';
import { router } from './plugins/router';
import './main.css';
import './plugins/socket';

const app = createApp(App).use(router);

app.mount('#app');


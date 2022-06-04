import { createApp, onMounted } from 'vue';
import 'tw-elements';

import App from './App.vue';
import { fetchAndSetUser } from './functions/fetchAndSetUser';
import { router } from './plugins/router';
import './main.css';
import './plugins/socket';

const app = createApp(App).use(router);

app.mount('#app');

if (router.currentRoute.value.path === '/') {
  const accessToken = localStorage.getItem('accessToken');
  const refreshToken = localStorage.getItem('refreshToken');
  if (accessToken && refreshToken) {
    fetchAndSetUser();
  }
}

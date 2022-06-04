import { createApp } from 'vue';

import App from './App.vue';
import { fetchAndSetUser } from './functions/fetchAndSetUser';
import { router } from './plugins/router';
import './main.css';
import './plugins/socket';

const accessToken = localStorage.getItem('accessToken');
const refreshToken = localStorage.getItem('refreshToken');
if (accessToken && refreshToken) {
  fetchAndSetUser();
}

createApp(App).use(router).mount('#app');

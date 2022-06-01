import { Quasar } from 'quasar';
import { createApp } from 'vue';

import '@quasar/extras/material-icons/material-icons.css';
import '@quasar/extras/material-icons-outlined/material-icons-outlined.css';
import '@quasar/extras/material-icons-round/material-icons-round.css';
import '@quasar/extras/material-icons-sharp/material-icons-sharp.css';
import '@quasar/extras/material-symbols-outlined/material-symbols-outlined.css';
import '@quasar/extras/material-symbols-rounded/material-symbols-rounded.css';
import '@quasar/extras/material-symbols-sharp/material-symbols-sharp.css';
import 'quasar/src/css/index.sass';
import 'quasar/src/css/index.sass';

import App from './App.vue';
import { fetchAndSetUser } from './functions/fetchAndSetUser';
import { router } from './plugins/router';
import './plugins/socket';

const accessToken = localStorage.getItem('accessToken');
const refreshToken = localStorage.getItem('refreshToken');
if (accessToken && refreshToken) {
  fetchAndSetUser();
}

createApp(App)
  .use(Quasar, {
    plugins: {}, // import Quasar plugins and add here
    config: {
      dark: 'auto', // or Boolean true/false
    },
  })
  .use(router)
  .mount('#app');

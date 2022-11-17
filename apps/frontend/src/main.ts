import MasonryWall from '@yeger/vue-masonry-wall';
import { createApp } from 'vue';
import Toast, { PluginOptions, POSITION, TYPE } from 'vue-toastification';
import { inject } from '@vercel/analytics';

import 'tw-elements';

import App from './App.vue';
import { i18n } from './plugins/i18n';
import { router } from './plugins/router';

import './main.css';
import 'vue-toastification/dist/index.css';
import 'flag-icons/css/flag-icons.css';

const app = createApp(App).use(MasonryWall).use(i18n).use(router);

app.use(Toast, {
  position: POSITION.TOP_RIGHT,
  pauseOnFocusLoss: false,
  closeOnClick: true,
  maxToasts: 5,
  toastDefaults: {
    [TYPE.ERROR]: {
      timeout: 5000,
      closeButton: false,
    },
    [TYPE.SUCCESS]: {
      timeout: 3000,
      hideProgressBar: true,
    },
  },
} as PluginOptions);

app.mount('#app');

/* async function checkIfUpdateAvailable() {
  if (!import.meta.env.PROD) return;

  const request = await fetch('/api/version');
  if (!request.ok) return;

  const data = await request.text();

  const sha = import.meta.env.VITE_VERCEL_GIT_COMMIT_SHA;
  if (sha != data) window.location.reload();
}

setInterval(checkIfUpdateAvailable, 20 * 1000);
checkIfUpdateAvailable(); */
inject();

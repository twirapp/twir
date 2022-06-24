import { createApp } from 'vue';
import 'tw-elements';
import Toast, { PluginOptions, POSITION, TYPE } from 'vue-toastification';

import App from './App.vue';
import { i18n } from './plugins/i18n';
import { router } from './plugins/router';
import './main.css';
import './plugins/socket';
import 'vue-toastification/dist/index.css';

const app = createApp(App).use(i18n).use(router);

app.use(Toast, {
  position: POSITION.TOP_RIGHT,
  pauseOnFocusLoss: false,
 'closeOnClick': true,
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


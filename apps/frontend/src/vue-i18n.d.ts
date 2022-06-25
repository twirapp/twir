import 'vue-i18n';

import enUS from './locales/gb.json';

type Messages = typeof enUS

declare module 'vue-i18n' {
  // eslint-disable-next-line @typescript-eslint/no-empty-interface
  export interface DefineLocaleMessage extends Messages { }
}
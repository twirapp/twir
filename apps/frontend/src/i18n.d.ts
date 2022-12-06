import {
  DefineLocaleMessage,
} from 'vue-i18n';

import gb from './locales/gb.json';

type Messages = typeof gb

declare module 'vue-i18n' {
  // define the locale messages schema
  // eslint-disable-next-line @typescript-eslint/no-empty-interface
  export interface DefineLocaleMessage extends Messages {}

}
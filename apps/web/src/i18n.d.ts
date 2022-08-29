import appEn from '@/locales/app/en.json';
import landingEn from '@/locales/landing/en.json';

import 'vue-i18n';

declare module 'vue-i18n' {
  type Messages = typeof appEn & typeof landingEn;

  // eslint-disable-next-line @typescript-eslint/no-empty-interface
  export interface DefineLocaleMessage extends Messages {}
}

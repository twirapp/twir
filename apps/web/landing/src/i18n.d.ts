import 'vue-i18n';

declare module 'vue-i18n' {
  import enUS from '@/locales/landing/en.json';

  type Messages = typeof enUS;

  // eslint-disable-next-line @typescript-eslint/no-empty-interface
  export interface DefineLocaleMessage extends Messages {}
}

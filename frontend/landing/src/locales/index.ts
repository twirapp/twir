import type IAppLocale from '@/locales/app/interface.js';
import type ILandingLocale from '@/locales/landing/interface.js';

export type Locale = 'en' | 'ru';

export const locales: Locale[] = ['en', 'ru'];
export const defaultLocale: Locale = 'en';

type Languages = { name: string; locale: Locale }[];

export const languages: Languages = [
  { name: 'English', locale: 'en' },
  { name: 'Русский', locale: 'ru' },
];

export interface LocaleTypes {
  landing: ILandingLocale;
  app: IAppLocale;
}

export type LocaleType = keyof LocaleTypes;

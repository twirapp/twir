import type { Locale } from '@/locales';

export interface SeoPageProps {
  title: string;
  description: string;
  keywords: string[];
}

export const author = 'Satont, me@satont.dev';
export const ogImage = ''; // TODO

export const seoLocales: { [K in Locale]: SeoPageProps } = {
  en: {
    title: 'Twir - Main page',
    description:
      'Powerful and useful Twitch bot that helps manage chat on big channels. Developed from streamers for streamers with love.',
    keywords: ['twir', 'twir bot', 'tsuwari', 'twitch bot', 'tsuwari bot'],
  },
  ru: {
    title: 'Twir - Главная страница',
    description:
      'Мощный и полезный бот для Twitch, который помогает управлять чатом на больших каналах. Разработан от стримеров для стримеров c любовью.',
    keywords: ['twir', 'twir bot', 'tsuwari', 'twitch bot', 'tsuwari bot'],
  },
};

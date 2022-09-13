import type { Locale } from '@/locales';

export interface SeoPageProps {
  title: string;
  description: string;
  keywords: string[];
}

export const author = 'Satont, me@satont.dev';
export const ogImage = ''; // TODO

export const landingPage: { [K in Locale]: SeoPageProps } = {
  en: {
    title: 'Tsuwari - Main page',
    description:
      'Very powerful and useful Twitch bot that helps manage chat on big channels. Developed from streamers for streamers with love.',
    keywords: ['tsuwari', 'twitch bot', 'tsuwari bot'],
  },
  ru: {
    title: 'Tsuwari - Главная страница',
    description:
      'Очень мощный и полезный бот для Twitch, который помогает управлять чатом на больших каналах. Разработан от стримеров для стримеров c любовью.',
    keywords: ['tsuwari', 'twitch bot', 'tsuwari bot'],
  },
};

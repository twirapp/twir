import type IAppLocale from '@/locales/app/interface.js';
import { AppMenu } from '@/pages/app/router.js';

const messages: IAppLocale = {
  hello: 'Привет',
  pages: {
    [AppMenu.dashboard]: 'Панель управления',
    [AppMenu.commands]: 'Комманды',
    [AppMenu.greetings]: 'Приветствия',
    [AppMenu.keywords]: 'Ключевые слова',
    [AppMenu.moderation]: 'Модерация',
    [AppMenu.settings]: 'Настройки',
    [AppMenu.timers]: 'Таймеры',
    [AppMenu.variables]: 'Переменные',
  },
};

export default messages;

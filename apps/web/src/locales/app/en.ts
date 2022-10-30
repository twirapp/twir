import type IAppLocale from '@/locales/app/interface.js';
import { AppMenu } from '@/pages/app/router.js';

const messages: IAppLocale = {
  hello: 'Hello',
  pages: {
    [AppMenu.dashboard]: 'Dashboard',
    [AppMenu.commands]: 'Commands',
    [AppMenu.greetings]: 'Greetings',
    [AppMenu.keywords]: 'Keywords',
    [AppMenu.moderation]: 'Moderation',
    [AppMenu.settings]: 'Settings',
    [AppMenu.timers]: 'Timers',
    [AppMenu.variables]: 'Variables',
  },
};

export default messages;

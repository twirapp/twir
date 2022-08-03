import { DefaultCommand } from '../types.js';

export const spam: DefaultCommand[] = [
  {
    name: 'spam',
    description: 'Spam into chat. Example usage: <b>!spam 5 https://tsuwari.tk</b>',
    permission: 'MODERATOR',
    visible: false,
    module: 'CHANNEL',
    handler(state, params?) {
      if (!params || !state.channelId) return;

      const paramsArr = params.split(' ');

      if (paramsArr.length < 2) return;
      const count = Number(paramsArr[0]);

      if (isNaN(count)) return;
      if (count > 10) return 'Max count is 10.';

      return Array(count).fill(paramsArr.slice(1).join(' '));
    },
  },
];
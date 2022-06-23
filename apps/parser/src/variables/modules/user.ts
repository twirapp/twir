import { TwitchApiService } from '@tsuwari/shared';
import { formatDuration, intervalToDuration } from 'date-fns';

import { app } from '../../index.js';
import { Module } from '../index.js';

const formatDistanceLocale: { [x: string]: string } = {
  xMinutes: '{{count}}m',
  xHours: '{{count}}h',
};
const shortEnLocale = { formatDistance: (token: string, count: string) => formatDistanceLocale[token]?.replace('{{count}}', count) };

const staticApi = app.get(TwitchApiService);

async function getuserFollowAge(fromId: string, toId: string) {
  const follow = await staticApi.users.getFollowFromUserToBroadcaster(fromId, toId);

  if (!follow) return 'not follower';

  const duration = intervalToDuration({ start: follow.followDate.getTime(), end: Date.now() });

  return formatDuration(duration);
}


export const user: Module[] = [
  {
    key: 'user.followage',
    description: 'User followage',
    handler: (_, state) => {
      if (!state.sender?.id || !state.channelId) return 'cannot fetch data';
      return getuserFollowAge(state.sender.id, state.channelId);
    },
  },
  {
    key: 'user.messages',
    description: 'User messages',
    handler: async (_, state) => {
      const stats = await state.cache.getUserStats();

      return stats?.messages ?? 0;
    },
  },
  {
    key: 'user.watched',
    description: 'User watched time',
    handler: async (_, state) => {
      const stats = await state.cache.getUserStats();
      if (!stats) return '0h0m';
      return formatDuration(
        intervalToDuration({ start: 0, end: Number(stats.watched ?? 0) }),
        {
          zero: true,
          format: ['hours', 'minutes'],
          delimiter: '',
          locale: shortEnLocale,
        },
      );
    },
  },
];

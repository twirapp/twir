import { RedisService, TwitchApiService, WEEK } from '@tsuwari/shared';
import { HelixUserData } from '@twurple/api';
import { getRawData } from '@twurple/common';
import { formatDuration, intervalToDuration } from 'date-fns';

import { app } from '../../index.js';
import { Module } from '../index.js';

const formatDistanceLocale: { [x: string]: string } = {
  xMinutes: '{{count}}m',
  xHours: '{{count}}h',
};
const shortEnLocale = { formatDistance: (token: string, count: string) => formatDistanceLocale[token]?.replace('{{count}}', count) };

const staticApi = app.get(TwitchApiService);
const redis = app.get(RedisService);

async function getuserFollowAge(fromId: string, toId: string) {
  const follow = await staticApi.users.getFollowFromUserToBroadcaster(fromId, toId);

  if (!follow) return 'not follower';

  const duration = intervalToDuration({ start: follow.followDate.getTime(), end: Date.now() });

  return formatDuration(duration);
}

async function getUserAge(userId: string) {
  const redisKey = `twitchUsersCache:${userId}`;
  let data: HelixUserData | null = null;

  const cachedData = await redis.get(redisKey);
  if (cachedData) {
    data = JSON.parse(cachedData) as HelixUserData;
  } else {
    const user = await staticApi.users.getUserById(userId);
    if (user) {
      data = getRawData(user);
      redis.set(redisKey, JSON.stringify(data), 'EX', (WEEK * 2) / 1000);
    }
  }

  if (!data) return 'Error on getting info.';

  const duration = intervalToDuration({ start: new Date(data.created_at).getTime(), end: Date.now() });
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
  {
    key: 'user.age',
    description: 'User account age',
    handler(key, state) {
      if (!state.sender?.id || !state.channelId) return 'cannot fetch data';
      return getUserAge(state.sender.id);
    },
  },
];

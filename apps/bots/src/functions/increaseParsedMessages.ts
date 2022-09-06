import { HelixStreamData } from '@twurple/api';

import { redis } from '../libs/redis.js';

export const increaseParsedMessages = async (channelId: string) => {
  const key = `streams:${channelId}`;
  const rawStream = await redis.get(`streams:${channelId}`);
  if (!rawStream) return;

  const stream = JSON.parse(rawStream) as HelixStreamData & { parsedMessages?: number };

  stream.parsedMessages = (stream.parsedMessages ?? 0) + 1;

  redis.set(key, JSON.stringify(stream));
};

import { Stream, streamSchema } from '@tsuwari/redis';

import { redisOm } from '../libs/redis.js';

const repository = redisOm.fetchRepository(streamSchema);

export const increaseParsedMessages = async (channelId: string) => {
  const stream = await repository.fetch(channelId);
  const json = stream.toRedisJson();
  if (Object.keys(json).length) {
    await repository.createAndSave({
      ...json as Stream,
      parsedMessages: (json.parsedMessages ?? 0) + 1,
    }, channelId);
  }
};
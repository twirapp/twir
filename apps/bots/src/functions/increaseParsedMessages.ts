import { StreamRepository } from '@tsuwari/redis';

import { redisSource, redlock } from '../libs/redis.js';

const repository = redisSource.getRepository(StreamRepository);

export const increaseParsedMessages = async (channelId: string) => {
  await redlock.using(
    [`increaseParsedMessages:${channelId}`],
    5000,
    { retryCount: 2 },
    async () => {
      console.log('increasing');
      const stream = await repository.read(channelId);
      if (!stream) return;

      stream.parsedMessages = (stream.parsedMessages ?? 0) + 1;

      repository.write(channelId, stream);
    },
  );
};

import { redis } from '../libs/redis.js';

export const increaseParsedMessages = async (channelId: string) => {
  const stream = await redis.get(`streams:${channelId}`);

  if (stream) {
    const streamData = JSON.parse(stream);
    await redis.set(`streams:${channelId}`, JSON.stringify({
      ...streamData,
      parsedMessages: (streamData.parsedMessages ?? 0) + 1,
    }));
  }
};
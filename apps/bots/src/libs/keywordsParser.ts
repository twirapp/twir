import { KeywordsRepository } from '@tsuwari/redis';
import { TwitchPrivateMessage } from '@twurple/chat/lib/commands/TwitchPrivateMessage.js';

import { redis, redisSource } from './redis.js';

export class KeywordsParser {
  async parse(message: string, state: TwitchPrivateMessage) {
    if (!state.channelId) return;
    const keywordsKeys = await redis.keys(`keywords:${state.channelId}:*`);
    if (!keywordsKeys?.length) return;

    const keywords = await redisSource
      .getRepository(KeywordsRepository)
      .readMany(keywordsKeys, true);
    const responses: string[] = [];

    message = message.toLowerCase();
    for (const keyword of keywords.filter((k) => k.enabled)) {
      const cooldownKey = `cooldowns:keywords:${keyword.id}`;
      const isOnCooldown = await redis.get(cooldownKey);
      if (message.includes(keyword.text.toLowerCase())) {
        if (keyword.cooldown && isOnCooldown) continue;
        responses.push(keyword.response);

        if (keyword.cooldown !== null) {
          redis.set(cooldownKey, 'true').then(() => redis.expire(cooldownKey, keyword.cooldown!));
        }
      }
    }

    return responses;
  }
}

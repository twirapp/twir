import { GreetingsRepository } from '@tsuwari/redis';
import type { TwitchPrivateMessage } from '@twurple/chat/lib/commands/TwitchPrivateMessage';

import { redisSource } from './redis.js';

export class GreetingsParser {
  async parse(state: TwitchPrivateMessage) {
    const repository = redisSource.getRepository(GreetingsRepository);

    const key = `${state.channelId}:${state.userInfo.userId}`;
    const item = await repository.read(key);

    if (!item || item.processed !== false || !item.enabled) return;

    await repository.write(key, {
      ...item,
      processed: true,
    });

    return item.text;
  }
}

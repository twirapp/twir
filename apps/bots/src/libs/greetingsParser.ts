import { Greeting } from '@tsuwari/prisma';
import type { TwitchPrivateMessage } from '@twurple/chat/lib/commands/TwitchPrivateMessage';

import { redis } from './redis.js';

type GreetingConditional = Greeting & { processed: 'true' | 'false' };

export class GreetingsParser {
  async parse(state: TwitchPrivateMessage) {
    const key = `greetings:${state.channelId}:${state.userInfo.userId}`;
    const item = await redis.hgetall(key) as unknown as GreetingConditional;

    if (!Object.keys(item).length || item.processed !== 'false') return;

    await redis.hset(key, 'processed', 'true');
    return item.text;
  }
}

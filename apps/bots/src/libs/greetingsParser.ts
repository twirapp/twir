import { greetingsSchema } from '@tsuwari/redis';
import type { TwitchPrivateMessage } from '@twurple/chat/lib/commands/TwitchPrivateMessage';

import { redisOm } from './redis.js';


export class GreetingsParser {
  #repository = redisOm.fetchRepository(greetingsSchema);

  async parse(state: TwitchPrivateMessage) {
    const key = `${state.channelId}:${state.userInfo.userId}`;
    const item = await this.#repository.fetch(key).then(i => i.toRedisJson());

    if (!Object.keys(item).length || item.processed !== false || !item.enabled) return;

    await this.#repository.createAndSave({
      ...item,
      processed: true,
    }, key);

    return item.text;
  }
}

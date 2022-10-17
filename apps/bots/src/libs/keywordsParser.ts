import { ChannelKeyword } from '@tsuwari/typeorm/entities/ChannelKeyword';
import { TwitchPrivateMessage } from '@twurple/chat/lib/commands/TwitchPrivateMessage.js';

import { typeorm } from './typeorm.js';

const repository = typeorm.getRepository(ChannelKeyword);

export class KeywordsParser {
  async parse(message: string, state: TwitchPrivateMessage) {
    if (!state.channelId) return;

    const keywords = await repository.findBy({
      channelId: state.channelId,
      enabled: true,
    });

    const responses: string[] = [];

    message = message.toLowerCase();
    for (const keyword of keywords) {
      if (!message.includes(keyword.text.toLowerCase())) {
        continue;
      }
      let isOnCooldown = false;
      if (keyword.cooldown && keyword.cooldownExpireAt) {
        isOnCooldown = keyword.cooldownExpireAt.getTime() >= Date.now();
      }

      if (isOnCooldown) continue;
      responses.push(keyword.response);

      if (keyword.cooldown !== null) {
        await repository.update(
          { id: keyword.id },
          { cooldownExpireAt: new Date(Date.now() + 1000 * keyword.cooldown) },
        );
      }
    }

    return responses;
  }
}

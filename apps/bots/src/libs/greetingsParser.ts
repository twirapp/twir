import { ChannelGreeting } from '@tsuwari/typeorm/entities/ChannelGreeting';
import type { TwitchPrivateMessage } from '@twurple/chat/lib/commands/TwitchPrivateMessage';

import { typeorm } from './typeorm.js';

const repository = typeorm.getRepository(ChannelGreeting);

export class GreetingsParser {
  async parse(state: TwitchPrivateMessage) {
    const item = await repository.findOneBy({
      channelId: state.channelId!,
      userId: state.userInfo.userId!,
    });

    if (item?.processed || !item?.enabled) return;

    await repository.update({ id: item.id }, { processed: true });

    return {
      text: item.text,
      isReply: item.isReply,
    };
  }
}

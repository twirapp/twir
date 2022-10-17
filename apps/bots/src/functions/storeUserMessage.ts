import { messageSchema } from '@tsuwari/redis';
import { TwitchPrivateMessage } from '@twurple/chat/lib/commands/TwitchPrivateMessage.js';

import { redisOm } from '../libs/redis.js';

const repository = redisOm.fetchRepository(messageSchema);

export function storeUserMessage(state: TwitchPrivateMessage, message: string) {
  const key = `${state.userInfo.userId}:${state.id}`;
  repository
    .createAndSave(
      {
        channelId: state.channelId!,
        userId: state.userInfo.userId!,
        message,
        messageId: state.id,
        userName: state.userInfo.userName,
        canBeDeleted:
          !state.userInfo.isSubscriber &&
          !state.userInfo.isFounder &&
          !state.userInfo.isVip &&
          !state.userInfo.isMod &&
          !state.userInfo.isBroadcaster,
      },
      key,
    )
    .then(() => {
      repository.expire(key, 600);
    })
    .catch(console.error);
}

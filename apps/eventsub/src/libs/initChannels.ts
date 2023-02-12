import { config } from '@tsuwari/config';
import { Channel } from '@tsuwari/typeorm/entities/Channel';
import { Token } from '@tsuwari/typeorm/entities/Token';

import { botsGrpcClient } from './botsGrpc.js';
import {
  apiClient,
  eventSubMiddleware,
  subscribeToEvents, unSubscribeFromEvents,
} from './middleware.js';
import { typeorm } from './typeorm.js';

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
export const initChannels = async (force = false) => {
  const channels = await typeorm.getRepository(Channel).find({
    where: {
      isEnabled: true,
      isBanned: false,
    },
  });

  if (config.isProd || force) {
    for (const channel of channels) {
      subscribeToEvents(channel.id);
    }
  } else {
    await apiClient.eventSub.deleteAllSubscriptions();

    return initChannels(true);
  }

  eventSubMiddleware.subscribeToUserAuthorizationRevokeEvents(config.TWITCH_CLIENTID, async (e) => {
    const repository = typeorm.getRepository(Channel);
    const channel = await repository.findOneBy({
      id: e.userId,
    });

    if (channel) {
      let username: string | null | undefined = null;
      username = e.userName;

      if (!username) {
        const user = await apiClient.users.getUserById(e.userId);
        username = user?.name;
      }

      if (e.userName) {
        botsGrpcClient.leave({
          userName: username,
          botId: channel.botId,
        });
      }
      await repository.update({ id: channel.id }, { isEnabled: false });

      const tokenRepository = typeorm.getRepository(Token);
      const token = await tokenRepository.findOneBy({
        user: {
          id: channel.id,
        },
      });
      if (token) {
        await tokenRepository.delete({
          id: token.id,
        });
      }
      unSubscribeFromEvents(e.userId);
    }
  });
};

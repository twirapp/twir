import { config } from '@tsuwari/config';
import { Channel } from '@tsuwari/typeorm/entities/Channel';

import {
  apiClient,
  eventSubHandlers,
  eventSubMiddleware,
  subscribeToEvents,
} from './middleware.js';
import { typeorm } from './typeorm.js';

export const initChannels = async (force = false) => {
  const channels = await typeorm.getRepository(Channel).find();
  if (config.isProd || force) {
    for (const channel of channels) {
      subscribeToEvents(channel.id);
    }
  } else {
    await apiClient.eventSub.deleteAllSubscriptions();

    initChannels(true);
  }

  eventSubMiddleware
    .subscribeToUserAuthorizationRevokeEvents(
      config.TWITCH_CLIENTID,
      eventSubHandlers.subscribeToUserAuthorizationRevokeEvents,
    )
    // eslint-disable-next-line @typescript-eslint/no-empty-function
    .catch(() => {});
};

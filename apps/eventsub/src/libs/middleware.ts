import { config } from '@tsuwari/config';
import { Channel } from '@tsuwari/typeorm/entities/Channel';
import { ChannelEvent, EventType } from '@tsuwari/typeorm/entities/ChannelEvent';
import { ChannelFollowEvent } from '@tsuwari/typeorm/entities/channelEvents/Follow';
import { Token } from '@tsuwari/typeorm/entities/Token';
import { ApiClient } from '@twurple/api';
import { ClientCredentialsAuthProvider } from '@twurple/auth';
import { getRawData } from '@twurple/common';
import {
  EventSubChannelFollowEvent,
  EventSubChannelModeratorEvent,
  EventSubChannelUpdateEvent,
  EventSubMiddleware,
  EventSubStreamOfflineEvent,
  EventSubStreamOnlineEvent,
  EventSubSubscription,
  EventSubUserAuthorizationRevokeEvent,
  EventSubUserUpdateEvent,
} from '@twurple/eventsub';

import { botsGrpcClient } from './botsGrpc.js';
import { getHostName } from './hostname.js';
import { pubsub } from './pubsub.js';
import { typeorm } from './typeorm.js';

const subScriptionValues = new Map([
  ['channel.update', 'subscribeToChannelUpdateEvents'],
  ['stream.online', 'subscribeToStreamOnlineEvents'],
  ['stream.offline', 'subscribeToStreamOfflineEvents'],
  ['user.update', 'subscribeToUserUpdateEvents'],
  ['channel.follow', 'subscribeToChannelFollowEvents'],
  ['channel.moderator.add', 'subscribeToChannelModeratorAddEvents'],
  ['channel.moderator.remove', 'subscribeToChannelModeratorRemoveEvents'],
]);

const authProvider = new ClientCredentialsAuthProvider(
  config.TWITCH_CLIENTID,
  config.TWITCH_CLIENTSECRET,
);
export const apiClient = new ApiClient({ authProvider });

export const eventSubMiddleware = new EventSubMiddleware({
  apiClient,
  hostName: await getHostName(),
  pathPrefix: '/twitch',
  secret: 'secretHere',
  logger: {
    minLevel: 'debug',
  },
  strictHostCheck: true,
});

const subscriptions: Map<string, EventSubSubscription[]> = new Map();

export const subscribeToEvents = (channelId: string) => {
  if (!subscriptions.has(channelId)) {
    subscriptions.set(channelId, []);
  }
  const cachedChannel = subscriptions.get(channelId);

  for (const type of subScriptionValues.keys()) {
    const typeValue = subScriptionValues.get(type);
    if (!typeValue) continue;
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore
    eventSubMiddleware[typeValue](channelId, (e) => {
      // eslint-disable-next-line @typescript-eslint/ban-ts-comment
      // @ts-ignore
      // eslint-disable-next-line @typescript-eslint/no-empty-function
      eventSubHandlers[typeValue](e);
    })
      .then((s: EventSubSubscription) => {
        subscriptions.set(channelId, [
          ...cachedChannel!,
          s,
        ]);
      })
      .catch((e: any) => {
        console.log(`${typeValue}#${channelId}`,  e);
      });

    console.log(`Subscribed to ${type}#${channelId} event.`);
  }
};

export const unSubscribeFromEvents = (channelId: string) => {
  const cachedChannel = subscriptions.get(channelId);
  if (!cachedChannel) return;

  for (const sub of cachedChannel) {
    sub.stop().catch(e => console.error(e));
  }
};

export const eventSubHandlers = {
  subscribeToChannelUpdateEvents: (e: EventSubChannelUpdateEvent) => {
    pubsub.publish('stream.update', getRawData(e));
  },
  subscribeToStreamOnlineEvents: async (e: EventSubStreamOnlineEvent) => {
    const stream = await e.getStream();

    pubsub.publish('streams.online', { channelId: e.broadcasterId, streamId: stream.id });
  },
  subscribeToStreamOfflineEvents: (e: EventSubStreamOfflineEvent) => {
    pubsub.publish('streams.offline', { channelId: e.broadcasterId });
  },
  subscribeToUserUpdateEvents: (e: EventSubUserUpdateEvent) => {
    pubsub.publish('user.update', getRawData(e));
  },
  subscribeToUserAuthorizationRevokeEvents: async (e: EventSubUserAuthorizationRevokeEvent) => {
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
          userName: e.userName,
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
  },
  subscribeToChannelFollowEvents: async (e: EventSubChannelFollowEvent) => {
    console.info(
      `New follow from ${e.userName}#${e.userId} to ${e.broadcasterName}#${e.broadcasterId}`,
    );
    const event = await typeorm.getRepository(ChannelEvent).save({
      type: EventType.FOLLOW,
      channelId: e.broadcasterId,
    });
    await typeorm.getRepository(ChannelFollowEvent).save({
      eventId: event.id,
      fromUserId: e.userId,
      toUserId: e.broadcasterId,
    });
  },
  subscribeToChannelModeratorAddEvents: (e: EventSubChannelModeratorEvent) => {
    updateBotModStatus(e, true);
  },
  subscribeToChannelModeratorRemoveEvents: (e: EventSubChannelModeratorEvent) => {
    updateBotModStatus(e, false);
  },
};

const updateBotModStatus = async (e: EventSubChannelModeratorEvent, state: boolean) => {
  const repository = typeorm.getRepository(Channel);
  const channel = await repository.findOneBy({
    id: e.broadcasterId,
  });

  if (!channel || channel.botId !== e.userId) return;
  await repository.update({ id: channel.id }, { isBotMod: state });
};

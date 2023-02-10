import { config } from '@tsuwari/config';
import { Channel } from '@tsuwari/typeorm/entities/Channel';
import { ChannelCommand } from '@tsuwari/typeorm/entities/ChannelCommand';
import { ChannelEvent, EventType } from '@tsuwari/typeorm/entities/ChannelEvent';
import { ChannelFollowEvent } from '@tsuwari/typeorm/entities/channelEvents/Follow';
import { ChannelModuleSettings, ModuleType } from '@tsuwari/typeorm/entities/ChannelModuleSettings';
import { Token } from '@tsuwari/typeorm/entities/Token';
import { type YouTubeSettings } from '@tsuwari/types/api';
import { ApiClient } from '@twurple/api';
import { ClientCredentialsAuthProvider } from '@twurple/auth';
import { getRawData } from '@twurple/common';
import {
  EventSubChannelFollowEvent,
  EventSubChannelModeratorEvent,
  EventSubChannelRedemptionAddEvent,
  EventSubChannelUpdateEvent,
  EventSubMiddleware,
  EventSubStreamOfflineEvent,
  EventSubStreamOnlineEvent,
  EventSubSubscription,
  EventSubUserAuthorizationRevokeEvent,
  EventSubUserUpdateEvent,
} from '@twurple/eventsub';

import { botsGrpcClient } from './botsGrpc.js';
import { eventsGrpcClient } from './eventsGrpc.js';
import { getHostName } from './hostname.js';
import { parserGrpcClient } from './parserGrpc.js';
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
  ['channel.channel_points_custom_reward_redemption.add', 'subscribeToChannelRedemptionAddEvents'],
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
        subscriptions.set(channelId, [...cachedChannel!, s]);
      })
      .catch((e: any) => {
        console.log(`${typeValue}#${channelId}`, e);
      });

    console.log(`Subscribed to ${type}#${channelId} event.`);
  }
};

export const unSubscribeFromEvents = (channelId: string) => {
  const cachedChannel = subscriptions.get(channelId);
  if (!cachedChannel) return;

  for (const sub of cachedChannel) {
    sub.stop().catch((e) => console.error(e));
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
    await eventsGrpcClient.follow({
      userName: e.userName,
      userDisplayName: e.userDisplayName,
      baseInfo: {
        channelId: e.broadcasterId,
      },
    });
  },
  subscribeToChannelModeratorAddEvents: (e: EventSubChannelModeratorEvent) => {
    updateBotModStatus(e, true);
  },
  subscribeToChannelModeratorRemoveEvents: (e: EventSubChannelModeratorEvent) => {
    updateBotModStatus(e, false);
  },
  subscribeToChannelRedemptionAddEvents: async (e: EventSubChannelRedemptionAddEvent) => {
    await eventsGrpcClient.redemptionCreated({
      id: e.rewardId,
      baseInfo: { channelId: e.broadcasterId },
      input: e.input,
      userName: e.userName,
      userDisplayName: e.userDisplayName,
      rewardCost: e.rewardCost.toString(),
      rewardName: e.rewardTitle,
    });

    if (!e.input) return;

    const repository = typeorm.getRepository(ChannelModuleSettings);
    const entity = await repository.findOne({
      where: {
        channelId: e.broadcasterId,
        type: ModuleType.YOUTUBE_SONG_REQUESTS,
      },
    });
    if (!entity) return;
    const settings = entity.settings as YouTubeSettings;
    if (!settings.enabled || e.rewardId !== settings.channelPointsRewardId) return;

    const command = await typeorm.getRepository(ChannelCommand).findOneBy({
      channelId: e.broadcasterId,
      defaultName: 'ytsr',
    });

    if (!command) return;

    const result = await parserGrpcClient.processCommand({
      channel: {
        id: e.broadcasterId,
        name: e.broadcasterName,
      },
      message: {
        id: e.id,
        text: `!${command.name} ${e.input}`,
      },
      sender: {
        id: e.userId,
        badges: ['VIEWER'],
        displayName: e.userDisplayName,
        name: e.userName,
      },
    });

    botsGrpcClient.sendMessage({
      channelId: e.broadcasterId,
      channelName: e.broadcasterName,
      isAnnounce: false,
      message: `@${e.userName} ${result.responses[0]}`,
    });
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

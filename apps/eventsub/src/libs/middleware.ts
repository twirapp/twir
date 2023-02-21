import { config } from '@tsuwari/config';
import { Channel } from '@tsuwari/typeorm/entities/Channel';
import { ChannelCommand } from '@tsuwari/typeorm/entities/ChannelCommand';
import { ChannelEvent, EventType } from '@tsuwari/typeorm/entities/ChannelEvent';
import { ChannelFollowEvent } from '@tsuwari/typeorm/entities/channelEvents/Follow';
import { ChannelModuleSettings, ModuleType } from '@tsuwari/typeorm/entities/ChannelModuleSettings';
import { type YouTubeSettings } from '@tsuwari/types/api';
import { ApiClient } from '@twurple/api';
import { ClientCredentialsAuthProvider } from '@twurple/auth';
import { getRawData } from '@twurple/common';
import {
  EventSubChannelModeratorEvent,
  EventSubMiddleware,
  EventSubSubscription,
} from '@twurple/eventsub';

import { countUserChannelPoints } from '../helpers/countUserChannelPoints.js';
import { decrementUserChannelPoints } from '../helpers/decrementUserChannelPoints.js';
import { botsGrpcClient } from './botsGrpc.js';
import { eventsGrpcClient } from './eventsGrpc.js';
import { getHostName } from './hostname.js';
import { parserGrpcClient } from './parserGrpc.js';
import { pubsub } from './pubsub.js';
import { typeorm } from './typeorm.js';

const authProvider = new ClientCredentialsAuthProvider(
  config.TWITCH_CLIENTID,
  config.TWITCH_CLIENTSECRET,
);
export const apiClient = new ApiClient({ authProvider });

export const eventSubMiddleware = new EventSubMiddleware({
  apiClient,
  hostName: await getHostName(),
  pathPrefix: '/twitch',
  secret: config.EVENTSUB_SECRET,
  logger: {
    minLevel: 'debug',
  },
  strictHostCheck: true,
});

const subscriptions: Map<string, UserSubscriptions> = new Map();

class UserSubscriptions {
  #subs: Array<EventSubSubscription | undefined> = [];

  constructor(private readonly userId: string) {}

  async init() {
    this.#subs = (await Promise.allSettled([
      this.#subscribeToChannelUpdateEvents(),
      this.#subscribeToStreamOnlineEvents(),
      this.#subscribeToStreamOfflineEvents(),
      this.#subscribeToUserUpdateEvents(),
      this.#subscribeToChannelFollowEvents(),
      this.#subscribeToChannelModeratorAddEvents(),
      this.#subscribeToChannelModeratorRemoveEvents(),
      this.#subscribeToChannelRedemptionAddEvents(),
      this.#subscribeToChannelRedemptionUpdateEvents(),
    ])).map(p => p.status === 'rejected' ? undefined : p.value);
  }

  async unsubscribe() {
    await Promise.all(this.#subs.map(s => s?.stop()));
  }

  #subscribeToChannelUpdateEvents() {
    return eventSubMiddleware.subscribeToChannelUpdateEvents(this.userId, (e) => {
      pubsub.publish('stream.update', getRawData(e));

      eventsGrpcClient.titleOrCategoryChanged({
        baseInfo: { channelId: e.broadcasterId },
        newCategory: e.categoryName,
        newTitle: e.streamTitle,
      });
    });
  }

  #subscribeToStreamOnlineEvents() {
    return eventSubMiddleware.subscribeToStreamOnlineEvents(this.userId, async (e) => {
      const stream = await e.getStream();

      pubsub.publish('streams.online', { channelId: e.broadcasterId, streamId: stream.id });
      await eventsGrpcClient.streamOnline({
        baseInfo: { channelId: e.broadcasterId },
        category: stream.gameName,
        title: stream.title,
      });
    });
  }

  #subscribeToStreamOfflineEvents() {
    return eventSubMiddleware.subscribeToStreamOfflineEvents(this.userId, async (e) => {
      pubsub.publish('streams.offline', { channelId: e.broadcasterId });
      await eventsGrpcClient.streamOffline({
        baseInfo: { channelId: e.broadcasterId },
      });
    });
  }

  #subscribeToUserUpdateEvents() {
    return eventSubMiddleware.subscribeToUserUpdateEvents(this.userId, async (e) => {
      pubsub.publish('user.update', getRawData(e));
    });
  }

  #subscribeToChannelFollowEvents() {
    return eventSubMiddleware.subscribeToChannelFollowEvents(this.userId, async (e) => {
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
    });
  }

  #subscribeToChannelModeratorAddEvents() {
    return eventSubMiddleware.subscribeToChannelModeratorAddEvents(this.userId, async (e) => {
      await updateBotModStatus(e, true);
    });
  }

  #subscribeToChannelModeratorRemoveEvents() {
    return eventSubMiddleware.subscribeToChannelModeratorRemoveEvents(this.userId, async (e) => {
      await updateBotModStatus(e, false);
    });
  }

  #subscribeToChannelRedemptionAddEvents() {
    return eventSubMiddleware.subscribeToChannelRedemptionAddEvents(this.userId, async (e) => {
      countUserChannelPoints(e.broadcasterId, e.userId, e.rewardCost);
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

      await botsGrpcClient.sendMessage({
        channelId: e.broadcasterId,
        channelName: e.broadcasterName,
        isAnnounce: false,
        message: `@${e.userName} ${result.responses[0]}`,
      });
    });
  }

  #subscribeToChannelRedemptionUpdateEvents() {
    return eventSubMiddleware.subscribeToChannelRedemptionUpdateEvents(this.userId, (data) => {
      if (data.status === 'canceled') {
        decrementUserChannelPoints(data.broadcasterId, data.userId, data.rewardCost);
      }
    });
  }
}

export const subscribeToEvents = async (channelId: string) => {
  if (subscriptions.has(channelId)) {
    const subscription = await subscriptions.get(channelId)!;
    await subscription.unsubscribe();
    subscriptions.delete(channelId);
  }

  const subs = new UserSubscriptions(channelId);
  await subs.init();
  subscriptions.set(channelId, subs);
};

export const unSubscribeFromEvents = async (channelId: string) => {
  const cachedChannel = subscriptions.get(channelId);
  if (!cachedChannel) return;

  await cachedChannel.unsubscribe();
};

const updateBotModStatus = async (e: EventSubChannelModeratorEvent, state: boolean) => {
  const repository = typeorm.getRepository(Channel);
  const channel = await repository.findOneBy({
    id: e.broadcasterId,
  });

  if (!channel || channel.botId !== e.userId) return;
  await repository.update({ id: channel.id }, { isBotMod: state });
};

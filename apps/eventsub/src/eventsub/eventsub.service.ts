import { Inject, Injectable, Logger } from '@nestjs/common';
import { config } from '@tsuwari/config';
import { PrismaService } from '@tsuwari/prisma';
import { TwitchApiService } from '@tsuwari/shared';
import { EventSubMiddleware } from '@twurple/eventsub';

import { HandlerService } from '../handler/handler.service.js';

const subScriptionValues = new Map([
  ['channel.update', 'subscribeToChannelUpdateEvents'],
  ['stream.online', 'subscribeToStreamOnlineEvents'],
  ['stream.offline', 'subscribeToStreamOfflineEvents'],
]);

@Injectable()
export class EventSub extends EventSubMiddleware {
  private readonly logger = new Logger(EventSub.name);

  constructor(
    readonly twitchApi: TwitchApiService,
    @Inject('HOSTNAME') hostName: string,
    private readonly prisma: PrismaService,
    private readonly handler: HandlerService,
  ) {
    super({
      apiClient: twitchApi,
      hostName: hostName,
      pathPrefix: '/twitch',
      secret: 'secretHere',
      logger: {
        minLevel: 'debug',
      },
      strictHostCheck: true,
    });
  }

  async subscribeToEvents(channelId: string) {
    for (const type of subScriptionValues.keys()) {
      const typeValue = subScriptionValues.get(type);
      if (!typeValue) continue;

      // eslint-disable-next-line @typescript-eslint/ban-ts-comment
      // @ts-ignore
      this[typeValue](channelId, (e) => {
        // eslint-disable-next-line @typescript-eslint/ban-ts-comment
        // @ts-ignore
        this.handler[typeValue](e);
      }).catch();

      this.logger.log(`Subsribed to ${type}#${channelId} event.`);
    }
  }

  async init(force = false) {
    const channels = await this.prisma.channel.findMany();
    if (config.isProd || force) {
      //const subscriptions = await this.twitchApi.eventSub.getSubscriptionsPaginated().getAll();

      for (const channel of channels) {
        /* for (const type of subScriptionValues.keys()) {
          const isExists = subscriptions.find(s => s.type === type && s.condition.broadcaster_user_id === channel.id);
          if (isExists) continue;
          const typeValue = subScriptionValues.get(type);
          if (!typeValue) continue;

          // eslint-disable-next-line @typescript-eslint/ban-ts-comment
          // @ts-ignore
          this[typeValue](channel.id, (e) => {
            // eslint-disable-next-line @typescript-eslint/ban-ts-comment
            // @ts-ignore
            this.handler[typeValue](e);
          });

          this.logger.log(`Subsribed to ${type}#${channel.id} event.`);
        } */
        this.subscribeToEvents(channel.id);
      }

    } else {
      await this.twitchApi.eventSub.deleteAllSubscriptions();

      this.init(true);
    }
  }
} 
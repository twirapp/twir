import { Inject, Injectable, Logger } from '@nestjs/common';
import { config } from '@tsuwari/config';
import { TwitchApiService } from '@tsuwari/shared';
import { Channel } from '@tsuwari/typeorm/entities/Channel';
import { EventSubMiddleware } from '@twurple/eventsub';

import { HandlerService } from '../handler/handler.service.js';
import { typeorm } from '../index.js';

// eslint-disable-next-line @typescript-eslint/no-empty-function
const noop = () => {};

const subScriptionValues = new Map([
  ['channel.update', 'subscribeToChannelUpdateEvents'],
  ['stream.online', 'subscribeToStreamOnlineEvents'],
  ['stream.offline', 'subscribeToStreamOfflineEvents'],
  ['user.update', 'subscribeToUserUpdateEvents'],
  ['channel.follow', 'subscribeToChannelFollowEvents'],
]);

@Injectable()
export class EventSub extends EventSubMiddleware {
  private readonly logger = new Logger(EventSub.name);

  constructor(
    readonly twitchApi: TwitchApiService,
    @Inject('HOSTNAME') hostName: string,
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
      }).catch(noop);

      this.logger.log(`Subsribed to ${type}#${channelId} event.`);
    }
  }

  async init(force = false) {
    const channels = await typeorm.getRepository(Channel).find();
    if (config.isProd || force) {
      for (const channel of channels) {
        this.subscribeToEvents(channel.id);
      }
    } else {
      await this.twitchApi.eventSub.deleteAllSubscriptions();

      this.init(true);
    }

    this.subscribeToUserAuthorizationRevokeEvents(
      config.TWITCH_CLIENTID,
      this.handler.subscribeToUserAuthorizationRevokeEvents,
    ).catch(noop);
  }
}

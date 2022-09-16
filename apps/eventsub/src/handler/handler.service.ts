import { Inject, Injectable } from '@nestjs/common';
import { ClientProxy, TwitchApiService } from '@tsuwari/shared';
import { Channel } from '@tsuwari/typeorm/entities/Channel';
import { Token } from '@tsuwari/typeorm/entities/Token';
import { getRawData } from '@twurple/common';
import {
  EventSubChannelUpdateEvent,
  EventSubStreamOfflineEvent,
  EventSubStreamOnlineEvent,
  EventSubUserAuthorizationRevokeEvent,
  EventSubUserUpdateEvent,
} from '@twurple/eventsub';
import { lastValueFrom } from 'rxjs';

import { typeorm } from '../index.js';

@Injectable()
export class HandlerService {
  constructor(
    @Inject('NATS') private readonly nats: ClientProxy,
    private readonly twitch: TwitchApiService,
  ) {}

  async subscribeToChannelUpdateEvents(e: EventSubChannelUpdateEvent) {
    await this.nats.emit('stream.update', getRawData(e)).toPromise();
  }

  async subscribeToStreamOnlineEvents(e: EventSubStreamOnlineEvent) {
    const stream = await e.getStream();

    await this.nats
      .emit('streams.online', {
        channelId: e.broadcasterId,
        streamId: stream.id,
      })
      .toPromise();
  }

  async subscribeToStreamOfflineEvents(e: EventSubStreamOfflineEvent) {
    await this.nats
      .emit('streams.offline', {
        channelId: e.broadcasterId,
      })
      .toPromise();
  }

  async subscribeToUserUpdateEvents(e: EventSubUserUpdateEvent) {
    await this.nats.emit('user.update', getRawData(e)).toPromise();
  }

  async subscribeToUserAuthorizationRevokeEvents(e: EventSubUserAuthorizationRevokeEvent) {
    const repository = typeorm.getRepository(Channel);
    const channel = await repository.findOneBy({
      id: e.userId,
    });

    if (channel) {
      let username: string | null | undefined = null;
      username = e.userName;

      if (!username) {
        const user = await this.twitch.users.getUserById(e.userId);
        username = user?.name;
      }

      if (e.userName) {
        await lastValueFrom(
          this.nats.emit('bots.joinOrLeave', {
            action: 'part',
            username: e.userName,
            botId: channel.botId,
          }),
        );
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
    }

    return;
  }
}

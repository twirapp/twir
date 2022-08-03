import { Inject, Injectable } from '@nestjs/common';
import { PrismaService } from '@tsuwari/prisma';
import { ClientProxy, TwitchApiService } from '@tsuwari/shared';
import { getRawData } from '@twurple/common';
import { EventSubChannelUpdateEvent, EventSubStreamOfflineEvent, EventSubStreamOnlineEvent, EventSubUserAuthorizationRevokeEvent, EventSubUserUpdateEvent } from '@twurple/eventsub';
import { lastValueFrom } from 'rxjs';


@Injectable()
export class HandlerService {
  constructor(
    @Inject('NATS') private readonly nats: ClientProxy,
    private readonly prisma: PrismaService,
    private readonly twitch: TwitchApiService,
  ) { }

  async subscribeToChannelUpdateEvents(e: EventSubChannelUpdateEvent) {
    await this.nats.emit('stream.update', getRawData(e)).toPromise();
  }

  async subscribeToStreamOnlineEvents(e: EventSubStreamOnlineEvent) {
    const stream = await e.getStream();

    await this.nats.emit('streams.online', {
      channelId: e.broadcasterId,
      streamId: stream.id,
    }).toPromise();
  }

  async subscribeToStreamOfflineEvents(e: EventSubStreamOfflineEvent) {
    await this.nats.emit('streams.offline', {
      channelId: e.broadcasterId,
    }).toPromise();
  }

  async subscribeToUserUpdateEvents(e: EventSubUserUpdateEvent) {
    await this.nats.emit('user.update', getRawData(e)).toPromise();
  }

  async subscribeToUserAuthorizationRevokeEvents(e: EventSubUserAuthorizationRevokeEvent) {
    const channel = await this.prisma.channel.findFirst({
      where: { id: e.userId },
    });

    if (channel) {
      let username: string | null | undefined = null;
      username = e.userName;

      if (!username) {
        const user = await this.twitch.users.getUserById(e.userId);
        username = user?.name;
      }

      if (e.userName) {
        await lastValueFrom(this.nats.emit('bots.joinOrLeave', { action: 'part', username: e.userName, botId: channel.botId }));
      }
      await this.prisma.channel.update({
        where: {
          id: channel.id,
        },
        data: {
          isEnabled: false,
        },
      });
      const token = await this.prisma.token.findFirst({
        where: { user: { id: channel.id } },
      });
      if (token) {
        await this.prisma.token.delete({
          where: {
            id: token.id,
          },
        });
      }
    }

    return;
  }
}

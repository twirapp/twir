import { Inject, Injectable } from '@nestjs/common';
import { ClientProxy } from '@tsuwari/shared';
import { getRawData } from '@twurple/common';
import { EventSubChannelUpdateEvent, EventSubStreamOfflineEvent, EventSubStreamOnlineEvent, EventSubUserUpdateEvent } from '@twurple/eventsub';


@Injectable()
export class HandlerService {
  constructor(@Inject('NATS') private readonly nats: ClientProxy) { }

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
}

import { Injectable } from '@nestjs/common';
import { Client, Transport } from '@nestjs/microservices';
import { config } from '@tsuwari/config';
import { ClientProxy, ClientProxyEvents, convertSnakeToCamel } from '@tsuwari/shared';
import { ChannelStream } from '@tsuwari/typeorm/entities/ChannelStream';
import { ApiClient } from '@twurple/api';
import { ClientCredentialsAuthProvider } from '@twurple/auth';
import { getRawData } from '@twurple/common';

import { typeorm } from './typeorm.js';

const authProvider = new ClientCredentialsAuthProvider(
  config.TWITCH_CLIENTID,
  config.TWITCH_CLIENTSECRET,
);
const api = new ApiClient({ authProvider });

@Injectable()
export class AppService {
  @Client({ transport: Transport.NATS, options: { servers: [config.NATS_URL] } })
  nats: ClientProxy;

  async delStream(userId: string) {
    await typeorm.getRepository(ChannelStream).delete({
      userId,
    });
  }

  async handleUpdate(e: ClientProxyEvents['stream.update']['input']) {
    const stream = await api.streams.getStreamByUserId(e.broadcaster_user_id);
    if (!stream) return;

    const repository = typeorm.getRepository(ChannelStream);
    const storedStream = await repository.findOneBy({ userId: e.broadcaster_user_id });

    if (storedStream) {
      await repository.update(
        { id: storedStream.id },
        {
          title: e.title,
          gameName: e.category_name,
          gameId: e.category_id,
        },
      );
    } else {
      await repository.save(convertSnakeToCamel(getRawData(stream)));
    }
  }

  async handleOnline(e: ClientProxyEvents['streams.online']['input']) {
    const stream = await api.streams.getStreamByUserId(e.channelId);
    if (!stream) return;
    const repository = typeorm.getRepository(ChannelStream);
    await repository.delete({
      userId: e.channelId,
    });
    await repository.save(convertSnakeToCamel(getRawData(stream)));
  }

  async handleOffline(e: ClientProxyEvents['streams.offline']['input']) {
    await typeorm.getRepository(ChannelStream).delete({
      userId: e.channelId,
    });
  }

  async handleChannels(channelsIds: string[]) {
    const { data: streams } = await api.streams.getStreams({
      userId: channelsIds,
    });

    const repository = typeorm.getRepository(ChannelStream);

    for (const channel of channelsIds) {
      const stream = streams.find((s) => s.userId === channel);

      const storedStream = await repository.findOneBy({
        userId: channel,
      });

      if (stream) {
        if (!storedStream) {
          await repository.save(convertSnakeToCamel(getRawData(stream)));
          await this.nats
            .emit('streams.online', { streamId: stream.id, channelId: channel })
            .toPromise();
        } else {
          await repository.update({ id: storedStream.id }, convertSnakeToCamel(getRawData(stream)));
        }
      } else {
        if (storedStream) {
          await repository.delete({ id: storedStream.id });
          await this.nats.emit('streams.offline', { channelId: channel }).toPromise();
        }
      }
    }

    return true;
  }
}

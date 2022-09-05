import { Injectable, OnModuleInit } from '@nestjs/common';
import { Client, Transport } from '@nestjs/microservices';
import { config } from '@tsuwari/config';
import { ClientProxy, ClientProxyEvents, RedisService } from '@tsuwari/shared';
import { ApiClient, HelixStreamData } from '@twurple/api';
import { ClientCredentialsAuthProvider } from '@twurple/auth';
import { getRawData } from '@twurple/common';
import _ from 'lodash';

const authProvider = new ClientCredentialsAuthProvider(
  config.TWITCH_CLIENTID,
  config.TWITCH_CLIENTSECRET,
);
const api = new ApiClient({ authProvider });

@Injectable()
export class AppService {
  @Client({ transport: Transport.NATS, options: { servers: [config.NATS_URL] } })
  nats: ClientProxy;

  constructor(private readonly redis: RedisService) {}

  async delStream(channelId: string) {
    await this.redis.del(`streams:${channelId}`);
  }

  async handleUpdate(e: ClientProxyEvents['stream.update']['input']) {
    const stream = await api.streams.getStreamByUserId(e.broadcaster_user_id);
    if (!stream) return;
    const cachedStream = await this.redis.get(`streams:${stream.userId}`);

    const rawData = getRawData(stream);
    rawData.title = e.title;
    rawData.game_name = e.category_name;
    rawData.game_id = e.category_id;

    await this.cacheStream(rawData, cachedStream ? JSON.parse(cachedStream) : null);
  }

  async handleOnline(e: ClientProxyEvents['streams.online']['input']) {
    const stream = await api.streams.getStreamByUserId(e.channelId);
    if (!stream) return;

    const cachedStream = await this.redis.get(`streams:${stream.userId}`);

    this.cacheStream(getRawData(stream), cachedStream ? JSON.parse(cachedStream) : null);
  }

  async handleOffline(e: ClientProxyEvents['streams.offline']['input']) {
    this.redis.del(`streams:${e.channelId}`);
  }

  async cacheStream(stream: HelixStreamData, cachedStream?: HelixStreamData | null) {
    let data = stream;

    if (cachedStream) {
      data = _.merge(cachedStream, data);
    }

    await this.redis.set(`streams:${stream.user_id}`, JSON.stringify(data), 'EX', 600);
  }

  async handleChannels(channelsIds: string[]) {
    const { data: streams } = await api.streams.getStreams({
      userId: channelsIds,
    });

    for (const channel of channelsIds) {
      const stream = streams.find((s) => s.userId === channel);

      const cachedStream = await this.redis.get(`streams:${channel}`);
      let parsedStream: HelixStreamData | null = null;

      if (cachedStream) {
        parsedStream = JSON.parse(cachedStream);
      }

      if (stream) {
        if (!parsedStream?.id) {
          await this.nats
            .emit('streams.online', { streamId: stream.id, channelId: channel })
            .toPromise();
        } else {
          await this.cacheStream(getRawData(stream), parsedStream);
        }
      } else {
        if (parsedStream?.id) {
          await this.nats.emit('streams.offline', { channelId: channel }).toPromise();
        }
      }
    }

    return true;
  }
}

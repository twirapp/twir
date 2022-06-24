import { Injectable } from '@nestjs/common';
import { Client, Transport } from '@nestjs/microservices';
import { config } from '@tsuwari/config';
import { ClientProxy, ClientProxyEvents } from '@tsuwari/shared';
import { ApiClient, HelixStream, HelixStreamData } from '@twurple/api';
import { ClientCredentialsAuthProvider } from '@twurple/auth';
import { getRawData } from '@twurple/common';
import _ from 'lodash';

import { RedisService } from './redis.service.js';

const authProvider = new ClientCredentialsAuthProvider(config.TWITCH_CLIENTID, config.TWITCH_CLIENTSECRET);
const api = new ApiClient({ authProvider });

@Injectable()
export class AppService {
  @Client({ transport: Transport.NATS, options: { servers: [config.NATS_URL] } })
  nats: ClientProxy;

  constructor(private readonly redis: RedisService) { }

  async delStream(channelId: string) {
    const key = `streams:${channelId}`;

    await this.redis.del(key);
  }

  async handleUpdate(e: ClientProxyEvents['stream.update']['input']) {
    const stream = await api.streams.getStreamByUserId(e.broadcaster_user_id);
    if (!stream) return;
    const cachedStream = await this.redis.get(`streams:${e.broadcaster_user_id}`);
    await this.cacheStream(stream, cachedStream);
  }

  async handleOnline(e: ClientProxyEvents['streams.online']['input']) {
    const stream = await api.streams.getStreamByUserId(e.channelId);
    if (!stream) return;
    const key = `streams:${e.channelId}`;
    const cachedStream = await this.redis.get(key);

    this.cacheStream(stream, cachedStream);
  }

  async handleOffline(e: ClientProxyEvents['streams.offline']['input']) {
    this.redis.del(`stream:${e.channelId}`);
  }

  async cacheStream(stream: HelixStream, cachedStream?: string | null) {
    const key = `streams:${stream.userId}`;
    let data = { ...getRawData(stream) };

    if (cachedStream) {
      data = _.merge(JSON.parse(cachedStream), data);
    }

    this.redis.set(
      key,
      JSON.stringify(data),
    ).then(() => {
      this.redis.expire(key, 600);
    });
  }

  async handleChannels(channelsIds: string[]) {
    const { data: streams } = await api.streams.getStreams({
      userId: channelsIds,
    });

    for (const channel of channelsIds) {
      const stream = streams.find(s => s.userId === channel);
      const key = `streams:${channel}`;
      const cachedStream = await this.redis.get(key);

      if (stream) {
        this.handleOnline({ streamId: stream.id, channelId: stream.userId });

        if (!cachedStream) {
          await this.nats.emit('streams.online', { streamId: stream.id, channelId: channel }).toPromise();
        }
      } else {
        if (cachedStream) {
          await this.nats.emit('streams.offline', { channelId: channel }).toPromise();
        }

        this.handleOffline({ channelId: channel });
      }
    }

    return true;
  }
}
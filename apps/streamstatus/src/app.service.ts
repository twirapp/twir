import { Injectable } from '@nestjs/common';
import { Client, ClientProxy, Transport } from '@nestjs/microservices';
import { config } from '@tsuwari/config';
import { ApiClient } from '@twurple/api';
import { ClientCredentialsAuthProvider } from '@twurple/auth';
import { getRawData } from '@twurple/common';

import { RedisService } from './redis.service.js';

const authProvider = new ClientCredentialsAuthProvider(config.TWITCH_CLIENTID, config.TWITCH_CLIENTSECRET);
const api = new ApiClient({ authProvider });

@Injectable()
export class AppService {
  @Client({ transport: Transport.NATS, options: { servers: [config.NATS_URL] } })
  nats: ClientProxy;

  constructor(private readonly redis: RedisService) {

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
        let data = { ...getRawData(stream) };
        if (cachedStream) {
          data = { ...data, ...JSON.parse(cachedStream) };
        } else {
          await this.nats.emit('streams.online', { streamId: stream.id, channelId: channel }).toPromise();
        }

        this.redis.set(
          key,
          JSON.stringify(data),
        ).then(() => {
          this.redis.expire(key, 600);
        });

      } else {
        if (cachedStream) {
          await this.nats.emit('streams.offline', { channelId: channel }).toPromise();
        }
        this.redis.del(key);
      }
    }
  }
}
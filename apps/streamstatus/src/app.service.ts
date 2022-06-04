import { Injectable } from '@nestjs/common';
import { config } from '@tsuwari/config';
import { ApiClient } from '@twurple/api';
import { ClientCredentialsAuthProvider } from '@twurple/auth';
import { getRawData } from '@twurple/common';

import { RedisService } from './redis.service.js';

const authProvider = new ClientCredentialsAuthProvider(config.TWITCH_CLIENTID, config.TWITCH_CLIENTSECRET);
const api = new ApiClient({ authProvider });

@Injectable()
export class AppService {
  constructor(private readonly redis: RedisService) {

  }


  async handleChannels(channelsIds: string[]) {
    const { data: streams } = await api.streams.getStreams({
      userId: channelsIds,
    });

    for (const channel of channelsIds) {
      const stream = streams.find(s => s.userId === channel);
      const key = `streams:${channel}`;

      if (stream) {
        this.redis.get(key).then(cachedStream => {
          let data = { ...getRawData(stream) };
          if (cachedStream) {
            data = { ...data, ...JSON.parse(cachedStream) };
          }

          this.redis.set(
            key,
            JSON.stringify(data),
          ).then(() => {
            this.redis.expire(key, 600);
          });
        });
      } else {
        this.redis.del(key);
      }
    }
  }
}
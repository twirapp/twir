import { Injectable, OnModuleInit } from '@nestjs/common';
import { Client, Transport } from '@nestjs/microservices';
import { config } from '@tsuwari/config';
import { RedisORMService, streamSchema, Stream, Repository } from '@tsuwari/redis';
import { ClientProxy, ClientProxyEvents } from '@tsuwari/shared';
import { ApiClient, HelixStreamData } from '@twurple/api';
import { ClientCredentialsAuthProvider } from '@twurple/auth';
import { getRawData } from '@twurple/common';
import _ from 'lodash';

const authProvider = new ClientCredentialsAuthProvider(config.TWITCH_CLIENTID, config.TWITCH_CLIENTSECRET);
const api = new ApiClient({ authProvider });

@Injectable()
export class AppService implements OnModuleInit {
  @Client({ transport: Transport.NATS, options: { servers: [config.NATS_URL] } })
  nats: ClientProxy;
  #redisRepository: Repository<Stream>;

  constructor(private readonly redisService: RedisORMService) { }

  onModuleInit() {
    this.#redisRepository = this.redisService.fetchRepository(streamSchema);
  }

  async delStream(channelId: string) {
    await this.redisService.fetchRepository(streamSchema).remove(channelId);
  }

  async handleUpdate(e: ClientProxyEvents['stream.update']['input']) {
    const stream = await api.streams.getStreamByUserId(e.broadcaster_user_id);
    if (!stream) return;
    const cachedStream = await this.#redisRepository.fetch(stream.userId);

    const rawData = getRawData(stream);
    rawData.title = e.title;
    rawData.game_name = e.category_name;
    rawData.game_id = e.category_id;

    await this.cacheStream(rawData, cachedStream);
  }

  async handleOnline(e: ClientProxyEvents['streams.online']['input']) {
    const stream = await api.streams.getStreamByUserId(e.channelId);
    if (!stream) return;

    const cachedStream = await this.#redisRepository.fetch(stream.userId);

    this.cacheStream(getRawData(stream), cachedStream);
  }

  async handleOffline(e: ClientProxyEvents['streams.offline']['input']) {
    this.#redisRepository.remove(e.channelId);
  }

  async cacheStream(stream: HelixStreamData, cachedStream?: Stream | null) {
    let data = stream;

    if (cachedStream) {
      data = _.merge(cachedStream.toRedisJson(), data);
    }

    await this.#redisRepository.createAndSave(data, stream.user_id);
    await this.#redisRepository.expire(stream.user_id, 600);
  }

  async handleChannels(channelsIds: string[]) {
    const { data: streams } = await api.streams.getStreams({
      userId: channelsIds,
    });

    for (const channel of channelsIds) {
      const stream = streams.find(s => s.userId === channel);

      const cachedStream = await this.#redisRepository.fetch(channel);

      if (stream) {
        if (!cachedStream.toRedisJson()?.id) {
          await this.nats.emit('streams.online', { streamId: stream.id, channelId: channel }).toPromise();
        } else {
          await this.cacheStream(getRawData(stream), cachedStream);
        }
      } else {
        if (cachedStream.toRedisJson()?.id) {
          await this.nats.emit('streams.offline', { channelId: channel }).toPromise();
        }
      }
    }

    return true;
  }
}
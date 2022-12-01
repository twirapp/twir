import { Injectable } from '@nestjs/common';
import { Interval } from '@nestjs/schedule';
import { config } from '@tsuwari/config';
import { convertSnakeToCamel } from '@tsuwari/shared';
import { Channel } from '@tsuwari/typeorm/entities/Channel';
import { ChannelStream } from '@tsuwari/typeorm/entities/ChannelStream';
import { ApiClient } from '@twurple/api';
import { ClientCredentialsAuthProvider } from '@twurple/auth';
import { getRawData } from '@twurple/common';
import _ from 'lodash';

import { pubSub } from './pubsub.js';
import { typeorm } from './typeorm.js';

const authProvider = new ClientCredentialsAuthProvider(
  config.TWITCH_CLIENTID,
  config.TWITCH_CLIENTSECRET,
);
const api = new ApiClient({ authProvider });

@Injectable()
export class AppService {
  async delStream(userId: string) {
    await typeorm.getRepository(ChannelStream).delete({
      userId,
    });
  }

  async handleUpdate(e: Record<string, any>) {
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

  async handleOnline(e: Record<string, any>) {
    const stream = await api.streams.getStreamByUserId(e.channelId);
    if (!stream) return;
    const repository = typeorm.getRepository(ChannelStream);
    await repository.delete({
      userId: e.channelId,
    });
    await repository.save(convertSnakeToCamel(getRawData(stream)));
  }

  async handleOffline(e: Record<string, any>) {
    await typeorm.getRepository(ChannelStream).delete({
      userId: e.channelId,
    });
  }

  @Interval(config.isDev ? 5000 : 5 * 60 * 1000)
  async handleChannels() {
    const channels = await typeorm.getRepository(Channel).find({
      where: { isEnabled: true },
      select: { id: true },
    });

    const chunks = _.chunk(
      channels.map((c) => c.id),
      100,
    );

    for (const chunk of chunks) {
      const { data: streams } = await api.streams.getStreams({
        userId: chunk,
      });

      const repository = typeorm.getRepository(ChannelStream);

      for (const channel of chunk) {
        const stream = streams.find((s) => s.userId === channel);

        const storedStream = await repository.findOneBy({
          userId: channel,
        });

        if (stream) {
          if (!storedStream) {
            await repository.save(convertSnakeToCamel(getRawData(stream)));
            pubSub.publish('stream.online', { streamId: stream.id, channelId: channel });
          } else {
            await repository.update(
              { id: storedStream.id },
              convertSnakeToCamel(getRawData(stream)),
            );
          }
        } else {
          if (storedStream) {
            await repository.delete({ id: storedStream.id });
            pubSub.publish('streams.offline', { channelId: channel });
          }
        }
      }
    }
  }
}

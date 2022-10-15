import { Injectable } from '@nestjs/common';
import { ModerationSettingsRepository, RedisDataSourceService } from '@tsuwari/redis';
import { ModerationSettingsDto } from '@tsuwari/shared';
import {
  ChannelModerationSetting,
  SettingsType,
} from '@tsuwari/typeorm/entities/ChannelModerationSetting';

import { typeorm } from '../../index.js';

@Injectable()
export class ModerationService {
  constructor(private readonly redisSource: RedisDataSourceService) {}

  async getSettings(channelId: string) {
    const keys = Object.values(SettingsType);
    const repository = typeorm.getRepository(ChannelModerationSetting);
    for (const key of keys) {
      await repository.upsert(
        {
          channelId,
          type: key,
        },
        {
          skipUpdateIfNoValuesChanged: true,
          conflictPaths: ['channelId', 'type'],
        },
      );
    }

    return repository.findBy({ channelId });
  }

  async update(channelId: string, data: ModerationSettingsDto[]) {
    const repository = typeorm.getRepository(ChannelModerationSetting);
    await Promise.all(
      data.map(async (item) => {
        const updateObject = {
          ...item,
          blackListSentences: item.blackListSentences as string[] | undefined,
        };

        const settings = await repository.findOneBy({
          channelId,
          type: item.type,
        });

        if (settings) {
          await repository.update({ id: settings.id }, updateObject);
        } else {
          await repository.save({
            channelId,
            ...updateObject,
          });
        }
      }),
    );

    const redisRepository = this.redisSource.getRepository(ModerationSettingsRepository);
    await Promise.all(data.map((item) => redisRepository.remove(`${channelId}:${item.type}`)));
    return repository.findBy({ channelId });
  }
}

import { Injectable } from '@nestjs/common';
import { PrismaService, SettingsType } from '@tsuwari/prisma';
import { ModerationUpdateDto } from '@tsuwari/shared';

import { RedisService } from '../../redis.service.js';

@Injectable()
export class ModerationService {
  constructor(
    private readonly redis: RedisService,
    private readonly prisma: PrismaService,
  ) { }


  async getSettings(channelId: string) {
    const keys = Object.values(SettingsType);
    const result = await Promise.all(keys.map(key => {
      return this.prisma.moderationSettings.upsert({
        where: {
          channelId_type: {
            channelId,
            type: key,
          },
        },
        update: {},
        create: {
          type: key,
          channelId,
        },
      });
    }));

    return result;
  }

  async update(channelId: string, data: ModerationUpdateDto) {
    const result = await Promise.all(data.items.map(item => {
      const updateObject = {
        ...item,
        blackListSentences: item.blackListSentences as string[] | undefined,
      };

      return this.prisma.moderationSettings.upsert({
        where: {
          channelId_type: {
            channelId,
            type: item.type,
          },
        },
        update: updateObject,
        create: {
          ...updateObject,
          channelId,
        },
      });
    }));

    return result;
  }
}

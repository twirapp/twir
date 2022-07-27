import { Injectable } from '@nestjs/common';
import { PrismaService, SettingsType } from '@tsuwari/prisma';
import { ModerationSettings, moderationSettingsSchema, RedisORMService, Repository } from '@tsuwari/redis';
import { ModerationSettingsDto, RedisService } from '@tsuwari/shared';


@Injectable()
export class ModerationService {
  #repository: Repository<ModerationSettings>;

  constructor(
    private readonly prisma: PrismaService,
    private readonly redis: RedisService,
    private readonly redisOrm: RedisORMService,
  ) { }

  onModuleInit() {
    this.#repository = this.redisOrm.fetchRepository(moderationSettingsSchema);
  }


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

  async update(channelId: string, data: ModerationSettingsDto[]) {
    const result = await Promise.all(data.map(item => {
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
          channelId: undefined,
          channel: {
            connect: {
              id: channelId,
            },
          },
        },
      });
    }));

    await Promise.all(data.map(item => this.#repository.remove(`${channelId}:${item.type}`)));
    return result;
  }
}

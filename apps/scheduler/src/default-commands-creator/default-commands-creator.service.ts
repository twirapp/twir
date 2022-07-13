import { Inject, Injectable, Logger } from '@nestjs/common';
import { Interval } from '@nestjs/schedule';
import { config } from '@tsuwari/config';
import { PrismaService, Channel } from '@tsuwari/prisma';
import { ClientProxy, RedisService } from '@tsuwari/shared';
import { lastValueFrom } from 'rxjs';

@Injectable()
export class DefaultCommandsCreatorService {
  #logger = new Logger(DefaultCommandsCreatorService.name);

  constructor(
    private readonly prisma: PrismaService,
    @Inject('NATS') private nats: ClientProxy,
    private readonly redis: RedisService,
  ) { }

  @Interval('defaultCommands', config.isDev ? 1000 : 1 * 60 * 1000)
  async createDefaultCommands(usersIds?: string[]) {
    const defaultCommands = await lastValueFrom(this.nats.send('bots.getDefaultCommands', {}));

    for (const command of defaultCommands) {
      // const usersForUpdate: Channel[] = await this.prisma.$queryRaw`SELECT * FROM "channels" where id not in (select "channelId" from "channels_commands" where "channels_commands"."defaultName" = ${command.name})`;
      const usersForUpdate = usersIds
        ? await this.prisma.channel.findMany({
          where: {
            id: {
              in: usersIds,
            },
            commands: {
              none: {
                default: true,
                defaultName: command.name,
              },
            },
          },
          select: {
            id: true,
          },
        })
        : await this.prisma.channel.findMany({
          where: {
            commands: {
              none: {
                default: true,
                defaultName: command.name,
              },
            },
          },
          select: {
            id: true,
          },
        });

      if (!usersForUpdate.length) continue;

      this.#logger.log(`Creating default command ${command.name} for ${usersForUpdate.length} users.`);

      for (const channel of usersForUpdate) {
        const newCommand = await this.prisma.command.create({
          data: {
            channelId: channel.id,
            default: true,
            defaultName: command.name,
            description: command.description,
            visible: command.visible,
            name: command.name,
            permission: command.permission,
            cooldown: 0,
            cooldownType: 'GLOBAL',
          },
        });

        const commandForSet = {
          ...newCommand,
          responses: JSON.stringify([]),
          aliases: JSON.stringify([]),
        };

        await this.redis.hmset(`commands:${channel.id}:${command.name}`, commandForSet);
      }
    }
  }
}

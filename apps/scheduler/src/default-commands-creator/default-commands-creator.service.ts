import { Inject, Injectable, Logger } from '@nestjs/common';
import { Interval } from '@nestjs/schedule';
import { config } from '@tsuwari/config';
import { PrismaService } from '@tsuwari/prisma';
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

  @Interval(config.isDev ? 1000 : 5 * 60 * 1000)
  async createDefaultCommands() {
    const defaultCommands = await lastValueFrom(this.nats.send('bots.getDefaultCommands', {}));
    const defaultCommandsNames = defaultCommands.map(c => c.name);

    const usersForUpdate = await this.prisma.channel.findMany({
      where: {
        commands: {
          none: {
            defaultName: {
              in: defaultCommandsNames,
            },
            default: true,
          },
        },
      },
      include: {
        commands: true,
      },
    });

    this.#logger.log(`Creating default commands for ${usersForUpdate.length} users.`);

    for (const user of usersForUpdate) {
      const existedCommands = user.commands.filter(c => c.default).map(c => c.defaultName);
      const commandsForCreate = defaultCommands.filter(c => !existedCommands.includes(c.name));

      for (const command of commandsForCreate) {
        const newCommand = await this.prisma.command.create({
          data: {
            channelId: user.id,
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

        await this.redis.hmset(`commands:${user.id}:${command.name}`, commandForSet);
      }
    }
  }
}

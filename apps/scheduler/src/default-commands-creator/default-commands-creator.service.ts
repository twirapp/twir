import { Inject, Injectable, Logger, OnModuleInit } from '@nestjs/common';
import { Interval } from '@nestjs/schedule';
import { config } from '@tsuwari/config';
import { PrismaService } from '@tsuwari/prisma';
import { ClientProxy, RedisService } from '@tsuwari/shared';
import * as Knex from 'knex';
import { lastValueFrom } from 'rxjs';

@Injectable()
export class DefaultCommandsCreatorService implements OnModuleInit {
  #logger = new Logger(DefaultCommandsCreatorService.name);
  #knex: Knex.Knex;

  constructor(
    private readonly prisma: PrismaService,
    @Inject('NATS') private nats: ClientProxy,
    private readonly redis: RedisService,
  ) { }

  async onModuleInit() {
    const knex = Knex.default({
      client: 'pg',
      connection: config.DATABASE_URL,
    });

    this.#knex = knex;
  }

  @Interval('defaultCommands', config.isDev ? 1000 : 1 * 60 * 1000)
  async createDefaultCommands(usersIds?: string[]) {
    const defaultCommands = await lastValueFrom(this.nats.send('bots.getDefaultCommands', {}));
    const defaultCommandsNames = defaultCommands.map(c => c.name);

    const channels: Array<{
      id: string,
      commands: string[],
    }> = await this.#knex
      .select('channels.id', this.#knex.raw('array_remove(array_agg("channels_commands"."defaultName"),null) as commands'))
      .from('channels')
      .leftJoin('channels_commands', function () {
        this
          .on('channels_commands.channelId', '=', 'channels.id')
          .andOnIn('channels_commands.defaultName', defaultCommandsNames);
      })
      .groupBy('channels.id')
      .modify(function (b) {
        if (usersIds) {
          b.whereIn('channels.id', usersIds);
        }
      })
      .having(this.#knex.raw(`count("defaultName") < ${defaultCommandsNames.length}`));


    if (channels.length) {
      this.#logger.log(`Creating default commands for ${channels.length} channels.`);
    }

    for (const channel of channels) {
      const commandsForCreate = defaultCommands.filter(c => !channel.commands.includes(c.name));

      for (const command of commandsForCreate) {
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
            module: command.module,
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

    /* for (const command of defaultCommands) {
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
    } */
  }
}

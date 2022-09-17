import { Injectable, Logger, OnModuleInit } from '@nestjs/common';
import { Interval } from '@nestjs/schedule';
import { config } from '@tsuwari/config';
import * as Parser from '@tsuwari/nats/parser';
import { RedisService } from '@tsuwari/shared';
import {
  ChannelCommand,
  CommandModule,
  CommandPermission,
  CooldownType,
} from '@tsuwari/typeorm/entities/ChannelCommand';
import * as Knex from 'knex';

import { typeorm } from '../index.js';
import { nats } from '../libs/nats.js';

@Injectable()
export class DefaultCommandsCreatorService implements OnModuleInit {
  #logger = new Logger(DefaultCommandsCreatorService.name);
  #knex: Knex.Knex;

  constructor(private readonly redis: RedisService) {}

  async onModuleInit() {
    const knex = Knex.default({
      client: 'pg',
      connection: config.DATABASE_URL,
    });

    this.#knex = knex;
  }

  @Interval('defaultCommands', config.isDev ? 1000 : 1 * 60 * 1000)
  async createDefaultCommands(usersIds?: string[]) {
    const msg = await nats.request('bots.getDefaultCommands', new Uint8Array());
    const { list: defaultCommands } = Parser.GetDefaultCommandsResponse.fromBinary(msg.data);

    const defaultCommandsNames = defaultCommands.map((c) => c.name);

    const channels: Array<{
      id: string;
      commands: string[];
    }> = await this.#knex
      .select(
        'channels.id',
        this.#knex.raw(
          'array_remove(array_agg("channels_commands"."defaultName"),null) as commands',
        ),
      )
      .from('channels')
      .leftJoin('channels_commands', function () {
        this.on('channels_commands.channelId', '=', 'channels.id').andOnIn(
          'channels_commands.defaultName',
          defaultCommandsNames,
        );
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

    const repository = typeorm.getRepository(ChannelCommand);
    for (const channel of channels) {
      const commandsForCreate = defaultCommands.filter((c) => !channel.commands.includes(c.name));

      for (const command of commandsForCreate) {
        const newCommand = await repository.save({
          channelId: channel.id,
          default: true,
          defaultName: command.name,
          description: command.description,
          visible: command.visible,
          name: command.name,
          permission: command.permission as unknown as CommandPermission,
          cooldown: 0,
          cooldownType: CooldownType.GLOBAL,
          module: command.module as unknown as CommandModule | undefined,
        });

        const commandForSet = {
          ...newCommand,
          responses: [],
          aliases: [],
        };

        await this.redis.set(
          `commands:${channel.id}:${command.name}`,
          JSON.stringify(commandForSet),
        );
      }
    }
  }
}

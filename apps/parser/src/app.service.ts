import { Injectable, OnApplicationBootstrap, OnModuleInit } from '@nestjs/common';
import { PrismaService } from '@tsuwari/prisma';
import { RedisORMService } from '@tsuwari/redis';
import { ClientProxyCommands, RedisService, TwitchApiService } from '@tsuwari/shared';
import { TwitchPrivateMessage } from '@twurple/chat/lib/commands/TwitchPrivateMessage.js';

import { DefaultCommand } from './defaultCommands/types.js';
import { CommandConditional, HelpersService } from './helpers.service.js';
import { FaceitIntegration } from './integrations/faceit.js';
import { ParserCache } from './variables/cache.js';
import { State, VariablesParser } from './variables/index.js';

@Injectable()
export class AppService implements OnModuleInit {
  #defaultCommands: Map<string, DefaultCommand> = new Map();

  constructor(
    private readonly helpers: HelpersService,
    private readonly prisma: PrismaService,
    private readonly redis: RedisService,
    private readonly staticApi: TwitchApiService,
    private readonly faceitIntegration: FaceitIntegration,
    private readonly variablesParser: VariablesParser,
    private readonly redisOm: RedisORMService,
  ) { }

  onModuleInit() {
    setTimeout(async () => {
      const defCommands = await import('./defaultCommands/index.js');

      for (const command of Object.values(defCommands).flat()) {
        this.#defaultCommands.set(command.name, command);
      }
    }, 1000);
  }

  async getResponses(message: string, state: TwitchPrivateMessage) {
    if (!state.channelId) return;

    const channelCommands = await this.helpers.getChannelCommands(state.channelId);
    if (!channelCommands || !channelCommands.length) return;

    const findCommand = this.helpers.findCommandInArrayOfNames(
      message,
      channelCommands.map(c => [c.name, ...c.aliases]).flat(),
    );
    if (!findCommand.isFound) return;

    const command = channelCommands.find(c => [c.name, ...c.aliases].includes(findCommand.commandName));
    if (!command || !command.enabled) return;

    const userPermissions = await this.helpers.getUserPermissions(state.userInfo, {
      checkAdmin: command.permission !== 'VIEWER',
      checkFollower: command.permission === 'FOLLOWER',
      channelId: state.channelId,
    });
    const hasPermission = this.helpers.hasPermission(userPermissions, command.permission);

    if (!hasPermission) {
      return;
    }

    if (command.default && command.defaultName) {
      const cmd = this.#defaultCommands.get(command.defaultName);
      if (cmd) {
        const result = await cmd.handler(state, findCommand.params);
        command.responses = result ? Array.isArray(result) ? result : [result] : [];
      }
    }

    // const lock = await redlock.acquire([`locks:commandsParser:msg:${state.id}`], 1000);
    try {
      const onCooldown = await this.helpers.isOnCooldown(command, state.userInfo.userId);
      if (onCooldown && !(userPermissions.BROADCASTER || userPermissions.MODERATOR)) return;

      if (!command.responses?.length) return;

      this.prisma.commandUsage.create({
        data: { commandId: command.id, channelId: state.channelId, userId: state.userInfo.userId },
      });
      this.helpers.setCommandCooldown(command, state.userInfo.userId);

      return {
        responses: command.responses,
        params: findCommand.params,
        commandName: findCommand.commandName,
      };
    } finally {
      // await lock.release();
    }
  }

  async parseResponses(state: ClientProxyCommands['parseResponse']['input'], response: {
    responses: (string | undefined)[];
    params?: string;
  }) {
    const msgObject = {
      channelId: state.channelId,
      sender: {
        id: state.userId,
        name: state.userName,
      },
      cache: new ParserCache(this.staticApi, this.prisma, this.redis, this.faceitIntegration, state.channelId, this.redisOm, state.userId),
    };

    const parsedResponses: string[] = [];

    for (const r of response.responses) {
      if (!r) continue;

      const parsed = await this.variablesParser.parse(r, msgObject, response.params);
      parsedResponses.push(parsed);
    }

    return parsedResponses;
  }
}

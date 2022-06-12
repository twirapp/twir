import { setTimeout as sleep } from 'node:timers/promises';

import { Command, CommandPermission, Response } from '@tsuwari/prisma';
import { ChatUser } from '@twurple/chat/lib';
import type { TwitchPrivateMessage } from '@twurple/chat/lib/commands/TwitchPrivateMessage';

import * as DefCommands from '../defaultCommands/index.js';
import { DefaultCommand } from '../defaultCommands/types.js';
import { getChannelCommandsNamesFromRedis } from '../functions/getChannelCommandListFromRedis.js';
import { prisma } from './prisma.js';
import { redis, redlock } from './redis.js';

const defaultCommands: Map<keyof typeof DefCommands | string, DefaultCommand> = new Map();

for (const command of Object.values(DefCommands)) {
  defaultCommands.set(command.name, command);
}

export type CommandConditional = Command & { responses: (string | undefined)[] | undefined };

export class CommandsParser {
  async parse(message: string, state: TwitchPrivateMessage) {
    if (!message.startsWith('!') || !state.channelId) return;
    const channelCommandsNames = await getChannelCommandsNamesFromRedis(state.channelId);
    if (!channelCommandsNames || !channelCommandsNames.length) return;

    const findCommand = this.#findCommandInArrayOfNames(message, channelCommandsNames as string[]);
    if (!findCommand.isFound) return;

    const command: CommandConditional = (await redis.hgetall(
      `commands:${state.channelId}:${findCommand.commandName}`,
    )) as unknown as CommandConditional;
    if (!command || !command.enabled) return;

    const userPermissions = this.#getUserPermissions(state.userInfo);
    const hasPermission = this.#hasPermission(userPermissions, command.permission);

    if (!hasPermission) {
      return;
    }

    if (command.default && command.defaultName) {
      const cmd = defaultCommands.get(command.defaultName);
      if (cmd) {
        const result = await cmd.handler(state, findCommand.params);
        command.responses = Array.isArray(result) ? result : [result];
      }
    } else {
      command.responses = JSON.parse(command.responses as unknown as string);
    }

    // const lock = await redlock.acquire([`locks:commandsParser:msg:${state.id}`], 1000);
    try {
      const onCooldown = await this.#isOnCooldown(command, state.userInfo.userId);
      if (onCooldown && !(userPermissions.BROADCASTER || userPermissions.MODERATOR)) return;

      if (!command.responses?.length) return;

      prisma.commandUsage.create({
        data: { commandId: command.id, channelId: state.channelId, userId: state.userInfo.userId },
      });
      this.#setCommandCooldown(command, state.userInfo.userId);

      return command.responses;
    } finally {
      // await lock.release();
    }
  }

  #findCommandInArrayOfNames(message: string, commands: string[]) {
    message = message.substring(1).trim();

    const msgArray = message.toLowerCase().split(' ');
    let isFound = false;
    let commandName = '';

    for (let i = 0, len = msgArray.length; i < len; i++) {
      const query = msgArray.join(' ');
      const find = commands.find((c) => c === query);
      if (!find) {
        msgArray.pop();
        continue;
      }

      commandName = find;
      isFound = true;
    }

    return {
      isFound,
      commandName,
      params: message.replace(new RegExp(`^${commandName}`), '').trim(),
    };
  }

  #getUserPermissions(userInfo: ChatUser): Record<CommandPermission, boolean> {
    return {
      BROADCASTER: userInfo.isBroadcaster,
      MODERATOR: userInfo.isMod,
      VIP: userInfo.isVip,
      SUBSCRIBER: userInfo.isSubscriber || userInfo.isFounder,
      FOLLOWER: true,
      VIEWER: true,
    };
  }

  #hasPermission(perms: Record<CommandPermission, boolean>, searchForPermission: CommandPermission) {
    if (!searchForPermission) return true;

    const userPerms = Object.entries(perms) as [CommandPermission, boolean][];
    const permissionIndex = userPerms.find((perm) => perm[0] === searchForPermission);
    const commandPermissionIndex = userPerms.indexOf(permissionIndex!);

    const hasPerm = userPerms.some((p, index) => p[1] && index <= commandPermissionIndex);
    return hasPerm;
  }

  #buildCooldownKey(command: CommandConditional, senderId: string) {
    if (command.cooldownType === 'GLOBAL') {
      return `commands:cooldowns:${command.id}`;
    } else {
      return `commands:cooldowns:${command.id}:${senderId}`;
    }
  }

  async #isOnCooldown(command: CommandConditional, senderId: string) {
    if (!command.cooldown) return false;
    const item = await redis.get(this.#buildCooldownKey(command, senderId));
    return item !== null;
  }

  #setCommandCooldown(command: CommandConditional, senderId: string) {
    if (command.cooldown && command.cooldown <= 0) return;

    if (command.cooldownType === 'GLOBAL') {
      redis.set(`commands:cooldowns:${command.id}`, '', 'EX', command.cooldown!);
    } else {
      redis.set(`commands:cooldowns:${command.id}:${senderId}`, '', 'EX', command.cooldown!);
    }
  }
}

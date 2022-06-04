import { setTimeout as sleep } from 'node:timers/promises';

import { Command, CommandPermission, Response } from '@tsuwari/prisma';
import { ChatUser } from '@twurple/chat/lib';
import type { TwitchPrivateMessage } from '@twurple/chat/lib/commands/TwitchPrivateMessage';

import { getChannelCommandsNamesFromRedis } from '../functions/getChannelCommandListFromRedis.js';
import { prisma } from './prisma.js';
import { redis, redlock } from './redis.js';


export type CommandConditional = Command & { responses: string[] };

export class CommandsParser {
  async #getCommandResponses(channelId: string, commandId: string) {
    const keys = await redis.keys(`commands:${channelId}:${commandId}:responses:*`);
    if (!keys.length) return;

    const responsesIds = keys.map((k) => k.split(':')[4]);
    if (!responsesIds?.length) return;

    const responses = await Promise.all(
      responsesIds.map((id) => redis.hgetall(`commands:${channelId}:${commandId}:responses:${id}`)),
    );

    return responses as Response[];
  }

  async parse(message: string, state: TwitchPrivateMessage) {
    if (!message.startsWith('!') || !state.channelId) return;
    const channelCommandsNames = await getChannelCommandsNamesFromRedis(state.channelId);
    if (!channelCommandsNames || !channelCommandsNames.length) return;

    const findCommand = this.#findCommandInArrayOfNames(message, channelCommandsNames as string[]);
    if (!findCommand.isFound) return;

    const command: CommandConditional = (await redis.hgetall(
      `commands:${state.channelId}:${findCommand.commandName}`,
    )) as unknown as CommandConditional;
    command.responses = JSON.parse(command.responses as unknown as string);

    if (!command || !command.enabled) return;

    const userPermissions = this.#getUserPermissions(state.userInfo);
    const hasPermission = this.#hasPermission(userPermissions, command.permission);

    if (!hasPermission) {
      return;
    }

    // const lock = await redlock.acquire([`locks:commandsParser:msg:${state.id}`], 1000);
    try {
      const onCooldown = await this.#isOnCooldown(command, state.userInfo.userId);
      if (onCooldown /* && !(userPermissions.BROADCASTER || userPermissions.MODERATOR) */) return;

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

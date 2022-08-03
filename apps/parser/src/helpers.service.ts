import { Injectable, OnModuleInit } from '@nestjs/common';
import { Command, CommandPermission, PrismaService } from '@tsuwari/prisma';
import { RedisORMService, Repository, Command as CommandCacheClass, commandSchema } from '@tsuwari/redis';
import { RedisService, TwitchApiService } from '@tsuwari/shared';
import { ChatUser } from '@twurple/chat';

export type CommandConditional = Command & { responses: (string | undefined)[] | undefined };

@Injectable()
export class HelpersService implements OnModuleInit {
  #commandsRepository: Repository<CommandCacheClass>;

  constructor(
    private readonly redis: RedisService,
    private readonly prisma: PrismaService,
    private readonly twitchApi: TwitchApiService,
    private readonly redisOrm: RedisORMService,
  ) { }

  onModuleInit() {
    this.#commandsRepository = this.redisOrm.fetchRepository(commandSchema);
  }

  async getChannelCommands(channelId: string) {
    const cmds = await this.#commandsRepository.search()
      .where('channelId').eq(channelId)
      .returnAll();

    return cmds.map(c => c.toRedisJson());
  }

  async getUserPermissions(userInfo: ChatUser,
    opts: {
      checkAdmin?: boolean,
      checkFollower?: boolean,
      channelId?: string
    } = {},
  ): Promise<Record<CommandPermission, boolean>> {
    const dbUser = opts.checkAdmin
      ? await this.prisma.user.findFirst({ where: { id: userInfo.userId } })
      : null;
    const twitchFollow = (opts.channelId && opts.checkFollower)
      ? await this.twitchApi.users.getFollowFromUserToBroadcaster(userInfo.userId, opts.channelId)
      : null;

    const perms = {
      BROADCASTER: userInfo.isBroadcaster,
      MODERATOR: userInfo.isMod,
      VIP: userInfo.isVip,
      SUBSCRIBER: userInfo.isSubscriber || userInfo.isFounder,
      FOLLOWER: !!twitchFollow,
      VIEWER: true,
    };

    if (dbUser?.isBotAdmin) {
      Object.keys(perms).forEach((key) => {
        perms[key as CommandPermission] = true;
      });
    }

    return perms;
  }

  hasPermission(perms: Record<CommandPermission, boolean>, searchForPermission: CommandPermission) {
    if (!searchForPermission) return true;

    const userPerms = Object.entries(perms) as [CommandPermission, boolean][];
    const permissionIndex = userPerms.find((perm) => perm[0] === searchForPermission);
    const commandPermissionIndex = userPerms.indexOf(permissionIndex!);

    const hasPerm = userPerms.some((p, index) => p[1] && index <= commandPermissionIndex);
    return hasPerm;
  }

  buildCooldownKey(command: CommandConditional, senderId: string) {
    if (command.cooldownType === 'GLOBAL') {
      return `commands:cooldowns:${command.id}`;
    } else {
      return `commands:cooldowns:${command.id}:${senderId}`;
    }
  }

  async isOnCooldown(command: CommandConditional, senderId: string) {
    if (!command.cooldown) return false;
    const item = await this.redis.get(this.buildCooldownKey(command, senderId));
    return item !== null;
  }

  setCommandCooldown(command: CommandConditional, senderId: string) {
    if (!command.cooldown || command.cooldown === 0) return;

    if (command.cooldownType === 'GLOBAL') {
      this.redis.set(`commands:cooldowns:${command.id}`, '', 'EX', command.cooldown);
    } else {
      this.redis.set(`commands:cooldowns:${command.id}:${senderId}`, '', 'EX', command.cooldown);
    }
  }

  findCommandInArrayOfNames(message: string, commands: string[]) {
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
}

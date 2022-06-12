import { HttpException, Inject, Injectable, OnModuleInit } from '@nestjs/common';
import { ClientGrpc, ClientProxy } from '@nestjs/microservices';
import { Bots } from '@tsuwari/grpc';
import { Command, PrismaService, Response } from '@tsuwari/prisma';

import { RedisService } from '../../redis.service.js';
import { UpdateOrCreateCommandDto } from './dto/create.js';

@Injectable()
export class CommandsService implements OnModuleInit {
  botsMicroservce: Bots.Commands;

  constructor(
    private readonly prisma: PrismaService,
    private readonly redis: RedisService,
    @Inject('BOTS_MICROSERVICE') private readonly botsClient: ClientGrpc,
  ) { }

  onModuleInit() {
    this.botsMicroservce = this.botsClient.getService<Bots.Commands>('Commands');
  }

  async getList(userId: string) {
    const commands: (Command & {
      responses?: Response[];
    })[] = await this.prisma.command.findMany({
      where: { channelId: userId },
      include: {
        responses: true,
      },
    });

    const defaultCommands = await this.botsMicroservce.getDefaultCommands({}).toPromise();
    if (defaultCommands?.commands) {
      for (const command of defaultCommands.commands) {
        if (!commands.some(c => c.defaultName === command.name)) {
          const newCommand = await this.prisma.command.create({
            data: {
              channelId: userId,
              default: true,
              defaultName: command.name!,
              description: command.description,
              name: command.name!,
              permission: command.permission! as any,
              cooldown: 0,
              cooldownType: 'GLOBAL',
            },
          });

          this.#setCommandCache(newCommand);
          commands.push(newCommand);
        }
      }
    }

    return commands;
  }

  async #setCommandCache(command: Command & { responses?: Response[] }, oldCommand?: Command & { responses?: Response[] }) {
    const preKey = `commands:${command.channelId}`;

    if (oldCommand) {
      await this.redis.del(`commands:${oldCommand.channelId}:${oldCommand.name}`);

      if (oldCommand.aliases && Array.isArray(command.aliases)) {
        for (const alias of command.aliases) {
          await this.redis.del(`${preKey}:${alias}`);
        }
      }
    }

    const commandForSet = {
      ...command,
      responses: command.responses ? JSON.stringify(command.responses.map(r => r.text)) : [],
      aliases: Array.isArray(command.aliases) ? JSON.stringify(command.aliases) : command.aliases,
    };

    await this.redis.hmset(`${preKey}:${command.name}`, commandForSet);

    if (command.aliases && Array.isArray(command.aliases)) {
      for (const alias of command.aliases) {
        await this.redis.hmset(`${preKey}:${alias}`, commandForSet);
      }
    }

  }

  async create(userId: string, data: UpdateOrCreateCommandDto & { defaultName?: string }) {
    const isExists = await this.prisma.command.findMany({
      where: {
        name: data.name,
        OR: {
          name: { in: data.aliases },
          aliases: {
            array_contains: data.aliases,
          },
        },
      },
    });

    if (isExists.length) {
      throw new HttpException(`Command already exists`, 400);
    }

    const command = await this.prisma.command.create({
      data: {
        ...data,
        channelId: userId,
        responses: {
          createMany: {
            data: data.responses.map((r) => ({ text: r.text?.trim().replace(/(\r\n|\n|\r)/, '') })),
          },
        },
      },
      include: {
        responses: true,
      },
    });

    await this.#setCommandCache(command);
    return command;
  }

  async delete(userId: string, commandId: string) {
    const command = await this.prisma.command.findFirst({ where: { channelId: userId, id: commandId } });

    if (!command) {
      throw new Error('Command not exists');
    }

    if (command.default) {
      throw new Error('You cannot delete default command.');
    }

    const result = await this.prisma.command.delete({
      where: {
        id: commandId,
      },
    });

    await this.redis.del(`commands:${userId}:${command.name}`);
    if (Array.isArray(command.aliases)) {
      for (const aliase of command.aliases as string[]) {
        await this.redis.del(`commands:${userId}:${aliase}`);
      }
    }

    return result;
  }

  async update(userId: string, commandId: string, data: UpdateOrCreateCommandDto) {
    const command = await this.prisma.command.findFirst({
      where: { channelId: userId, id: commandId },
      include: { responses: true },
    });

    if (!command) {
      throw new Error('Command not exists');
    }

    data.responses = data.responses?.map(r => ({ ...r, text: r.text ? r.text.trim().replace(/(\r\n|\n|\r)/, '') : null }));

    const responsesForUpdate = data.responses
      .filter(r => command.responses.some(c => c.id === r.id && r.text && r.id))
      .map(r => ({ id: r.id, text: r.text }))
      .map(r => this.prisma.response.update({ where: { id: r.id }, data: { text: r.text } }));

    const [newCommand] = await this.prisma.$transaction([
      this.prisma.command.update({
        where: { id: commandId },
        data: {
          ...data,
          channelId: userId,
          responses: {
            deleteMany: command.responses.filter(r => !data.responses.map(s => s.id).includes(r.id)),
            createMany: {
              data: data.responses
                .filter((r) => !command.responses.some((c) => c.id === r.id)),
              skipDuplicates: true,
            },
          },
        },
        include: {
          responses: true,
        },
      }),
      ...responsesForUpdate,
    ]);

    const newResponses = await this.prisma.response.findMany({ where: { commandId: command.id } });

    await this.#setCommandCache({
      ...newCommand,
      responses: newResponses.flat(),
    }, command);

    return {
      ...newCommand,
      responses: newResponses.flat(),
    };
  }
}

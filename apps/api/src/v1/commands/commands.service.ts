import { HttpException, Injectable, OnModuleInit } from '@nestjs/common';
import { Client, Transport } from '@nestjs/microservices';
import { config } from '@tsuwari/config';
import { Command, PrismaService, Response } from '@tsuwari/prisma';
import { RedisORMService, Repository, Command as CommandCacheClass, commandSchema } from '@tsuwari/redis';
import { ClientProxy, RedisService } from '@tsuwari/shared';

import { UpdateOrCreateCommandDto } from './dto/create.js';

@Injectable()
export class CommandsService implements OnModuleInit {
  @Client({ transport: Transport.NATS, options: { servers: [config.NATS_URL] } })
  nats: ClientProxy;
  #commandsRepository: Repository<CommandCacheClass>;

  constructor(
    private readonly prisma: PrismaService,
    private readonly redis: RedisService,
    private readonly redisOrm: RedisORMService,
  ) { }

  async onModuleInit() {
    await this.redisOrm.open(config.REDIS_URL);
    this.#commandsRepository = this.redisOrm.fetchRepository(commandSchema);
  }

  async getList(userId: string) {
    await this.nats.send('bots.createDefaultCommands', [userId]).toPromise();

    const commands: (Command & {
      responses?: Response[];
    })[] = await this.prisma.command.findMany({
      where: { channelId: userId },
      include: {
        responses: true,
      },
    });

    return commands;
  }

  async setCommandCache(command: Command & { responses?: Response[] }, oldCommand?: Command & { responses?: Response[] }) {

    if (oldCommand) {
      await this.#commandsRepository.remove(`${oldCommand.channelId}:${oldCommand.name}`);

      if (oldCommand.aliases && Array.isArray(command.aliases)) {
        for (const alias of command.aliases) {
          await this.#commandsRepository.remove(`${oldCommand.channelId}:${alias}`);
        }
      }
    }

    const commandForSet = {
      ...command,
      responses: command.responses ? command.responses.filter(r => r.text).map(r => r.text!) : [] as string[],
      aliases: command.aliases as string[],
      defaultName: command.defaultName ?? null,
    };

    await this.#commandsRepository.createAndSave(commandForSet, `${command.channelId}:${command.name}`);

    /* if (command.aliases && Array.isArray(command.aliases)) {
      for (const alias of command.aliases) {
        await this.#commandsRepository.createAndSave(commandForSet, `${command.channelId}:${alias}`);
      }
    } */

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

    if (!data.responses?.length) {
      throw new HttpException(`You should add atleast 1 response to command.`, 400);
    }

    const command = await this.prisma.command.create({
      data: {
        ...data,
        channelId: userId,
        responses: {
          createMany: {
            data: data.responses.filter(r => r.text).map((r) => ({ text: r.text?.trim().replace(/(\r\n|\n|\r)/, '') })),
          },
        },
      },
      include: {
        responses: true,
      },
    });

    await this.setCommandCache(command);
    return command;
  }

  async delete(userId: string, commandId: string) {
    const command = await this.prisma.command.findFirst({ where: { channelId: userId, id: commandId } });

    if (!command) {
      throw new HttpException('Command not exists', 404);
    }

    if (command.default) {
      throw new HttpException('You cannot delete default command.', 400);
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
      throw new HttpException('Command not exists', 404);
    }

    if (!data.responses?.length && !command.default) {
      throw new HttpException(`You should add atleast 1 response to command.`, 400);
    }

    data.responses = data.responses?.filter(r => r.text).map(r => ({ ...r, text: r.text ? r.text.trim().replace(/(\r\n|\n|\r)/, '') : null }));

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

    await this.setCommandCache({
      ...newCommand,
      responses: newResponses.flat(),
    }, command);

    return {
      ...newCommand,
      responses: newResponses.flat(),
    };
  }
}

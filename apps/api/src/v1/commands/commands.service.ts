import { HttpException, Injectable } from '@nestjs/common';
import { Client, Transport } from '@nestjs/microservices';
import { config } from '@tsuwari/config';
import { ClientProxy, RedisService } from '@tsuwari/shared';
import { ChannelCommand } from '@tsuwari/typeorm/entities/ChannelCommand';
import { CommandResponse } from '@tsuwari/typeorm/entities/CommandResponse';

import { typeorm } from '../../index.js';
import { UpdateOrCreateCommandDto } from './dto/create.js';

@Injectable()
export class CommandsService {
  @Client({ transport: Transport.NATS, options: { servers: [config.NATS_URL] } })
  nats: ClientProxy;

  constructor(private readonly redis: RedisService) {}

  async getList(userId: string) {
    await this.nats.send('bots.createDefaultCommands', [userId]).toPromise();

    const commands = await typeorm.getRepository(ChannelCommand).find({
      where: { channelId: userId },
      relations: {
        responses: true,
      },
    });

    return commands;
  }

  async setCommandCache(command: ChannelCommand, oldCommand?: ChannelCommand) {
    if (oldCommand) {
      await this.redis.del(`commands:${oldCommand.channelId}:${oldCommand.name}`);
    }

    const commandForSet = {
      ...command,
      responses: command.responses
        ? command.responses.filter((r) => r.text).map((r) => r.text!)
        : ([] as string[]),
      aliases: command.aliases as string[],
      defaultName: command.defaultName ?? null,
    };

    await this.redis.set(
      `commands:${command.channelId}:${command.name}`,
      JSON.stringify(commandForSet),
    );
  }

  async #isCommandWithThatNameExists(opts: {
    userId: string;
    name: string;
    aliases?: string[];
    commandId?: string;
  }) {
    opts.aliases = opts.aliases ?? [];
    const unique = [...new Set([opts.name, ...opts.aliases])];

    const commands = await this.getList(opts.userId);
    return commands
      .filter((c) => c.id !== opts.commandId)
      .some((c) => {
        return unique.includes(c.name) || opts.aliases!.some((a) => c.aliases.includes(a));
      });
  }

  async create(userId: string, data: UpdateOrCreateCommandDto & { defaultName?: string }) {
    const isExists = await this.#isCommandWithThatNameExists({
      userId,
      name: data.name,
      aliases: data.aliases,
    });

    if (isExists) {
      throw new HttpException(`Command with that name or aliase already exists`, 400);
    }

    if (!data.responses?.length) {
      throw new HttpException(`You should add atleast 1 response to command.`, 400);
    }

    const command = await typeorm.getRepository(ChannelCommand).save({
      ...data,
      channelId: userId,
    });

    command.responses = await typeorm.getRepository(CommandResponse).save(
      data.responses
        .filter((r) => r.text)
        .map((r) => ({
          commandId: command.id,
          text: r.text?.trim().replace(/(\r\n|\n|\r)/, ''),
        })),
    );

    await this.setCommandCache(command);
    return command as ChannelCommand;
  }

  async delete(userId: string, commandId: string) {
    const command = await typeorm
      .getRepository(ChannelCommand)
      .findOneBy({ channelId: userId, id: commandId });

    if (!command) {
      throw new HttpException('Command not exists', 404);
    }

    if (command.default) {
      throw new HttpException('You cannot delete default command.', 400);
    }

    const result = await typeorm.getRepository(ChannelCommand).delete({
      id: commandId,
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
    const isExists = await this.#isCommandWithThatNameExists({
      userId,
      name: data.name,
      aliases: data.aliases,
      commandId,
    });

    if (isExists) {
      throw new HttpException(`Command with this name or aliase already exists`, 400);
    }

    const command = await typeorm
      .getRepository(ChannelCommand)
      .findOneBy({ channelId: userId, id: commandId });

    if (!command) {
      throw new HttpException('Command not exists', 404);
    }

    if (!data.responses?.length && !command.default) {
      throw new HttpException(`You should add atleast 1 response to command.`, 400);
    }

    await typeorm.getRepository(CommandResponse).delete({
      commandId: command.id,
    });

    await typeorm.getRepository(CommandResponse).save(
      data.responses
        ?.filter((r) => r.text)
        .map((r) => ({
          commandId: command.id,
          text: r.text ? r.text.trim().replace(/(\r\n|\n|\r)/, '') : null,
        })),
    );

    await typeorm.getRepository(ChannelCommand).update(
      { id: command.id },
      {
        ...data,
        responses: undefined,
      },
    );

    const newCommand = await typeorm.getRepository(ChannelCommand).findOne({
      where: { id: command.id },
      relations: {
        responses: true,
      },
    });

    await this.setCommandCache(newCommand!, command);

    return newCommand;
  }
}

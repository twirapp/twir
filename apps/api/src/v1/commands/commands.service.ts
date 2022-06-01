import { Inject, Injectable, OnModuleInit } from '@nestjs/common';
import { ClientGrpc } from '@nestjs/microservices';
import { Bots } from '@tsuwari/grpc';
import { PrismaService } from '@tsuwari/prisma';

import { UpdateOrCreateCommandDto } from './dto/create.js';

@Injectable()
export class CommandsService implements OnModuleInit {
  private botsMicroservice: Bots.Commands;

  constructor(private readonly prisma: PrismaService, @Inject('BOTS_MICROSERVICE') private client: ClientGrpc) { }

  onModuleInit(): void {
    this.botsMicroservice = this.client.getService<Bots.Commands>('Commands');
  }

  findOne(userId: string, commandId: string) {
    return this.prisma.command.findFirst({
      where: {
        channelId: userId,
        id: commandId,
      },
      include: {
        responses: true,
      },
    });
  }

  async getList(userId: string) {
    const commands = await this.prisma.command.findMany({
      where: { channelId: userId },
      include: {
        responses: true,
      },
    });
    return commands;
  }

  async create(userId: string, data: UpdateOrCreateCommandDto) {
    if (await this.prisma.command.count({ where: { channelId: userId, name: data.name } })) {
      throw new Error('Command already exists');
    }

    const command = await this.prisma.command.create({
      data: {
        ...data,
        channelId: userId,
        responses: {
          createMany: {
            data: data.responses.map((r) => ({ text: r.text })),
          },
        },
      },
      include: {
        responses: true,
      },
    });

    await this.botsMicroservice.updateByChannelId({ userId }).toPromise();

    return command;
  }

  async delete(userId: string, commandId: string) {
    const command = await this.prisma.command.findFirst({ where: { channelId: userId, id: commandId } });

    if (!command) {
      throw new Error('Command not exists');
    }

    const result = await this.prisma.command.delete({
      where: {
        id: commandId,
      },
    });

    await this.botsMicroservice.updateByChannelId({ userId }).toPromise();

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

    const responsesForUpdate = data.responses
      .filter(r => command.responses.some(c => c.id === r.id && r.text && r.id))
      .map(r => ({ id: r.id, text: r.text }))
      .map(r => this.prisma.response.update({ where: { id: r.id }, data: { text: r.text } }));

    const [newCommand, , ...newResponses] = await this.prisma.$transaction([
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
      this.prisma.response.findMany({ where: { commandId: command.id } }),
    ]);

    await this.botsMicroservice.updateByChannelId({ userId }).toPromise();

    return {
      ...newCommand,
      responses: newResponses.flat(),
    };
  }
}

import { Inject, Injectable, OnModuleInit } from '@nestjs/common';
import { ClientGrpc } from '@nestjs/microservices';
import { Bots } from '@tsuwari/grpc';
import { PrismaService } from '@tsuwari/prisma';

import { staticApi } from '../../twitchApi.js';
import { GreetingCreateDto } from './dto/create.js';

@Injectable()
export class GreetingsService implements OnModuleInit {
  private botsMicroservice: Bots.Greetings;

  constructor(private readonly prisma: PrismaService, @Inject('BOTS_MICROSERVICE') private client: ClientGrpc) { }

  onModuleInit(): void {
    this.botsMicroservice = this.client.getService<Bots.Greetings>('Greetings');
  }

  async getList(userId: string) {
    const greetings = await this.prisma.greeting.findMany({
      where: { channelId: userId },
    });

    const users = await staticApi.users.getUsersByIds(greetings.map(g => g.userId));

    return greetings.map(g => ({ ...g, username: users.find(u => u.id === g.userId)?.name }));
  }

  async create(userId: string, data: GreetingCreateDto) {
    const user = await staticApi.users.getUserByName(data.username);

    if (!user) throw new Error(`User ${data.username} not found on twitch`);

    const isExists = await this.prisma.greeting.count({
      where: {
        channelId: userId,
        userId: user.id,
      },
    });

    if (isExists) {
      throw new Error('Greeting already exists');
    }

    const greeting = await this.prisma.greeting.create({
      data: {
        channelId: userId,
        userId: user.id,
        text: data.text,
      },
    });

    await this.botsMicroservice.updateByChannelId({ userId }).toPromise();

    return {
      ...greeting,
      username: user.name,
    };
  }

  async update(userId: string, greetingId: string, data: GreetingCreateDto) {
    const currentGreeting = await this.prisma.greeting.count({
      where: {
        id: greetingId,
        channelId: userId,
      },
    });

    if (!currentGreeting) throw new Error(`Greeting with id ${greetingId} not found.`);

    const user = await staticApi.users.getUserByName(data.username);

    if (!user) throw new Error(`User ${data.username} not found on twitch`);

    const greeting = await this.prisma.greeting.update({
      where: {
        id: greetingId,
      },
      data: {
        text: data.text,
        userId: user.id,
        enabled: data.enabled,
      },
    });

    await this.botsMicroservice.updateByChannelId({ userId }).toPromise();

    return {
      ...greeting,
      username: user.name,
    };
  }

  async delete(userId: string, greetingId: string) {
    const command = await this.prisma.greeting.findFirst({ where: { channelId: userId, id: greetingId } });

    if (!command) {
      throw new Error('Greeting not exists');
    }

    const result = await this.prisma.greeting.delete({
      where: {
        id: greetingId,
      },
    });

    await this.botsMicroservice.updateByChannelId({ userId }).toPromise();

    return result;
  }
}

import { HttpException, Injectable } from '@nestjs/common';
import { PrismaService } from '@tsuwari/prisma';

import { RedisService } from '../../redis.service.js';
import { staticApi } from '../../twitchApi.js';
import { GreetingCreateDto } from './dto/create.js';

@Injectable()
export class GreetingsService {

  constructor(
    private readonly prisma: PrismaService,
    private readonly redis: RedisService,
  ) { }

  async getList(userId: string) {
    const greetings = await this.prisma.greeting.findMany({
      where: { channelId: userId },
    });

    const users = await staticApi.users.getUsersByIds(greetings.map(g => g.userId));

    return greetings.map(g => ({ ...g, username: users.find(u => u.id === g.userId)?.name }));
  }

  async create(userId: string, data: GreetingCreateDto) {
    const user = await staticApi.users.getUserByName(data.username);

    if (!user) throw new HttpException(`User ${data.username} not found on twitch`, 404);

    const isExists = await this.prisma.greeting.count({
      where: {
        channelId: userId,
        userId: user.id,
      },
    });

    if (isExists) {
      throw new HttpException(`Greeting for user ${user.name} already exists`, 400);
    }

    const greeting = await this.prisma.greeting.create({
      data: {
        channelId: userId,
        userId: user.id,
        text: data.text,
      },
    });

    await this.redis.hset(`greetings:${greeting.channelId}:${greeting.userId}`, {
      ...greeting,
      processed: false,
    });

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

    await this.redis.hset(`greetings:${greeting.channelId}:${greeting.userId}`, {
      ...greeting,
      processed: false,
    });

    return {
      ...greeting,
      username: user.name,
    };
  }

  async delete(userId: string, greetingId: string) {
    const greeting = await this.prisma.greeting.findFirst({ where: { channelId: userId, id: greetingId } });

    if (!greeting) {
      throw new Error('Greeting not exists');
    }

    const result = await this.prisma.greeting.delete({
      where: {
        id: greetingId,
      },
    });

    await this.redis.del(`greetings:${greeting.channelId}:${greeting.userId}`);

    return result;
  }
}

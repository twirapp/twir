import { HttpException, Injectable, OnModuleInit } from '@nestjs/common';
import { Greetings, greetingsSchema, RedisORMService, Repository } from '@tsuwari/redis';
import { RedisService } from '@tsuwari/shared';
import { ChannelGreeting } from '@tsuwari/typeorm/entities/ChannelGreeting';

import { typeorm } from '../../index.js';
import { staticApi } from '../../twitchApi.js';
import { GreetingCreateDto } from './dto/create.js';

@Injectable()
export class GreetingsService implements OnModuleInit {
  #repository: Repository<Greetings>;

  constructor(private readonly redis: RedisService, private readonly redisOrm: RedisORMService) {}

  async onModuleInit() {
    await this.redisOrm.onModuleInit();
    this.#repository = this.redisOrm.fetchRepository(greetingsSchema);
  }

  async getList(userId: string) {
    const greetings = await typeorm.getRepository(ChannelGreeting).findBy({ channelId: userId });

    const users = await staticApi.users.getUsersByIds(greetings.map((g) => g.userId));

    return greetings.map((g) => ({ ...g, username: users.find((u) => u.id === g.userId)?.name }));
  }

  async create(userId: string, data: GreetingCreateDto) {
    const user = await staticApi.users.getUserByName(data.username);

    if (!user) throw new HttpException(`User ${data.username} not found on twitch`, 404);

    const isExists = await typeorm.getRepository(ChannelGreeting).findOneBy({
      channelId: userId,
      userId: user.id,
    });

    if (isExists) {
      throw new HttpException(`Greeting for user ${user.name} already exists`, 400);
    }

    const greeting = await typeorm.getRepository(ChannelGreeting).save({
      channelId: userId,
      userId: user.id,
      text: data.text,
    });

    await this.#repository.createAndSave(
      {
        ...greeting,
        processed: false,
      },
      `${greeting.channelId}:${greeting.userId}`,
    );

    return {
      ...greeting,
      username: user.name,
    };
  }

  async update(userId: string, greetingId: string, data: GreetingCreateDto) {
    const repository = typeorm.getRepository(ChannelGreeting);
    const currentGreeting = await repository.findOneBy({
      id: greetingId,
      channelId: userId,
    });

    if (!currentGreeting) throw new HttpException(`Greeting with id ${greetingId} not found.`, 404);

    const user = await staticApi.users.getUserByName(data.username);

    if (!user) throw new HttpException(`User ${data.username} not found on twitch`, 404);

    await repository.update(
      { id: greetingId },
      { text: data.text, userId: user.id, enabled: data.enabled },
    );

    const greeting = await repository.findOneBy({ id: greetingId });

    await this.#repository.createAndSave(
      {
        ...greeting!,
        processed: false,
      },
      `${greeting!.channelId}:${greeting!.userId}`,
    );

    return {
      ...greeting,
      username: user.name,
    };
  }

  async delete(userId: string, greetingId: string) {
    const repository = typeorm.getRepository(ChannelGreeting);
    const greeting = await repository.findOneBy({
      channelId: userId,
      id: greetingId,
    });

    if (!greeting) {
      throw new HttpException('Greeting not exists', 404);
    }

    const result = await repository.delete({
      id: greetingId,
    });

    await this.#repository.remove(`${greeting.channelId}:${greeting.userId}`);

    return result;
  }
}

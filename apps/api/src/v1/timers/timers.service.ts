import { HttpException, Injectable } from '@nestjs/common';
import * as TimersNats from '@tsuwari/nats/timers';
import { ChannelTimer } from '@tsuwari/typeorm/entities/ChannelTimer';

import { typeorm } from '../../index.js';
import { nats } from '../../libs/nats.js';
import { CreateTimerDto } from './dto/create.js';

@Injectable()
export class TimersService {
  getList(channelId: string) {
    return typeorm.getRepository(ChannelTimer).findBy({
      channelId,
    });
  }

  findOne(channelId: string, id: string) {
    return typeorm.getRepository(ChannelTimer).findOneBy({
      channelId,
      id,
    });
  }

  async create(userId: string, data: CreateTimerDto) {
    const repository = typeorm.getRepository(ChannelTimer);
    const isExists = await repository.count({
      where: {
        channelId: userId,
        name: data.name,
      },
    });

    if (isExists) throw new HttpException(`Timer with name ${data.name} already exists`, 400);

    const timer = await repository.save({
      ...data,
      channelId: userId,
    });

    if (timer.enabled) {
      const data = TimersNats.AddTimerToQueue.toBinary({
        timerId: timer.id,
      });
      await nats.request('addTimerToQueue', data);
    }

    return timer;
  }

  async delete(userId: string, timerId: string) {
    const repository = typeorm.getRepository(ChannelTimer);
    const isExists = await repository.count({
      where: {
        channelId: userId,
        id: timerId,
      },
    });

    if (!isExists) throw new HttpException(`Timer with id ${timerId} not exists`, 404);

    await repository.delete({
      id: timerId,
    });

    const data = TimersNats.RemoveTimerFromQueue.toBinary({
      timerId,
    });
    await nats.request('removeTimerFromQueue', data);

    return true;
  }

  async update(userId: string, timerId: string, data: CreateTimerDto) {
    const repository = typeorm.getRepository(ChannelTimer);
    const isExists = await repository.count({
      where: {
        channelId: userId,
        id: timerId,
      },
    });

    if (!isExists) throw new HttpException(`Timer with id ${timerId} not exists`, 404);

    await repository.update({ id: timerId }, data);

    const timer = await this.findOne(userId, timerId);
    if (timer!.enabled) {
      const data = TimersNats.AddTimerToQueue.toBinary({
        timerId: timer!.id,
      });
      await nats.request('addTimerToQueue', data);
    } else {
      const data = TimersNats.RemoveTimerFromQueue.toBinary({
        timerId: timer!.id,
      });
      await nats.request('removeTimerFromQueue', data);
    }

    return timer!;
  }
}

import { HttpException, Injectable } from '@nestjs/common';
import * as TimersNats from '@tsuwari/nats/timers';
import { ChannelTimer } from '@tsuwari/typeorm/entities/ChannelTimer';
import { ChannelTimerResponse } from '@tsuwari/typeorm/entities/ChannelTimerResponse';

import { typeorm } from '../../index.js';
import { nats } from '../../libs/nats.js';
import { CreateTimerDto } from './dto/create.js';

@Injectable()
export class TimersService {
  getList(channelId: string) {
    return typeorm.getRepository(ChannelTimer).find({
      where: { channelId },
      relations: { responses: true },
    });
  }

  findOne(channelId: string, id: string) {
    return typeorm.getRepository(ChannelTimer).findOne({
      where: { channelId, id },
      relations: { responses: true },
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

    const timer = repository.create({
      channelId: userId,
      enabled: data.enabled,
      messageInterval: data.messageInterval,
      name: data.name,
      timeInterval: data.timeInterval,
      lastTriggerMessageNumber: 0,
    });

    const created = await repository.save(timer);

    const responses = data.responses.map((r) => ({ ...r, timerId: created.id }));

    await typeorm.getRepository(ChannelTimerResponse).save(responses);

    if (created.enabled) {
      const data = TimersNats.AddTimerToQueue.toBinary({
        timerId: created.id,
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

    await typeorm.getRepository(ChannelTimerResponse).delete({
      timerId,
    });
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

    await repository.update(
      { id: timerId },
      {
        enabled: data.enabled,
        messageInterval: data.messageInterval,
        name: data.name,
        timeInterval: data.timeInterval,
        lastTriggerMessageNumber: 0,
      },
    );

    const responsesRepository = typeorm.getRepository(ChannelTimerResponse);
    await responsesRepository.delete({ timerId: timerId });
    const responses = data.responses.map((r) => ({ ...r, timerId: timerId }));
    await responsesRepository.save(responses);

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

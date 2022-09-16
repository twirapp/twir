import { HttpException, Injectable } from '@nestjs/common';
import { Client, Transport } from '@nestjs/microservices';
import { config } from '@tsuwari/config';
import { ClientProxy } from '@tsuwari/shared';
import { ChannelTimer } from '@tsuwari/typeorm/entities/ChannelTimer';

import { typeorm } from '../../index.js';
import { CreateTimerDto } from './dto/create.js';

@Injectable()
export class TimersService {
  @Client({ transport: Transport.NATS, options: { servers: [config.NATS_URL] } })
  nats: ClientProxy;

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
      await this.nats.emit('bots.addTimerToQueue', timer.id).toPromise();
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

    await this.nats.emit('bots.removeTimerFromQueue', timerId).toPromise();

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
      await this.nats.emit('bots.addTimerToQueue', timer!.id).toPromise();
    } else {
      await this.nats.emit('bots.removeTimerFromQueue', timer!.id).toPromise();
    }

    return timer!;
  }
}

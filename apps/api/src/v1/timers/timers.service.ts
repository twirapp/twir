import { HttpException, Injectable } from '@nestjs/common';
import { Client, Transport } from '@nestjs/microservices';
import { config } from '@tsuwari/config';
import { PrismaService } from '@tsuwari/prisma';
import { ClientProxy } from '@tsuwari/shared';

import { CreateTimerDto } from './dto/create.js';

@Injectable()
export class TimersService {
  @Client({ transport: Transport.NATS, options: { servers: [config.NATS_URL] } })
  nats: ClientProxy;

  constructor(private readonly prisma: PrismaService) { }

  getList(userId: string) {
    return this.prisma.timer.findMany({
      where: {
        channelId: userId,
      },
    });
  }

  findOne(userId: string, timerId: string) {
    return this.prisma.timer.findFirst({
      where: {
        channelId: userId,
        id: timerId,
      },
    });
  }

  async create(userId: string, data: CreateTimerDto) {
    const isExists = await this.prisma.timer.count({
      where: {
        channelId: userId,
        name: data.name,
      },
    });

    if (isExists) throw new HttpException(`Timer with name ${data.name} already exists`, 400);

    const timer = await this.prisma.timer.create({
      data: {
        ...data,
        channelId: userId,
      },
    });

    await this.nats.emit('bots.addTimerToQueue', timer.id).toPromise();
    return timer;
  }

  async delete(userId: string, timerId: string) {
    const isExists = await this.prisma.timer.count({
      where: {
        channelId: userId,
        id: timerId,
      },
    });

    if (!isExists) throw new HttpException(`Timer with id ${timerId} not exists`, 404);

    const timer = await this.prisma.timer.delete({
      where: {
        id: timerId,
      },
    });

    await this.nats.emit('bots.removeTimerFromQueue', timer.id).toPromise();

    return timer;
  }

  async update(userId: string, timerId: string, data: CreateTimerDto) {
    const isExists = await this.prisma.timer.count({
      where: {
        channelId: userId,
        id: timerId,
      },
    });

    if (!isExists) throw new HttpException(`Timer with id ${timerId} not exists`, 404);

    const updated = await this.prisma.timer.update({
      where: {
        id: timerId,
      },
      data: {
        ...data,
      },
    });

    return updated;
  }
}

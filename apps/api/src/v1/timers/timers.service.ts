import { Inject, Injectable, OnModuleInit } from '@nestjs/common';
import { ClientGrpc } from '@nestjs/microservices';
import { Bots } from '@tsuwari/grpc';
import { PrismaService } from '@tsuwari/prisma';

import { CreateTimerDto } from './dto/create.js';

@Injectable()
export class TimersService implements OnModuleInit {
  botsMicroservice: Bots.Timers;

  constructor(private readonly prisma: PrismaService, @Inject('BOTS_MICROSERVICE') private readonly botsClient: ClientGrpc) { }


  onModuleInit() {
    this.botsMicroservice = this.botsClient.getService<Bots.Timers>('Timers');
  }


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

    if (isExists) throw new Error(`Timer with name ${data.name} already exists`);

    const timer = await this.prisma.timer.create({
      data: {
        ...data,
        channelId: userId,
      },
    });

    await this.botsMicroservice.addTimerToQueue({ timerId: timer.id }).toPromise();

    return timer;
  }

  async delete(userId: string, timerId: string) {
    const isExists = await this.prisma.timer.count({
      where: {
        channelId: userId,
        id: timerId,
      },
    });

    if (!isExists) throw new Error(`Timer with id ${timerId} not exists`);

    const timer = await this.prisma.timer.delete({
      where: {
        id: timerId,
      },
    });

    await this.botsMicroservice.removeTimerFromQueue({ timerId: timer.id }).toPromise();

    return timer;
  }

  async update(userId: string, timerId: string, data: CreateTimerDto) {
    const isExists = await this.prisma.timer.count({
      where: {
        channelId: userId,
        id: timerId,
      },
    });

    if (!isExists) throw new Error(`Timer with id ${timerId} not exists`);

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

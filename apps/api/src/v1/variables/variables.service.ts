import { HttpException, Injectable } from '@nestjs/common';
import { Client, ClientProxy, Transport } from '@nestjs/microservices';
import { config } from '@tsuwari/config';
import { PrismaService } from '@tsuwari/prisma';


import { RedisService } from '../../redis.service.js';
import { CreateVariableDto } from './dto/create.js';

@Injectable()
export class VariablesService {
  @Client({ transport: Transport.NATS, options: { servers: [config.NATS_URL] } })
  nats: ClientProxy;

  constructor(
    private readonly prisma: PrismaService,
    private readonly redis: RedisService,
  ) { }

  async getBuildInVariables() {
    const list = await this.nats.send('bots.getVariables', {}).toPromise();

    return list?.variables;
  }

  getList(channelId: string) {
    return this.prisma.customVar.findMany({ where: { channelId } });
  }

  async create(channelId: string, data: CreateVariableDto) {
    const isExists = await this.prisma.customVar.count({
      where: {
        channelId,
        name: data.name,
      },
    });

    if (isExists) throw new HttpException(`Variable with name ${data.name} already exists`, 400);

    const variable = await this.prisma.customVar.create({
      data: {
        channelId,
        ...data,
      },
    });

    await this.redis.set(`variables:${channelId}:${variable.name}`, JSON.stringify(variable));

    return variable;
  }

  async delete(channelId: string, variableId: string) {
    const variable = await this.prisma.customVar.findFirst({
      where: {
        channelId,
        id: variableId,
      },
    });

    if (!variable) throw new HttpException(`Variable with id ${variableId} not exists`, 404);

    await this.prisma.customVar.delete({
      where: { id: variableId },
    });

    await this.redis.del(`variables:${channelId}:${variable.name}`);

    return variable;
  }

  async update(channelId: string, variableId: string, data: CreateVariableDto) {
    const variable = await this.prisma.customVar.findFirst({
      where: {
        channelId,
        id: variableId,
      },
    });

    if (!variable) throw new HttpException(`Variable with id ${variableId} not exists`, 404);

    const newVariable = await this.prisma.customVar.update({
      where: {
        id: variable.id,
      },
      data,
    });

    await this.redis.set(`variables:${channelId}:${variable.name}`, JSON.stringify(newVariable));

    return newVariable;
  }
}

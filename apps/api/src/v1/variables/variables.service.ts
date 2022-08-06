import { HttpException, Injectable, OnModuleInit } from '@nestjs/common';
import { Client, Transport } from '@nestjs/microservices';
import { config } from '@tsuwari/config';
import { PrismaService } from '@tsuwari/prisma';
import { CustomVar, customVarSchema, RedisORMService, Repository } from '@tsuwari/redis';
import { ClientProxy, RedisService } from '@tsuwari/shared';


import { CreateVariableDto } from './dto/create.js';

@Injectable()
export class VariablesService implements OnModuleInit {
  @Client({ transport: Transport.NATS, options: { servers: [config.NATS_URL] } })
  nats: ClientProxy;
  #repository: Repository<CustomVar>;

  constructor(
    private readonly prisma: PrismaService,
    private readonly redis: RedisService,
    private readonly redisOrm: RedisORMService,
  ) { }

  onModuleInit() {
    this.#repository = this.redisOrm.fetchRepository(customVarSchema);
  }

  async getBuildInVariables() {
    const list = await this.nats.send('bots.getVariables', {}).toPromise();

    return list;
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

    await this.#repository.createAndSave(variable, `${channelId}:${variable.id}`);

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

    await this.#repository.remove(`${channelId}:${variable.id}`);

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

    await this.#repository.createAndSave(newVariable, `${channelId}:${variable.id}`);

    return newVariable;
  }
}

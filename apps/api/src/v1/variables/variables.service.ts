import { HttpException, Injectable } from '@nestjs/common';
import { Client, Transport } from '@nestjs/microservices';
import { config } from '@tsuwari/config';
import * as Parser from '@tsuwari/nats/parser';
import { ClientProxy, RedisService } from '@tsuwari/shared';
import { ChannelCustomvar } from '@tsuwari/typeorm/entities/ChannelCustomvar';

import { typeorm } from '../../index.js';
import { nats } from '../../libs/nats.js';
import { CreateVariableDto } from './dto/create.js';

@Injectable()
export class VariablesService {
  @Client({ transport: Transport.NATS, options: { servers: [config.NATS_URL] } })
  nats: ClientProxy;

  readonly #vulnerableScriptWords = ['prototype', 'contructor'];

  constructor(private readonly redis: RedisService) {}

  #checkNotSecureVariable(script: string) {
    if (this.#vulnerableScriptWords.some((w) => script.includes(w))) {
      throw new HttpException(
        `You cannot use such vulnerable words in script. You will be banned, if you attemp to abuse this feature.`,
        400,
      );
    }

    return true;
  }

  async getBuildInVariables() {
    const msg = await nats.request('bots.getVariables', new Uint8Array());
    const list = Parser.GetVariablesResponse.fromBinary(msg.data);

    return list.list;
  }

  getList(channelId: string) {
    return typeorm.getRepository(ChannelCustomvar).findBy({ channelId });
  }

  async create(channelId: string, data: CreateVariableDto) {
    const repository = typeorm.getRepository(ChannelCustomvar);
    const isExists = await repository.count({
      where: {
        channelId,
        name: data.name,
      },
    });

    if (isExists) throw new HttpException(`Variable with name ${data.name} already exists`, 400);

    const variable = await repository.save({
      channelId,
      ...data,
    });

    if (data.type === 'SCRIPT' && data.evalValue) {
      this.#checkNotSecureVariable(data.evalValue);
    }

    await this.redis.set(`variables:${channelId}:${variable.name}`, JSON.stringify(variable));
    return variable;
  }

  async delete(channelId: string, variableId: string) {
    const repository = typeorm.getRepository(ChannelCustomvar);

    const variable = await repository.findOneBy({
      channelId,
      id: variableId,
    });

    if (!variable) throw new HttpException(`Variable with id ${variableId} not exists`, 404);

    await repository.delete({
      id: variableId,
    });

    await this.redis.del(`variables:${channelId}:${variable.name}`);

    return variable;
  }

  async update(channelId: string, variableId: string, data: CreateVariableDto) {
    const repository = typeorm.getRepository(ChannelCustomvar);

    const variable = await repository.findOneBy({
      channelId,
      id: variableId,
    });

    if (!variable) throw new HttpException(`Variable with id ${variableId} not exists`, 404);

    if (data.type === 'SCRIPT' && data.evalValue) {
      this.#checkNotSecureVariable(data.evalValue);
    }

    await repository.update({ id: variable.id }, data);

    const newVariable = await repository.findOneBy({ id: variable.id });
    await this.redis.set(`variables:${channelId}:${variable!.name}`, JSON.stringify(newVariable));

    return newVariable;
  }
}

import { HttpException, Injectable, OnModuleInit } from '@nestjs/common';
import { ChannelKeyword } from '@tsuwari/typeorm/entities/ChannelKeyword';

import { typeorm } from '../../index.js';
import { CreateKeywordDto } from './dto/create.js';

@Injectable()
export class KeywordsService {
  async getList(channelId: string) {
    return typeorm.getRepository(ChannelKeyword).find({
      where: { channelId },
    });
  }

  async delete(channelId: string, keywordId: string) {
    const repository = typeorm.getRepository(ChannelKeyword);
    const keyword = await repository.findOneBy({
      channelId,
      id: keywordId,
    });

    if (!keyword) throw new HttpException(`Keyword with id ${keywordId} not found`, 404);

    await repository.delete({
      id: keyword.id,
    });
  }

  async create(channelId: string, data: CreateKeywordDto) {
    const repository = typeorm.getRepository(ChannelKeyword);
    const isExists = await repository.findOneBy({
      channelId,
      text: data.text,
    });

    if (isExists) throw new HttpException(`Keyword with text ${data.text} already exists`, 400);

    const keyword = await repository.save({
      channelId,
      ...data,
    });

    return keyword;
  }

  async patch(channelId: string, keywordId: string, data: CreateKeywordDto) {
    const repository = typeorm.getRepository(ChannelKeyword);
    const isExists = await repository.findOneBy({
      id: keywordId,
    });

    if (!isExists) throw new HttpException(`Keyword with id ${keywordId} not exists`, 404);

    await repository.update({ id: keywordId }, data);
    const newKeyword = await repository.findOneBy({ id: keywordId });

    return newKeyword;
  }
}

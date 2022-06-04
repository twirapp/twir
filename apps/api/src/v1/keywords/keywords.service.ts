import { HttpException, Injectable } from '@nestjs/common';
import { PrismaService } from '@tsuwari/prisma';

import { RedisService } from '../../redis.service.js';
import { CreateKeywordDto } from './dto/create.js';

@Injectable()
export class KeywordsService {
  constructor(
    private readonly prisma: PrismaService,
    private readonly redis: RedisService,
  ) { }

  async getList(channelId: string) {
    return this.prisma.keyword.findMany({
      where: {
        channelId,
      },
    });
  }

  async delete(channelId: string, keywordId: string) {
    const keyword = await this.prisma.keyword.findFirst({
      where: {
        channelId,
        id: keywordId,
      },
    });

    if (!keyword) throw new HttpException(`Keyword with id ${keywordId} not found`, 404);

    await this.prisma.keyword.delete({
      where: {
        id: keyword.id,
      },
    });

    await this.redis.del(`keywords:${channelId}:${keyword.id}`);
  }

  async create(channelId: string, data: CreateKeywordDto) {
    const isExists = await this.prisma.keyword.findUnique({
      where: {
        channelId_text: {
          channelId,
          text: data.text,
        },
      },
    });

    if (isExists) throw new HttpException(`Keyword with text ${data.text} already exists`, 400);

    const keyword = await this.prisma.keyword.create({
      data: {
        channelId,
        ...data,
      },
    });

    await this.redis.hmset(`keywords:${channelId}:${keyword.id}`, keyword);
    return keyword;
  }

  async patch(channelId: string, keywordId: string, data: CreateKeywordDto) {
    const isExists = await this.prisma.keyword.findUnique({
      where: {
        id: keywordId,
      },
    });

    if (!isExists) throw new HttpException(`Keyword with id ${keywordId} not exists`, 404);

    const keyword = await this.prisma.keyword.update({
      where: {
        id: keywordId,
      },
      data,
    });

    await this.redis.hmset(`keywords:${channelId}:${keyword.id}`, keyword);
    return keyword;
  }
}

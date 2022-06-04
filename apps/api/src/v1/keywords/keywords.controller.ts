import { Body, CacheTTL, CACHE_MANAGER, Controller, Get, Inject, Param, Patch, Post, UseGuards, UseInterceptors } from '@nestjs/common';
import { Request } from 'express';

import { DashboardAccessGuard } from '../../guards/DashboardAccess.guard.js';
import { CustomCacheInterceptor } from '../../helpers/customCacheInterceptor.js';
import { JwtAuthGuard } from '../../jwt/jwt.guard.js';
import { CreateKeywordDto } from './dto/create.js';
import { KeywordsService } from './keywords.service.js';


@Controller('v1/channels/:channelId/keywords')
export class KeywordsController {
  constructor(
    @Inject(CACHE_MANAGER) private cacheManager: Cache,
    @Inject() private readonly keywordsService: KeywordsService,
  ) { }

  async #delCache(channelId: string) {
    await this.cacheManager.delete(`nest:cache:v1/channels/${channelId}/keywords`);
  }

  @CacheTTL(600)
  @UseInterceptors(CustomCacheInterceptor(ctx => {
    const req = ctx.switchToHttp().getRequest() as Request;
    return `nest:cache:v1/channels/${req.params.channelId}/keywords`;
  }))
  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Get()
  root(@Param('channelId') channelId: string) {
    return this.keywordsService.getList(channelId);
  }

  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Get()
  async delete(@Param('channelId') channelId: string, @Param('keywordId') keywordId: string) {
    const result = await this.keywordsService.delete(channelId, keywordId);
    await this.#delCache(channelId);
    return result;
  }

  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Post()
  async create(@Param('channelId') channelId: string, @Body() body: CreateKeywordDto) {
    const result = await this.keywordsService.create(channelId, body);
    await this.#delCache(channelId);
    return result;
  }

  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Patch(':keywordId')
  async patch(@Param('channelId') channelId: string, @Param('keywordId') keywordId: string, @Body() body: CreateKeywordDto) {
    const result = await this.keywordsService.patch(channelId, keywordId, body);
    await this.#delCache(channelId);
    return result;
  }

}

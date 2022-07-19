import { Body, CacheTTL, CACHE_MANAGER, Controller, Delete, Get, Inject, Param, Post, Put, UseGuards, UseInterceptors } from '@nestjs/common';
import CacheManager from 'cache-manager';
import Express from 'express';

import { DashboardAccessGuard } from '../../guards/DashboardAccess.guard.js';
import { CustomCacheInterceptor } from '../../helpers/customCacheInterceptor.js';
import { JwtAuthGuard } from '../../jwt/jwt.guard.js';
import { GreetingCreateDto } from './dto/create.js';
import { GreetingsService } from './greetings.service.js';

@Controller('v1/channels/:channelId/greetings')
export class GreetingsController {
  constructor(
    private readonly greetingsService: GreetingsService,
    @Inject(CACHE_MANAGER) private cacheManager: CacheManager.Cache,
  ) { }

  async #delCache(channelId: string) {
    await this.cacheManager.del(`nest:cache:v1/channels/${channelId}/greetings`);
  }

  @CacheTTL(600)
  @UseInterceptors(CustomCacheInterceptor(ctx => {
    const req = ctx.switchToHttp().getRequest() as Express.Request;
    return `nest:cache:v1/channels/${req.params.channelId}/greetings`;
  }))
  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Get()
  root(@Param('channelId') channelId: string) {
    return this.greetingsService.getList(channelId);
  }

  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Post()
  async create(@Param('channelId') channelId: string, @Body() body: GreetingCreateDto) {
    const result = await this.greetingsService.create(channelId, body);
    await this.#delCache(channelId);
    return result;
  }

  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Put(':greetingId')
  async update(@Param('channelId') channelId: string, @Param('greetingId') greetingId: string, @Body() body: GreetingCreateDto) {
    const result = await this.greetingsService.update(channelId, greetingId, body);
    await this.#delCache(channelId);
    return result;
  }

  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Delete(':greetingId')
  async delete(@Param('channelId') channelId: string, @Param('greetingId') greetingId: string) {
    const result = await this.greetingsService.delete(channelId, greetingId);
    await this.#delCache(channelId);
    return result;
  }
}

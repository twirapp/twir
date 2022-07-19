import { Body, CacheTTL, CACHE_MANAGER, Controller, Delete, Get, Inject, Param, Post, Put, UseGuards, UseInterceptors } from '@nestjs/common';
import CacheManager from 'cache-manager';
import Express from 'express';

import { DashboardAccessGuard } from '../../guards/DashboardAccess.guard.js';
import { CustomCacheInterceptor } from '../../helpers/customCacheInterceptor.js';
import { JwtAuthGuard } from '../../jwt/jwt.guard.js';
import { CreateTimerDto } from './dto/create.js';
import { TimersService } from './timers.service.js';

@Controller('v1/channels/:channelId/timers')
export class TimersController {
  constructor(
    private readonly timersService: TimersService,
    @Inject(CACHE_MANAGER) private cacheManager: CacheManager.Cache,
  ) { }

  private async delCache(channelId: string) {
    await this.cacheManager.del(`nest:cache:v1/channels/${channelId}/timers`);
  }

  @CacheTTL(600)
  @UseInterceptors(CustomCacheInterceptor(ctx => {
    const req = ctx.switchToHttp().getRequest() as Express.Request;
    return `nest:cache:v1/channels/${req.params.channelId}/timers`;
  }))
  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Get()
  root(@Param('channelId') channelId: string) {
    return this.timersService.getList(channelId);
  }

  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Get(':id')
  async findOne(@Param('channelId') channelId: string, @Param('id') id: string) {
    const result = await this.timersService.findOne(channelId, id);
    this.delCache(channelId);
    return result;
  }

  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Post()
  async create(@Param('channelId') channelId: string, @Body() body: CreateTimerDto) {
    const result = await this.timersService.create(channelId, body);
    await this.delCache(channelId);
    return result;
  }

  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Delete(':id')
  async delete(@Param('channelId') channelId: string, @Param('id') id: string) {
    const result = await this.timersService.delete(channelId, id);
    await this.delCache(channelId);
    return result;
  }

  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Put(':id')
  async update(@Param('channelId') channelId: string, @Param('id') id: string, @Body() body: CreateTimerDto) {
    const result = await this.timersService.update(channelId, id, body);
    await this.delCache(channelId);
    return result;
  }
}

import { Body, CacheTTL, CACHE_MANAGER, Controller, Get, Inject, Param, Post, UseGuards, UseInterceptors } from '@nestjs/common';
import { ModerationUpdateDto } from '@tsuwari/shared';
import { Cache } from 'cache-manager';
import { Request } from 'express';

import { DashboardAccessGuard } from '../../guards/DashboardAccess.guard.js';
import { CustomCacheInterceptor } from '../../helpers/customCacheInterceptor.js';
import { JwtAuthGuard } from '../../jwt/jwt.guard.js';
import { ModerationService } from './moderation.service.js';

@Controller('v1/channels/:channelId/moderation')
export class ModerationController {
  constructor(
    @Inject(CACHE_MANAGER) private readonly cacheManager: Cache,
    private readonly moderationService: ModerationService,
  ) { }

  async #delCache(channelId: string) {
    await this.cacheManager.del(`nest:cache:v1/channels/${channelId}/moderation`);
  }

  @CacheTTL(600)
  @UseInterceptors(CustomCacheInterceptor(ctx => {
    const req = ctx.switchToHttp().getRequest() as Request;
    return `nest:cache:v1/channels/${req.params.channelId}/moderation`;
  }))
  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Get()
  root(@Param('channelId') channelId: string) {
    return this.moderationService.getSettings(channelId);
  }

  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Post()
  async update(@Param('channelId') channelId: string, @Body() data: ModerationUpdateDto) {
    const result = await this.moderationService.update(channelId, data);
    await this.#delCache(channelId);
    return result;
  }
}

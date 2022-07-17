import { Body, CacheTTL, CACHE_MANAGER, Controller, Delete, Get, Inject, Param, Post, UseGuards, UseInterceptors } from '@nestjs/common';
import { Cache } from 'cache-manager';
import Express from 'express';

import { DashboardAccessGuard } from '../../../guards/DashboardAccess.guard.js';
import { CustomCacheInterceptor } from '../../../helpers/customCacheInterceptor.js';
import { JwtAuthGuard } from '../../../jwt/jwt.guard.js';
import { DashboardAccessService } from './dashboard-access.service.js';

@Controller('v1/channels/:channelId/settings/dashboardAccess')
export class DashboardAccessController {
  constructor(private readonly service: DashboardAccessService, @Inject(CACHE_MANAGER) private readonly cacheManager: Cache) { }

  async #delCache(channelId: string, userId?: string) {
    await this.cacheManager.del(`nest:cache:v1/channels/${channelId}/settings/dashboardAccess`);
    if (userId) {
      await this.cacheManager.del(`nest:cache:auth/profile:${userId}`);
    }
  }

  @CacheTTL(600)
  @UseInterceptors(CustomCacheInterceptor(ctx => {
    const req = ctx.switchToHttp().getRequest() as Express.Request;
    return `nest:cache:v1/channels/${req.params.channelId}/settings/dashboardAccess`;
  }))
  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Get()
  async root(@Param('channelId') channelId: string) {
    return this.service.getMembers(channelId);
  }

  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Post()
  async addMember(@Param('channelId') channelId: string, @Body() body: { username: string }) {
    const result = await this.service.addMember(channelId, body.username);
    await this.#delCache(channelId);
    return result;
  }

  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Delete(':userId')
  async deleteMember(@Param('channelId') channelId: string, @Param('userId') userId: string) {
    const result = await this.service.deleteMember(channelId, userId);
    await this.#delCache(channelId, userId);
    return result;
  }
}

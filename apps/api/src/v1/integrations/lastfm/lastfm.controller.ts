import { Body, CacheTTL, CACHE_MANAGER, Controller, Get, Inject, Param, Post, UseGuards, UseInterceptors } from '@nestjs/common';
import { Cache } from 'cache-manager';
import Express from 'express';

import { DashboardAccessGuard } from '../../../guards/DashboardAccess.guard.js';
import { CustomCacheInterceptor } from '../../../helpers/customCacheInterceptor.js';
import { JwtAuthGuard } from '../../../jwt/jwt.guard.js';
import { LastfmUpdateDto } from './dto/update.js';
import { LastfmService } from './lastfm.service.js';


@Controller('v1/channels/:channelId/integrations/lastfm')
export class LastfmController {
  constructor(
    private readonly service: LastfmService,
    @Inject(CACHE_MANAGER) private cacheManager: Cache,
  ) { }


  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @CacheTTL(600)
  @UseInterceptors(CustomCacheInterceptor(ctx => {
    const req = ctx.switchToHttp().getRequest() as Express.Request;
    return `nest:cache:v1/channels/${req.params.channelId}/integrations/lastfm`;
  }))
  @Get()
  getIntegration(@Param('channelId') channelId: string) {
    return this.service.getIntegration(channelId);
  }

  @Post()
  async updateIntegration(@Param('channelId') channelId: string, @Body() body: LastfmUpdateDto) {
    const result = await this.service.updateIntegration(channelId, body);
    await this.cacheManager.del(`nest:cache:v1/channels/${channelId}/integrations/lastfm`);
    return result;
  }
}

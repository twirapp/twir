import {
  Body,
  CacheTTL,
  CACHE_MANAGER,
  Controller,
  Get,
  Inject,
  Param,
  Patch,
  Post,
  UseGuards,
  UseInterceptors,
} from '@nestjs/common';
import CacheManager from 'cache-manager';
import Express from 'express';

import { DashboardAccessGuard } from '../../../guards/DashboardAccess.guard.js';
import { CustomCacheInterceptor } from '../../../helpers/customCacheInterceptor.js';
import { JwtAuthGuard } from '../../../jwt/jwt.guard.js';
import { UpdateStreamlabsIntegrationDto } from './dto/patch.js';
import { StreamlabsService } from './streamlabs.service.js';

@Controller('v1/channels/:channelId/integrations/streamlabs')
export class StreamlabsController {
  constructor(
    private readonly service: StreamlabsService,
    @Inject(CACHE_MANAGER) private cacheManager: CacheManager.Cache,
  ) {}

  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @CacheTTL(600)
  @UseInterceptors(
    CustomCacheInterceptor((ctx) => {
      const req = ctx.switchToHttp().getRequest() as Express.Request;
      return `nest:cache:v1/channels/${req.params.channelId}/integrations/streamlabs`;
    }),
  )
  @Get()
  getIntegration(@Param('channelId') channelId: string) {
    return this.service.getIntegration(channelId);
  }

  @Get('auth')
  async auth() {
    const result = await this.service.getAuthLink();

    return result;
  }

  @Patch()
  async updateIntegration(
    @Param('channelId') channelId: string,
    @Body() body: UpdateStreamlabsIntegrationDto,
  ) {
    const result = await this.service.updateIntegration(channelId, body);
    await this.cacheManager.del(
      `nest:cache:v1/channels/${channelId}/integrations/streamlabs/profile`,
    );
    await this.cacheManager.del(`nest:cache:v1/channels/${channelId}/integrations/streamlabs`);
    return result;
  }

  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Post('token')
  async getToken(@Param('channelId') channelId: string, @Body() body: { code: string }) {
    const result = await this.service.getTokens(channelId, body.code);
    await this.cacheManager.del(
      `nest:cache:v1/channels/${channelId}/integrations/streamlabs/profile`,
    );
    await this.cacheManager.del(`nest:cache:v1/channels/${channelId}/integrations/streamlabs`);
    return result;
  }

  @CacheTTL(600)
  @UseInterceptors(
    CustomCacheInterceptor((ctx) => {
      const req = ctx.switchToHttp().getRequest() as Express.Request;
      return `nest:cache:v1/channels/${req.params.channelId}/integrations/streamlabs/profile`;
    }),
  )
  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Get('profile')
  profile(@Param('channelId') channelId: string) {
    return this.service.getProfile(channelId);
  }
}

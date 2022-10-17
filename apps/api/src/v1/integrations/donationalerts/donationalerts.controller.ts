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
import { Throttle } from '@nestjs/throttler';
import CacheManager from 'cache-manager';
import Express from 'express';

import { DashboardAccessGuard } from '../../../guards/DashboardAccess.guard.js';
import { CustomCacheInterceptor } from '../../../helpers/customCacheInterceptor.js';
import { JwtAuthGuard } from '../../../jwt/jwt.guard.js';
import { DonationAlertsService } from './donationalerts.service.js';
import { UpdateDonationAlertsIntegrationDto } from './dto/patch.js';

@Controller('v1/channels/:channelId/integrations/donationalerts')
export class DonationAlertsController {
  constructor(
    private readonly service: DonationAlertsService,
    @Inject(CACHE_MANAGER) private cacheManager: CacheManager.Cache,
  ) {}

  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
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
  @Throttle(1, 5)
  async updateIntegration(
    @Param('channelId') channelId: string,
    @Body() body: UpdateDonationAlertsIntegrationDto,
  ) {
    const result = await this.service.updateIntegration(channelId, body);
    await this.cacheManager.del(
      `nest:cache:v1/channels/${channelId}/integrations/donationalerts/profile`,
    );
    await this.cacheManager.del(`nest:cache:v1/channels/${channelId}/integrations/donationalerts`);
    return result;
  }

  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Post('token')
  async getToken(@Param('channelId') channelId: string, @Body() body: { code: string }) {
    const result = await this.service.getTokens(channelId, body.code);
    await this.cacheManager.del(
      `nest:cache:v1/channels/${channelId}/integrations/donationalerts/profile`,
    );
    await this.cacheManager.del(`nest:cache:v1/channels/${channelId}/integrations/donationalerts`);
    return result;
  }

  @CacheTTL(600)
  @UseInterceptors(
    CustomCacheInterceptor((ctx) => {
      const req = ctx.switchToHttp().getRequest() as Express.Request;
      return `nest:cache:v1/channels/${req.params.channelId}/integrations/donationalerts/profile`;
    }),
  )
  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Get('profile')
  profile(@Param('channelId') channelId: string) {
    return this.service.getProfile(channelId);
  }
}

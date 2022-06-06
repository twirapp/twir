import { Body, CacheTTL, CACHE_MANAGER, Controller, Get, Inject, Param, Patch, Post, Req, Res, UseGuards, UseInterceptors } from '@nestjs/common';
import { Cache } from 'cache-manager';
import { Request, Response } from 'express';

import { DashboardAccessGuard } from '../../../guards/DashboardAccess.guard.js';
import { CustomCacheInterceptor } from '../../../helpers/customCacheInterceptor.js';
import { JwtAuthGuard } from '../../../jwt/jwt.guard.js';
import { UpdateSpotifyIntegrationDto } from './dto/patch.js';
import { SpotifyService } from './spotify.service.js';

@Controller('v1/channels/:channelId/integrations/spotify')
export class SpotifyController {
  constructor(
    private readonly spotifyService: SpotifyService,
    @Inject(CACHE_MANAGER) private cacheManager: Cache,
  ) { }

  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @CacheTTL(600)
  @UseInterceptors(CustomCacheInterceptor(ctx => {
    const req = ctx.switchToHttp().getRequest() as Request;
    return `nest:cache:v1/channels/${req.params.channelId}/integrations/spotify`;
  }))
  @Get()
  getIntegration(@Param('channelId') channelId: string) {
    return this.spotifyService.getIntegration(channelId);
  }

  @Get('auth')
  async auth(@Res() res: Response) {
    const result = await this.spotifyService.getAuthLink();
    res.redirect(result);
  }

  @Patch()
  async updateIntegration(@Param('channelId') channelId: string, @Body() body: UpdateSpotifyIntegrationDto) {
    const result = await this.spotifyService.updateIntegration(channelId, body);
    await this.cacheManager.del(`nest:cache:v1/channels/${channelId}/integrations/spotify/profile`);
    await this.cacheManager.del(`nest:cache:v1/channels/${channelId}/integrations/spotify`);
    return result;
  }

  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Post('token')
  async getToken(@Param('channelId') channelId: string, @Body() body: { code: string }) {
    const result = await this.spotifyService.getTokens(channelId, body.code);
    await this.cacheManager.del(`nest:cache:v1/channels/${channelId}/integrations/spotify/profile`);
    await this.cacheManager.del(`nest:cache:v1/channels/${channelId}/integrations/spotify`);
    return result;
  }

  @CacheTTL(600)
  @UseInterceptors(CustomCacheInterceptor(ctx => {
    const req = ctx.switchToHttp().getRequest() as Request;
    return `nest:cache:v1/channels/${req.params.channelId}/integrations/spotify/profile`;
  }))
  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Get('profile')
  profile(@Param('channelId') channelId: string) {
    return this.spotifyService.getProfile(channelId);
  }
}

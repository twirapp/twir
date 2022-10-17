import {
  Body,
  CacheTTL,
  CACHE_MANAGER,
  Controller,
  Get,
  Inject,
  Param,
  Post,
  UseGuards,
  UseInterceptors,
} from '@nestjs/common';
import CacheManager from 'cache-manager';
import Express from 'express';

import { DashboardAccessGuard } from '../../../guards/DashboardAccess.guard.js';
import { CustomCacheInterceptor } from '../../../helpers/customCacheInterceptor.js';
import { JwtAuthGuard } from '../../../jwt/jwt.guard.js';
import { VkUpdateDto } from './dto/update.js';
import { VkService } from './vk.service.js';

@Controller('v1/channels/:channelId/integrations/vk')
export class VkController {
  constructor(
    private readonly service: VkService,
    @Inject(CACHE_MANAGER) private cacheManager: CacheManager.Cache,
  ) {}

  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Get()
  getIntegration(@Param('channelId') channelId: string) {
    return this.service.getIntegration(channelId);
  }

  @Post()
  async updateIntegration(@Param('channelId') channelId: string, @Body() body: VkUpdateDto) {
    const result = await this.service.updateIntegration(channelId, body);
    await this.cacheManager.del(`nest:cache:v1/channels/${channelId}/integrations/vk`);
    return result;
  }
}

import { Body, Controller, Get, Param, Post, UseGuards } from '@nestjs/common';

import { DashboardAccessGuard } from '../../../guards/DashboardAccess.guard.js';
import { JwtAuthGuard } from '../../../jwt/jwt.guard.js';
import { LastfmUpdateDto } from './dto/update.js';
import { LastfmService } from './lastfm.service.js';

@Controller('v1/channels/:channelId/integrations/lastfm')
export class LastfmController {
  constructor(private readonly service: LastfmService) {}

  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Get()
  getIntegration(@Param('channelId') channelId: string) {
    return this.service.getIntegration(channelId);
  }

  @Post()
  async updateIntegration(@Param('channelId') channelId: string, @Body() body: LastfmUpdateDto) {
    const result = await this.service.updateIntegration(channelId, body);
    return result;
  }
}

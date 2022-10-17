import { Body, Controller, Get, Param, Post, UseGuards } from '@nestjs/common';

import { DashboardAccessGuard } from '../../../guards/DashboardAccess.guard.js';
import { JwtAuthGuard } from '../../../jwt/jwt.guard.js';
import { FaceitUpdateDto } from './dto/update.js';
import { FaceitService } from './faceit.service.js';

@Controller('v1/channels/:channelId/integrations/faceit')
export class FaceitController {
  constructor(private readonly service: FaceitService) {}

  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Get()
  getIntegration(@Param('channelId') channelId: string) {
    return this.service.getIntegration(channelId);
  }

  @Post()
  async updateIntegration(@Param('channelId') channelId: string, @Body() body: FaceitUpdateDto) {
    const result = await this.service.updateIntegration(channelId, body);
    return result;
  }
}

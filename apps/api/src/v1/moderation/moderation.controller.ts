import {
  Body,
  Controller,
  Get,
  Param,
  ParseArrayPipe,
  Post,
  UseGuards,
  UsePipes,
  ValidationPipe,
} from '@nestjs/common';
import { ModerationSettingsDto } from '@tsuwari/shared';

import { DashboardAccessGuard } from '../../guards/DashboardAccess.guard.js';
import { JwtAuthGuard } from '../../jwt/jwt.guard.js';
import { ModerationService } from './moderation.service.js';

@Controller('v1/channels/:channelId/moderation')
export class ModerationController {
  constructor(private readonly moderationService: ModerationService) {}

  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Get()
  root(@Param('channelId') channelId: string) {
    return this.moderationService.getSettings(channelId);
  }

  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @UsePipes(new ValidationPipe({ transform: false }))
  @Post()
  async update(
    @Param('channelId') channelId: string,
    @Body(new ParseArrayPipe({ items: ModerationSettingsDto })) data: ModerationSettingsDto[],
  ) {
    const result = await this.moderationService.update(channelId, data);
    return result;
  }
}

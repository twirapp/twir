import { Body, Controller, Get, Param, Post, Req, UseGuards } from '@nestjs/common';
import { Request } from 'express';

import { DashboardAccessGuard } from '../../guards/DashboardAccess.guard.js';
import { JwtAuthGuard } from '../../jwt/jwt.guard.js';
import { MarkAsReadedDto } from './dto/markAsReaded.dto.js';
import { NotificationsService } from './notifications.service.js';

@Controller('v1/channels/:channelId/notifications')
export class NotificationsController {
  constructor(private readonly service: NotificationsService) {

  }

  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Get('new')
  getNotViewed(@Param('channelId') channelId: string) {
    return this.service.getNotViewed(channelId);
  }

  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Get('viewed')
  getViewed(@Param('channelId') channelId: string) {
    return this.service.getViewed(channelId);
  }

  @UseGuards(JwtAuthGuard, DashboardAccessGuard)
  @Post('viewed')
  markAsReaded(@Param('channelId') channelId: string, @Body() body: MarkAsReadedDto, @Req() req: Request) {
    if (req.user.id !== channelId) return;
    return this.service.markAsRead(channelId, body.notificationId);
  }
}

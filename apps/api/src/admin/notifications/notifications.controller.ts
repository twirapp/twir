import { Body, Controller, Delete, Param, Post, UseGuards } from '@nestjs/common';

import { IsAdminGuard } from '../../guards/IsAdmin.guard.js';
import { JwtAuthGuard } from '../../jwt/jwt.guard.js';
import { CreateNotificationDto } from './dto/create.js';
import { NotificationsService } from './notifications.service.js';

@Controller('admin/notifications')
export class NotificationsController {
  constructor(private readonly service: NotificationsService) { }

  @UseGuards(JwtAuthGuard, IsAdminGuard)
  @Post()
  create(@Body() data: CreateNotificationDto) {
    return this.service.create(data);
  }

  @UseGuards(JwtAuthGuard, IsAdminGuard)
  @Delete(':id')
  del(@Param('id') id: string) {
    return this.service.delete(id);
  }
}

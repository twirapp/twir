import { Module } from '@nestjs/common';
import { RedisModule } from '@tsuwari/shared';

import { NotificationsController } from './notifications.controller.js';
import { NotificationsService } from './notifications.service.js';

@Module({
  imports: [RedisModule],
  controllers: [NotificationsController],
  providers: [NotificationsService],
})
export class NotificationsModule { }

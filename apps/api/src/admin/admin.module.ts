import { Module } from '@nestjs/common';

import { NotificationsModule } from './notifications/notifications.module.js';

@Module({
  imports: [NotificationsModule],
})
export class AdminModule { }

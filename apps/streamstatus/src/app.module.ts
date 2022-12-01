import { Module } from '@nestjs/common';
import { ScheduleModule } from '@nestjs/schedule';
import { RedisModule } from '@tsuwari/shared';

import { AppController } from './app.controller.js';
import { AppService } from './app.service.js';

@Module({
  imports: [RedisModule, ScheduleModule.forRoot()],
  providers: [AppService],
  controllers: [AppController],
})
export class AppModule {}

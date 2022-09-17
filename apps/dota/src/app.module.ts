import { Module } from '@nestjs/common';
import { RedisService } from '@tsuwari/shared';

import { AppController } from './app.controller.js';
import { AppService } from './app.service.js';

@Module({
  imports: [],
  controllers: [AppController],
  providers: [RedisService, AppService],
})
export class AppModule {}

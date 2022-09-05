import { Module } from '@nestjs/common';
import { PrismaModule } from '@tsuwari/prisma';
import { RedisModule } from '@tsuwari/shared';

import { AppController } from './app.controller.js';
import { AppService } from './app.service.js';

@Module({
  imports: [PrismaModule, RedisModule],
  providers: [AppService],
  controllers: [AppController],
})
export class AppModule {}

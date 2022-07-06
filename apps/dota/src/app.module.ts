import { Module } from '@nestjs/common';
import { PrismaClient, PrismaModule } from '@tsuwari/prisma';
import { RedisService } from '@tsuwari/shared';
import { AppController } from './app.controller.js';

import { AppService } from './app.service.js';

@Module({
  imports: [PrismaModule],
  controllers: [AppController],
  providers: [PrismaClient, RedisService, AppService],
})
export class AppModule { }

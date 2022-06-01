import { Module } from '@nestjs/common';
import { PrismaModule, PrismaService } from '@tsuwari/prisma';

import { AppController } from './app.controller.js';
import { AppService } from './app.service.js';
import { RedisService } from './redis.service.js';


@Module({
  imports: [PrismaModule],
  providers: [RedisService, PrismaService, AppService],
  controllers: [AppController],
})
export class AppModule { }
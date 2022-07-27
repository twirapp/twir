import { Module } from '@nestjs/common';
import { PrismaModule } from '@tsuwari/prisma';
import { RedisORMModule } from '@tsuwari/redis';

import { AppController } from './app.controller.js';
import { AppService } from './app.service.js';


@Module({
  imports: [PrismaModule, RedisORMModule],
  providers: [AppService],
  controllers: [AppController],
})
export class AppModule { }
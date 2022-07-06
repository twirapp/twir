import { Module } from '@nestjs/common';
import { PrismaModule } from '@tsuwari/prisma';
import { RedisService } from '@tsuwari/shared';

import { AppService } from './app.service.js';

@Module({
  imports: [PrismaModule],
  controllers: [],
  providers: [RedisService, AppService],
})
export class AppModule { }

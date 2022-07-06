import { Module } from '@nestjs/common';
import { ScheduleModule } from '@nestjs/schedule';
import { PrismaModule } from '@tsuwari/prisma';

import { DotaModule } from './dota/dota.module.js';
import { MicroservicesModule } from './microservices/microservices.module.js';
import { RedisService } from './redis.service.js';
import { StreamStatusModule } from './streamstatus/streamstatus.module.js';

@Module({
  imports: [
    PrismaModule,
    ScheduleModule.forRoot(),
    StreamStatusModule,
    MicroservicesModule,
    DotaModule,
    // IncreaseWatchedModule,
  ],
  providers: [RedisService],
})
export class AppModule { }

import { Module } from '@nestjs/common';
import { ScheduleModule } from '@nestjs/schedule';
import { PrismaModule } from '@tsuwari/prisma';
import { RedisORMModule } from '@tsuwari/redis';
import { RedisService, RedisModule, TwitchApiService } from '@tsuwari/shared';

import { DefaultCommandsCreatorModule } from './default-commands-creator/default-commands-creator.module.js';
import { DotaModule } from './dota/dota.module.js';
import { MicroservicesModule } from './microservices/microservices.module.js';
import { OnlineUsersModule } from './online-users/online-users.module.js';
import { StreamStatusModule } from './streamstatus/streamstatus.module.js';

@Module({
  imports: [
    RedisORMModule,
    PrismaModule,
    RedisModule,
    ScheduleModule.forRoot(),
    StreamStatusModule,
    MicroservicesModule,
    DotaModule,
    DefaultCommandsCreatorModule,
    OnlineUsersModule,
  ],
  providers: [TwitchApiService, RedisService],
})
export class AppModule { }

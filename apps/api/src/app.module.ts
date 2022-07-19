import { CacheModule, CacheModuleOptions, Module } from '@nestjs/common';
import { ThrottlerModule } from '@nestjs/throttler';
import { config } from '@tsuwari/config';
import { PrismaModule } from '@tsuwari/prisma';
import { TwitchApiService } from '@tsuwari/shared';
import cacheRedisStore from 'cache-manager-ioredis';
import Redis from 'ioredis';

import { AdminModule } from './admin/admin.module.js';
import { AppController } from './app.controller.js';
import { AppService } from './app.service.js';
import { AuthModule } from './auth/auth.module.js';
import { JwtAuthModule } from './jwt/jwt.module.js';
import { BotsMicroserviceModule } from './microservices/bots/bots.module.js';
import { RedisModule } from './redis.module.js';
import { RedisService } from './redis.service.js';
import { SocketModule } from './socket/socket.module.js';
import { V1Module } from './v1/v1.module.js';
import { VersionModule } from './version/version.module.js';


const redis = new class extends Redis {
  constructor() {
    super(config.REDIS_URL, {
      isCacheableValue: () => true,
    } as any);
  }
};

@Module({
  imports: [
    CacheModule.register<CacheModuleOptions>({
      store: cacheRedisStore,
      redisInstance: redis,
      isCacheableValue: () => true,
      isGlobal: true,
      ttl: 60,
    }),
    PrismaModule,
    RedisModule,
    BotsMicroserviceModule,
    AuthModule,
    ThrottlerModule.forRoot({
      ttl: 60,
      limit: 100,
    }),
    JwtAuthModule,
    V1Module,
    SocketModule,
    AdminModule,
    VersionModule,
  ],
  controllers: [AppController],
  providers: [
    RedisModule,
    RedisService,
    PrismaModule,
    AppService,
    TwitchApiService,
  ],
})
export class AppModule { }

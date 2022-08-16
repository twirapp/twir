import { CacheModule, CacheModuleOptions, Module, Scope } from '@nestjs/common';
import { APP_GUARD } from '@nestjs/core';
import { ThrottlerModule } from '@nestjs/throttler';
import { config } from '@tsuwari/config';
import { PrismaModule } from '@tsuwari/prisma';
import { RedisORMModule, RedisORMService } from '@tsuwari/redis';
import { RedisModule, TwitchApiService } from '@tsuwari/shared';
import cacheRedisStore from 'cache-manager-ioredis';
import Redis from 'ioredis';

import { AdminModule } from './admin/admin.module.js';
import { AppController } from './app.controller.js';
import { AppService } from './app.service.js';
import { AuthModule } from './auth/auth.module.js';
import { ThrottlerGuard } from './guards/Throttle.guard.js';
import { JwtAuthModule } from './jwt/jwt.module.js';
import { BotsMicroserviceModule } from './microservices/bots/bots.module.js';
import { V1Module } from './v1/v1.module.js';
import { VersionModule } from './version/version.module.js';

const redis = new (class extends Redis {
  constructor() {
    super(config.REDIS_URL, {
      isCacheableValue: () => true,
    } as any);
  }
})();

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
    RedisORMModule,
    BotsMicroserviceModule,
    AuthModule,
    ThrottlerModule.forRoot({
      ttl: 60,
      limit: 100,
    }),
    JwtAuthModule,
    V1Module,
    AdminModule,
    VersionModule,
  ],
  controllers: [AppController],
  providers: [
    PrismaModule,
    RedisORMService,
    AppService,
    TwitchApiService,
    {
      provide: APP_GUARD,
      useClass: ThrottlerGuard,
      scope: Scope.DEFAULT,
    },
  ],
})
export class AppModule {}

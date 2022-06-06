import { CacheModule, CacheModuleOptions, Module } from '@nestjs/common';
import { ThrottlerModule } from '@nestjs/throttler';
import { config } from '@tsuwari/config';
import { PrismaModule } from '@tsuwari/prisma';
import cacheRedisStore from 'cache-manager-ioredis';
import Redis, { RedisOptions } from 'ioredis';

import { AppController } from './app.controller.js';
import { AppService } from './app.service.js';
import { AuthModule } from './auth/auth.module.js';
import { JwtAuthModule } from './jwt/jwt.module.js';
import { BotsMicroserviceModule } from './microservices/bots/bots.module.js';
import { RedisModule } from './redis.module.js';
import { RedisService } from './redis.service.js';
import { SocketModule } from './socket/socket.module.js';
import { V1Module } from './v1/v1.module.js';


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
      limit: 60,
    }),
    JwtAuthModule,
    V1Module,
    SocketModule,
  ],
  controllers: [AppController],
  providers: [
    RedisModule,
    RedisService,
    PrismaModule,
    AppService,
  ],
})
export class AppModule { }

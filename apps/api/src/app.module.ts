import { CacheModule, Module } from '@nestjs/common';
import { ThrottlerModule } from '@nestjs/throttler';
import { config } from '@tsuwari/config';
import { PrismaModule, PrismaService } from '@tsuwari/prisma';
import cacheRedisStore from 'cache-manager-ioredis';
import Redis, { RedisOptions } from 'ioredis';

import { AppController } from './app.controller.js';
import { AppService } from './app.service.js';
import { AuthModule } from './auth/auth.module.js';
import { JwtAuthModule } from './jwt/jwt.module.js';
import { BotsMicroserviceModule } from './microservices/bots/bots.module.js';
import { RedisService } from './redis.service.js';
import { SocketModule } from './socket/socket.module.js';
import { V1Module } from './v1/v1.module.js';

export const redis = new Redis(config.REDIS_URL);


@Module({
  imports: [
    CacheModule.register<RedisOptions>({
      store: cacheRedisStore,
      redisInstance: redis,
      isGlobal: true,
      ttl: 60,
    } as any),
    PrismaModule,
    RedisService,
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
    PrismaModule,
    PrismaService,
    AppService,
  ],
})
export class AppModule { }

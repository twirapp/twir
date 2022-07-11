import { ValidationPipe } from '@nestjs/common';
import { NestFactory } from '@nestjs/core';
import { NatsOptions, Transport } from '@nestjs/microservices';
import * as Sentry from '@sentry/node';
import '@sentry/tracing';
import { config } from '@tsuwari/config';
import cookieParser from 'cookie-parser';
import { Express } from 'express';

import { AppModule } from './app.module.js';

Sentry.init({
  dsn: 'https://1c78d79f3bcb443680e4d5550005e3ac@o324161.ingest.sentry.io/6485379',
  tracesSampleRate: 1.0,
});

export async function initHttp() {
  const app = await NestFactory.create(AppModule);
  app.connectMicroservice({
    transport: Transport.REDIS,
    options: {
      url: config.REDIS_URL,
    },
  });
  app.connectMicroservice<NatsOptions>({
    transport: Transport.NATS,
    options: {
      servers: [config.NATS_URL],
      timeout: 100,
    },
  });

  const adapter = app.getHttpAdapter() as unknown as Express;
  adapter.disable('x-powered-by');
  adapter.disable('etag');

  app.use(cookieParser());
  app.useGlobalPipes(
    new ValidationPipe({
      transform: true,
      whitelist: true,
    }),
  );

  await app.startAllMicroservices();
  await app.listen(3002);
}

initHttp();
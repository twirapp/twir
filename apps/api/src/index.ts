import { ValidationPipe } from '@nestjs/common';
import { NestFactory } from '@nestjs/core';
import { Transport } from '@nestjs/microservices';
import { config } from '@tsuwari/config';
import cookieParser from 'cookie-parser';
import { Express } from 'express';

import { AppModule } from './app.module.js';

export async function initHttp() {
  const app = await NestFactory.create(AppModule);
  app.connectMicroservice({
    transport: Transport.REDIS,
    options: {
      url: config.REDIS_URL,
    },
  });

  const adapter = app.getHttpAdapter() as unknown as Express;
  adapter.disable('x-powered-by');

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
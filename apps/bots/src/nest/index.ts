

import { INestApplication } from '@nestjs/common';
import { NestFactory } from '@nestjs/core';
import { Transport } from '@nestjs/microservices';
import { config } from '@tsuwari/config';

import { AppModule } from './app.module.js';

export const nestApp: INestApplication = await NestFactory.create(AppModule);
nestApp.connectMicroservice({
  transport: Transport.NATS,
  options: {
    servers: [config.NATS_URL],
    timeout: 100,
  },
});

export const startNest = async () => {
  await nestApp.startAllMicroservices();
  await nestApp.listen(3001);
};
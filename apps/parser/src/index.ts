import 'reflect-metadata';

import { INestApplication } from '@nestjs/common';
import { NestFactory } from '@nestjs/core';
import { Transport } from '@nestjs/microservices';
import { config } from '@tsuwari/config';

import { AppModule } from './app.module.js';

export const app: INestApplication = await NestFactory.create(AppModule);

app.connectMicroservice({
  transport: Transport.NATS,
  options: {
    servers: [config.NATS_URL],
    queue: 'parser_queue',
    timeout: 100,
  },
});

await app.startAllMicroservices();
await app.listen(3004);

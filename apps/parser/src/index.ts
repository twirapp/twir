import 'reflect-metadata';

import { NestFactory } from '@nestjs/core';
import { MicroserviceOptions, Transport } from '@nestjs/microservices';
import { config } from '@tsuwari/config';

import { AppModule } from './app.module.js';

export const app = await NestFactory.createMicroservice<MicroserviceOptions>(AppModule, {
  transport: Transport.NATS,
  options: {
    servers: [config.NATS_URL],
  },
});

await app.listen();


import { INestMicroservice } from '@nestjs/common';
import { NestFactory } from '@nestjs/core';
import { MicroserviceOptions, Transport } from '@nestjs/microservices';
import { config } from '@tsuwari/config';

import { AppModule } from './app.module.js';

export const nestApp: INestMicroservice = await NestFactory.createMicroservice<MicroserviceOptions>(AppModule, {
  transport: Transport.NATS,
  options: {
    servers: [config.NATS_URL],
  },
});

export const startNest = async () => await nestApp.listen();
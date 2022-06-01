

import { NestFactory } from '@nestjs/core';
import { MicroserviceOptions, Transport } from '@nestjs/microservices';
import { config } from '@tsuwari/config';
import { resolveProtoPath } from '@tsuwari/grpc';

import { AppModule } from './app.module.js';

export const startNest = async () => {
  const app = await NestFactory.createMicroservice<MicroserviceOptions>(AppModule, {
    transport: Transport.GRPC,
    options: {
      package: 'Bots',
      protoPath: resolveProtoPath('bots'),
      url: config.MICROSERVICE_BOTS_URL,
    },
  });

  await app.listen();
};
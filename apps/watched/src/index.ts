/* import 'reflect-metadata';

import { NestFactory } from '@nestjs/core';
import { MicroserviceOptions, Transport } from '@nestjs/microservices';
import { config } from '@tsuwari/config';

import '@tsuwari/config';
import { AppModule } from './app.module.js';

const app = await NestFactory.createMicroservice<MicroserviceOptions>(AppModule, {
  transport: Transport.GRPC,
  options: {
    package: 'Watched',
    url: config.MICROSERVICE_WATCHED_URL,
    protoPath: resolveProtoPath('watched'),
  },
});

await app.listen(); */
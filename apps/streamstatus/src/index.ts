import 'reflect-metadata';

import { NestFactory } from '@nestjs/core';
import { MicroserviceOptions, Transport } from '@nestjs/microservices';
import { config } from '@tsuwari/config';
import { resolveProtoPath } from '@tsuwari/grpc';

import '@tsuwari/config';
import { AppModule } from './app.module.js';



const app = await NestFactory.createMicroservice<MicroserviceOptions>(AppModule, {
  transport: Transport.GRPC,
  options: {
    package: 'StreamStatus',
    url: config.MICROSERVICE_STREAM_STATUS_URL,
    protoPath: resolveProtoPath('streamstatus'),
  },
});

await app.listen();
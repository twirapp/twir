import 'reflect-metadata';

import { NestFactory } from '@nestjs/core';
import { MicroserviceOptions, Transport } from '@nestjs/microservices';
import * as Sentry from '@sentry/node';
import '@sentry/tracing';
import { config } from '@tsuwari/config';

import '@tsuwari/config';
import { AppModule } from './app.module.js';



Sentry.init({
  dsn: 'https://1c78d79f3bcb443680e4d5550005e3ac@o324161.ingest.sentry.io/6485379',
  tracesSampleRate: 1.0,
});

const app = await NestFactory.createMicroservice<MicroserviceOptions>(AppModule, {
  transport: Transport.NATS,
  options: {
    servers: [config.NATS_URL],
  },
});

await app.listen();
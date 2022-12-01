import 'reflect-metadata';

import { NestFactory } from '@nestjs/core';
import * as Sentry from '@sentry/node';
import '@sentry/tracing';

import '@tsuwari/config';
import { AppModule } from './app.module.js';

Sentry.init({
  dsn: 'https://1c78d79f3bcb443680e4d5550005e3ac@o324161.ingest.sentry.io/6485379',
  tracesSampleRate: 1.0,
});

const app = await NestFactory.createApplicationContext(AppModule);

await app.init();

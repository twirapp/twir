import 'reflect-metadata';

import { NestFactory } from '@nestjs/core';
import { MicroserviceOptions, Transport } from '@nestjs/microservices';
import * as Sentry from '@sentry/node';
import '@sentry/tracing';
import { config } from '@tsuwari/config';
import * as Scheduler from '@tsuwari/grpc/generated/scheduler/scheduler';
import { PORTS } from '@tsuwari/grpc/servers/constants';
import { AppDataSource } from '@tsuwari/typeorm';
import { createServer } from 'nice-grpc';

import { AppModule } from './app.module.js';
import { DefaultCommandsCreatorService } from './default-commands-creator/default-commands-creator.service.js';

Sentry.init({
  dsn: 'https://1c78d79f3bcb443680e4d5550005e3ac@o324161.ingest.sentry.io/6485379',
  tracesSampleRate: 1.0,
});

export const typeorm = await AppDataSource.initialize();

export const app = await NestFactory.createApplicationContext(AppModule);

await app.init();

const schedulerService: Scheduler.SchedulerServiceImplementation = {
  async createDefaultCommands(request: Scheduler.CreateDefaultCommandsRequest) {
    const service = await app.resolve(DefaultCommandsCreatorService);
    service.createDefaultCommands([request.userId]);
    return {};
  },
};

const server = createServer({
  'grpc.keepalive_time_ms': 1 * 60 * 1000,
});

server.add(Scheduler.SchedulerDefinition, schedulerService);

await server.listen(`0.0.0.0:${PORTS.SCHEDULER_SERVER_PORT}`);

process.on('unhandledRejection', (e) => {
  console.error(e);
});
process.on('uncaughtException', (e) => {
  console.error(e);
});

import 'reflect-metadata';

import { NestFactory } from '@nestjs/core';
import { NatsOptions, Transport } from '@nestjs/microservices';
import { ExpressAdapter } from '@nestjs/platform-express';
import { config } from '@tsuwari/config';
import { AppDataSource } from '@tsuwari/typeorm';
import Express from 'express';

import { AppModule } from './app.module.js';
import { EventSub } from './eventsub/eventsub.service.js';

export const typeorm = await AppDataSource.initialize();

const e = Express();
const app = await NestFactory.create(AppModule, new ExpressAdapter(e), { bodyParser: false });
app.connectMicroservice<NatsOptions>({
  transport: Transport.NATS,
  options: {
    servers: [config.NATS_URL],
    timeout: 100,
  },
});

const eventSub = await app.resolve(EventSub);
await eventSub.apply(e);

await app.startAllMicroservices();
e.listen(3003, async () => {
  await eventSub.markAsReady();
  await eventSub.init();
});

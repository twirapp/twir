import 'reflect-metadata';

import { NestFactory } from '@nestjs/core';
import { ExpressAdapter } from '@nestjs/platform-express';
import Express from 'express';

import { AppModule } from './app.module.js';
import { EventSub } from './eventsub/eventsub.service.js';

const e = Express();
const app = await NestFactory.create(AppModule, new ExpressAdapter(e), { bodyParser: false });

const eventSub = await app.resolve(EventSub);
await eventSub.apply(e);

e.listen(3003, async () => {
  await eventSub.markAsReady();
  await eventSub.init();
});

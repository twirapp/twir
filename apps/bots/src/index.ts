import 'reflect-metadata';

import * as Sentry from '@sentry/node';
import '@sentry/tracing';

import { Bots } from './bots.js';
import { prisma } from './libs/prisma.js';
import { initTimers } from './libs/timers.js';
import { startNest } from './nest/index.js';

Sentry.init({
  dsn: 'https://1c78d79f3bcb443680e4d5550005e3ac@o324161.ingest.sentry.io/6485379',
  tracesSampleRate: 1.0,
});

console.info('Starting application...');
await prisma.$connect();
await Bots.init();
await startNest();
initTimers();

/* process.on('unhandledRejection', (r) => console.log(r)); */
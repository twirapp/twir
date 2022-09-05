import 'reflect-metadata';

import * as Sentry from '@sentry/node';
import '@sentry/tracing';

import './libs/nats.js';
import { prisma } from './libs/prisma.js';
import { startNest } from './nest/index.js';

Sentry.init({
  dsn: 'https://1c78d79f3bcb443680e4d5550005e3ac@o324161.ingest.sentry.io/6485379',
  tracesSampleRate: 1.0,
});

await startNest();
await prisma.$connect();
await import('./bots.js').then((b) => b.Bots.init());

/* process.on('unhandledRejection', (r) => console.log(r)); */

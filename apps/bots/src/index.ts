import 'reflect-metadata';

import { Bots } from './bots.js';
import { prisma } from './libs/prisma.js';
import { initTimers } from './libs/timers.js';
import { startNest } from './nest/index.js';

console.info('Starting application...');
await prisma.$connect();
await Bots.init();
await startNest();
initTimers();

/* process.on('unhandledRejection', (r) => console.log(r)); */
import { NestFactory } from '@nestjs/core';

import { AppModule } from './app.module.js';

const app = await NestFactory.createApplicationContext(AppModule);

await app.init();

process.on('unhandledRejection', (e) => {
  console.error(e);
});
process.on('uncaughtException', (e) => {
  console.error(e);
});
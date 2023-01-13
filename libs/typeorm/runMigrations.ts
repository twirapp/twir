import { resolve } from 'path';

import * as dotenv from 'dotenv';

import { AppDataSource } from './src';

dotenv.config({ path: resolve(process.cwd(), '../../.env') });

async function bootstrap() {
  await AppDataSource.initialize();
  await AppDataSource.runMigrations();
}

bootstrap();

import { resolve } from 'node:path';

import dotenv from 'dotenv';

import { AppDataSource } from './src/index.js';

dotenv.config({ path: resolve(process.cwd(), '../../.env') });

if (!process.env.DATABASE_URL) {
  console.error('🚨 Missed database url env.');
  process.exit(1);
}

await AppDataSource.initialize();
await AppDataSource.runMigrations();
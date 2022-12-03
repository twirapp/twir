import { AppDataSource } from './src/index.js';

await AppDataSource.initialize();
await AppDataSource.runMigrations();
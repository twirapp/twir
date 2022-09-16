import { config } from '@tsuwari/config';
import { DataSource } from 'typeorm';

export const AppDataSource = new DataSource({
  type: 'postgres',
  url: config.DATABASE_URL,
  logging: config.isDev,
  entities: ['src/entities/*.ts'],
  subscribers: [],
  migrations: ['src/migrations/*.ts'],
  migrationsTableName: 'typeorm_migrations',
});

import 'reflect-metadata';
import { config } from '@tsuwari/config';
import { DataSource } from 'typeorm';

export * from './entities';
import * as entities from './entities';

export const options = {
  url: config.DATABASE_URL,
  synchronize: config.isDev,
  logging: false,
  entities: Object.keys(entities),
  migrations: ['./src/migrations/*.{js,ts}'],
  subscribers: [],
};

export const AppDataSource = new DataSource({
  type: 'postgres',
  ...options,
});

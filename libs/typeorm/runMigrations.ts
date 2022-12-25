import { AppDataSource } from './src';

async function bootstrap() {
  await AppDataSource.initialize();
  await AppDataSource.runMigrations();
}

bootstrap();

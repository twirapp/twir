import { resolve } from 'path';

import * as dotenv from 'dotenv';

import { AppDataSource } from './src';

dotenv.config({ path: resolve(process.cwd(), '../../.env') });

async function bootstrap() {
	try {
		await AppDataSource.initialize();
		await AppDataSource.runMigrations();
	} catch (e) {
		console.error(e);
		process.exit(1);
	}
}

bootstrap();

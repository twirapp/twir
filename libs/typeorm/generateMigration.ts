import { execSync } from 'node:child_process';
import { resolve } from 'node:path';
import { parseArgs } from 'node:util';

import dotenv from 'dotenv';

dotenv.config({ path: resolve(process.cwd(), '../../.env') });

const {
  values: { name },
} = parseArgs({
  options: {
    name: {
      type: 'string',
      short: 'n',
    },
  },
});

if (!name) {
  console.error('🚨 Name not provided.');
}

if (!process.env.DATABASE_URL) {
  console.error('🚨 Database url not provded');
  process.exit(1);
}

execSync(
  `DATABASE_URL=${process.env.DATABASE_URL} pnpm typeorm-ts-node-esm -d ./src/index.ts migration:generate ./src/migrations/${name}`,
);

console.info(`✅ Migration with name ${name} created`);

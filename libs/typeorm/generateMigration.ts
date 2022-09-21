import { execSync } from 'node:child_process';
import { parseArgs } from 'node:util';

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

const { DATABASE_URL } = process.env;

const localDbUrl = DATABASE_URL?.replace('@postgres', '@localhost');

execSync(
  `DATABASE_URL=${localDbUrl} pnpm typeorm-ts-node-esm -d ./src/index.ts migration:generate ./src/migrations/${name}`,
);

console.info(`✅ Migration with name ${name} created`);

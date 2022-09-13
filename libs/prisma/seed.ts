import { parseArgs } from 'node:util';
import { PrismaClient } from './src/index.js';

const prisma = new PrismaClient();

const cliOptions = {
  id: {
    type: 'string' as const,
    short: 'i',
  },
  accessToken: {
    type: 'string' as const,
    short: 'a',
  },
  refreshToken: {
    type: 'string' as const,
    short: 'r',
  },
};

async function main() {
  const isExists = await prisma.bot.findFirst({
    where: { type: 'DEFAULT' },
  });

  if (isExists) {
    console.info('❌ Bot already exists.');
    return;
  }

  const { values } = parseArgs({ options: cliOptions });

  if (!values.accessToken || !values.refreshToken || !values.id) {
    console.error('❌ Missed params');
    process.exit(1);
  }

  await prisma.bot.create({
    data: {
      id: values.id,
      type: 'DEFAULT',
      token: {
        create: {
          accessToken: values.accessToken,
          refreshToken: values.refreshToken,
          expiresIn: 1000,
          obtainmentTimestamp: new Date(),
        },
      },
    },
  });

  console.info('✅ Done.');
}

main()
  .then(async () => {
    await prisma.$disconnect();
  })
  .catch(async (e) => {
    console.error(e);
    await prisma.$disconnect();
    process.exit(1);
  });

import { streamSchema, usersStatsSchema } from '@tsuwari/redis';

import { USERS_STATUS_CACHE_TTL } from '../constants.js';
import { prisma } from '../libs/prisma.js';
import { redisOm } from '../libs/redis.js';


const repository = redisOm.fetchRepository(usersStatsSchema);
const streamRepository = redisOm.fetchRepository(streamSchema);

export async function increaseUserMessages(userId: string, channelId: string) {
  const stream = await streamRepository.fetch(channelId);
  if (!Object.keys(stream.toRedisJson()).length) return;

  const stats = await prisma.userStats.upsert({
    where: {
      userId_channelId: {
        userId,
        channelId,
      },
    },
    create: {
      user: {
        connectOrCreate: {
          where: {
            id: userId,
          },
          create: {
            id: userId,
          },
        },
      },
      channel: {
        connect: {
          id: channelId,
        },
      },
    },
    update: {
      messages: {
        increment: 1,
      },
    },
  });

  const key = `${channelId}:${userId}`;
  await repository.createAndSave({
    ...stats,
    watched: stats.watched.toString(),
  }, key);
  await repository.expire(key, USERS_STATUS_CACHE_TTL);
}

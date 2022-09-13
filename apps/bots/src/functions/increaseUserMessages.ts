import { usersStatsSchema } from '@tsuwari/redis';

import { USERS_STATUS_CACHE_TTL } from '../constants.js';
import { prisma } from '../libs/prisma.js';
import { redis, redisOm } from '../libs/redis.js';

const repository = redisOm.fetchRepository(usersStatsSchema);

export async function increaseUserMessages(userId: string, channelId: string) {
  const rawStream = await redis.get(`streams:${channelId}`);
  if (!rawStream) return;

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
  await repository.createAndSave(
    {
      ...stats,
      watched: stats.watched.toString(),
    },
    key,
  );
  await repository.expire(key, USERS_STATUS_CACHE_TTL);
}

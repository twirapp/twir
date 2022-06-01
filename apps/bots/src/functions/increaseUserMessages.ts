import { USERS_STATUS_CACHE_TTL } from '../constants.js';
import { prisma } from '../libs/prisma.js';
import { redis } from '../libs/redis.js';

export async function increaseUserMessages(userId: string, channelId: string) {
  const stream = await redis.get(`streams:${channelId}`);

  if (!stream) return;

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

  const key = `usersStats:${channelId}:${userId}`;

  redis.hset(key, 'messages', stats.messages).then(() => {
    redis.expire(key, USERS_STATUS_CACHE_TTL);
  });
}

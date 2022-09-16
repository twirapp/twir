import { UsersStats, usersStatsSchema } from '@tsuwari/redis';
import { User } from '@tsuwari/typeorm/entities/User';
import { UserStats } from '@tsuwari/typeorm/entities/UserStats';

import { USERS_STATUS_CACHE_TTL } from '../constants.js';
import { redis, redisOm } from '../libs/redis.js';
import { typeorm } from '../libs/typeorm.js';

const repository = redisOm.fetchRepository(usersStatsSchema);

const typeormRepository = typeorm.getRepository(UserStats);
const userRepository = typeorm.getRepository(User);

export async function increaseUserMessages(userId: string, channelId: string) {
  const rawStream = await redis.get(`streams:${channelId}`);
  if (!rawStream) return;

  let stats: UserStats;
  const currentStats = await typeormRepository.findOneBy({
    userId,
    channelId,
  });
  if (currentStats) {
    await typeormRepository.update({ userId, channelId }, { messages: currentStats.messages + 1 });
    currentStats.messages + 1;
    stats = currentStats;
  } else {
    const user = await userRepository.findOneBy({
      id: userId,
    });

    if (!user) await userRepository.save({ id: userId });
    stats = await typeormRepository.save({
      userId,
      channelId,
      messages: 1,
    });
  }

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

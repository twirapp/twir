import { StreamRepository, UsersStatsRepository } from '@tsuwari/redis';
import { User } from '@tsuwari/typeorm/entities/User';
import { UserStats } from '@tsuwari/typeorm/entities/UserStats';

import { USERS_STATUS_CACHE_TTL } from '../constants.js';
import { redisSource } from '../libs/redis.js';
import { typeorm } from '../libs/typeorm.js';

const typeormRepository = typeorm.getRepository(UserStats);
const userRepository = typeorm.getRepository(User);

export async function increaseUserMessages(userId: string, channelId: string) {
  const stream = await redisSource.getRepository(StreamRepository).read(channelId);
  if (!stream) return;

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
  const repository = redisSource.getRepository(UsersStatsRepository);

  await repository.write(
    key,
    {
      ...stats,
      watched: stats.watched.toString(),
    },
    USERS_STATUS_CACHE_TTL,
  );
}

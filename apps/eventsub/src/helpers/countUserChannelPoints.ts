import { User } from '@tsuwari/typeorm/entities/User';
import { UserStats } from '@tsuwari/typeorm/entities/UserStats';

import { typeorm } from '../libs/typeorm.js';

export const countUserChannelPoints = async (channelId: string, userId: string, count: number) => {
  const usersRepository = typeorm.getRepository(User);
  const userStatsRepository = typeorm.getRepository(UserStats);

  const user = await usersRepository.findOne({
    where: {
      id: userId,
      stats: {
        channelId,
        userId,
      },
    },
    relations: {
      stats: true,
    },
  });

  if (!user) {
    const newUser = usersRepository.create({
      id: userId,
    });
    await usersRepository.save(newUser);
    const userStats = userStatsRepository.create({
      userId,
      channelId,
      usedChannelPoints: BigInt(count),
    });
    await userStatsRepository.save(userStats);
  } else if (!user.stats?.length) {
    const userStats = userStatsRepository.create({
      userId,         
      channelId,
      usedChannelPoints: BigInt(count),
    });
    await userStatsRepository.save(userStats);
  } else {
    const currentStats = user.stats.at(0)!;
    await userStatsRepository.update(
      { userId, channelId },
      {
        usedChannelPoints: BigInt(currentStats.usedChannelPoints) + BigInt(count),
      },
    );
  }

};
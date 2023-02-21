import { UserStats } from '@tsuwari/typeorm/entities/UserStats';

import { typeorm } from '../libs/typeorm.js';

export const decrementUserChannelPoints = async (channelId: string, userId: string, count: number) => {
  const repository = typeorm.getRepository(UserStats);

  const stats = await repository.findOneBy({
    channelId,
    userId,
  });

  if (!stats) return;
  if ((BigInt(stats.usedChannelPoints) - BigInt(count)) <= BigInt(0)) return;

  repository.update(
    { userId, channelId },
    { usedChannelPoints: BigInt(stats.usedChannelPoints) - BigInt(count) },
  );
};
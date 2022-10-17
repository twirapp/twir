import { ChannelStream } from '@tsuwari/typeorm/entities/ChannelStream';
import { User } from '@tsuwari/typeorm/entities/User';
import { UserStats } from '@tsuwari/typeorm/entities/UserStats';

import { USERS_STATUS_CACHE_TTL } from '../constants.js';
import { typeorm } from '../libs/typeorm.js';

const statsRepository = typeorm.getRepository(UserStats);
const userRepository = typeorm.getRepository(User);
const streamRepository = typeorm.getRepository(ChannelStream);

export async function increaseUserMessages(userId: string, channelId: string) {
  const stream = await streamRepository.findOneBy({ userId: channelId });
  if (!stream) return;

  const currentStats = await statsRepository.findOneBy({
    userId,
    channelId,
  });
  if (currentStats) {
    await statsRepository.increment({ userId, channelId }, 'messages', 1);
  } else {
    const user = await userRepository.findOneBy({
      id: userId,
    });

    if (!user) await userRepository.save({ id: userId });
    await statsRepository.save({
      userId,
      channelId,
      messages: 1,
    });
  }
}

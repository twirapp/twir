import { ChannelStream } from '@tsuwari/typeorm/entities/ChannelStream';

import { typeorm } from '../libs/typeorm.js';

const repository = typeorm.getRepository(ChannelStream);

export const increaseParsedMessages = async (userId: string) => {
  const stream = await repository.findOne({
    where: { userId },
  });
  if (!stream) return;

  await repository.increment({ userId }, 'parsedMessages', 1);
};

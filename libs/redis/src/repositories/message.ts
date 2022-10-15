import Redis from 'ioredis';

import { BaseRepository } from '../base.js';

type Message = {
  messageId: string;
  userId: string;
  channelId: string;
  userName: string;
  message: string;
  canBeDeleted: boolean;
};

export class MessagesRepository extends BaseRepository<Message> {
  constructor(redis: Redis) {
    super('variables', redis);
  }
}

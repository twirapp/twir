import { Injectable, Logger } from '@nestjs/common';
import { Interval } from '@nestjs/schedule';
import { config } from '@tsuwari/config';
import { ChannelChatMessage } from '@tsuwari/typeorm/entities/ChannelChatMessage';

import { typeorm } from '../index.js';

@Injectable()
export class DeleteOldMessagesService {
  @Interval('oldMessages', config.isDev ? 10000 : 1 * 60 * 60 * 1000)
  async delete() {
    const deleteTime = Date.now() - 6 * 60 * 60 * 1000;

    typeorm
      .getRepository(ChannelChatMessage)
      .createQueryBuilder('msg')
      .delete()
      .where('msg.createdAt <= :deleteTime', { deleteTime })
      .execute();
  }
}

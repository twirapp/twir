import { RedisORMService, messageSchema } from '@tsuwari/redis';
import { ClientProxy } from '@tsuwari/shared';
import { lastValueFrom } from 'rxjs';

import { app } from '../../index.js';
import { DefaultCommand } from '../types.js';

const redisOm = app.get(RedisORMService);
const repository = redisOm.fetchRepository(messageSchema);
const nats = app.get('NATS').providers[0].useValue as ClientProxy;

export const nuke: DefaultCommand[] = [
  {
    name: 'nuke',
    description: 'Mass remove messages in chat by message content. Usage: <b>!nuke phrase</b>',
    permission: 'MODERATOR',
    visible: false,
    module: 'CHANNEL',
    handler: async (state, params) => {
      if (!state.channelId || !params?.length) return;

      const messages = await repository.search()
        .where('channelId').equal(state.channelId)
        .and('message').match(params)
        .return.all();

      if (!messages.length) return;

      const result = await lastValueFrom(nats.send('bots.deleteMessages', {
        channelId: state.channelId,
        channelName: state.target.value.substring(1),
        messageIds: messages.map(m => m.toRedisJson()).filter(m => m.canBeDeleted).map(m => m.messageId),
      }));

      return result ? '✅' : '❌';
    },
  },
];
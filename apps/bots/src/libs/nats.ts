import { config } from '@tsuwari/config';
import * as NatsBots from '@tsuwari/nats/bots';
import { Channel } from '@tsuwari/typeorm/entities/Channel';
import { connect } from 'nats';

import { Bots } from '../bots.js';
import { typeorm } from './typeorm.js';

export const nats = await connect({
  servers: [config.NATS_URL],
});
const sub = nats.subscribe('bots.deleteMessages');
(async () => {
  for await (const m of sub) {
    const data = NatsBots.DeleteMessagesRequest.fromBinary(m.data);
    const channel = await typeorm.getRepository(Channel).findOneBy({
      id: data.channelId,
    });

    if (!channel) continue;

    const bot = Bots.cache.get(channel.botId);
    if (!bot) continue;

    for (const id of data.messageIds) {
      bot.deleteMessage(data.channelName, id);
    }

    continue;
  }
})();

import { config } from '@tsuwari/config';
import * as NatsBots from '@tsuwari/nats/bots';
import { connect } from 'nats';

import { Bots } from '../bots.js';
import { prisma } from './prisma.js';

export const nats = await connect({
  servers: [config.NATS_URL],
});
const sub = nats.subscribe('bots.deleteMessages');
(async () => {
  for await (const m of sub) {
    console.log('recieved');
    const data = NatsBots.DeleteMessagesRequest.fromBinary(m.data);
    console.log(data);
    const channel = await prisma.channel.findFirst({
      where: { id: data.channelId },
    });

    console.log(channel);

    if (!channel) continue;

    const bot = Bots.cache.get(channel?.botId);
    console.log(bot);
    if (!bot) continue;

    for (const id of data.messageIds) {
      bot.deleteMessage(data.channelName, id);
    }

    continue;
  }
})();

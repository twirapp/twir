import { config } from '@tsuwari/config';
import * as NatsBots from '@tsuwari/nats/bots';
import { Channel } from '@tsuwari/typeorm/entities/Channel';
import { connect } from 'nats';

import { Bots } from '../bots.js';
import { typeorm } from './typeorm.js';

export const nats = await connect({
  servers: [config.NATS_URL],
});

async function deleteMessagesQueue() {
  const sub = nats.subscribe('bots.deleteMessages');
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
}

async function sendMessagesQueue() {
  const sub = nats.subscribe('bots.sendMessage');
  for await (const m of sub) {
    const data = NatsBots.SendMessage.fromBinary(m.data);
    const channel = await typeorm.getRepository(Channel).findOneBy({
      id: data.channelId,
    });

    if (!channel) continue;

    const bot = Bots.cache.get(channel.botId);
    if (!bot) continue;

    bot.say(data.channelName, data.message);

    continue;
  }
}

deleteMessagesQueue();
sendMessagesQueue();

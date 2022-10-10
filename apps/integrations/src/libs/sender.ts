import { config } from '@tsuwari/config';
import * as Bots from '@tsuwari/nats/bots';
import { MyRefreshingProvider } from '@tsuwari/shared';
import { Channel } from '@tsuwari/typeorm/entities/Channel';
import { Token } from '@tsuwari/typeorm/entities/Token';
import { ApiClient, HelixChatAnnoucementColor } from '@twurple/api';

import { typeorm } from '../index.js';

export async function sendMessage(opts: {
  channelId: string;
  message: string | null;
  color?: HelixChatAnnoucementColor;
}) {
  const channel = await typeorm.getRepository(Channel).findOne({
    where: { id: opts.channelId },
    relations: {
      bot: { token: true },
    },
  });

  if (!channel || !channel?.bot?.token) return;

  const botApi = new ApiClient({
    authProvider: new MyRefreshingProvider({
      clientId: config.TWITCH_CLIENTID,
      clientSecret: config.TWITCH_CLIENTSECRET,
      repository: typeorm.getRepository(Token),
      token: channel.bot.token,
    }),
  });

  await botApi.chat.sendAnnouncement(opts.channelId, channel.botId, {
    message: !opts.message || opts.message === 'null' ? '' : opts.message,
    color: opts.color,
  });
}

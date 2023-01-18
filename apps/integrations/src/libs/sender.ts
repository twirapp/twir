import { config } from '@tsuwari/config';
import { Channel } from '@tsuwari/typeorm/entities/Channel';
import { ApiClient, HelixChatAnnoucementColor } from '@twurple/api';
import { StaticAuthProvider } from '@twurple/auth';

import { typeorm } from '../index.js';
import { tokensGrpcClient } from './tokensGrpc.js';

export async function sendMessage(opts: {
  channelId: string;
  message: string | null;
  color?: HelixChatAnnoucementColor;
}) {
  const channel = await typeorm.getRepository(Channel).findOne({
    where: { id: opts.channelId },
  });

  if (!channel) return;

  const token = await tokensGrpcClient.requestBotToken({
    botId: channel.botId,
  }).catch(() => null);

  if (!token) return;

  const botApi = new ApiClient({
    authProvider: new StaticAuthProvider(config.TWITCH_CLIENTID, token.accessToken, token.scopes),
  });

  await botApi.chat.sendAnnouncement(opts.channelId, channel.botId, {
    message: !opts.message ? '' : opts.message,
    color: opts.color,
  });
}

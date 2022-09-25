import { ChannelEvent, EventType } from '@tsuwari/typeorm/entities/ChannelEvent';
import { ChannelDonationEvent } from '@tsuwari/typeorm/entities/channelEvents/Donation';
import { ChannelIntegration } from '@tsuwari/typeorm/entities/ChannelIntegration';
import Centrifuge from 'centrifuge';
import WebSocket from 'ws';

import { typeorm } from '../index.js';
import { sendMessage } from '../libs/sender.js';

export class DonationAlerts {
  #socket: Centrifuge;
  constructor(
    private readonly accessToken: string,
    private readonly donationAlertsUserId: string,
    private readonly socketConnectionToken: string,
    private readonly twitchUserId: string,
  ) {}

  async init() {
    this.#socket = new Centrifuge('wss://centrifugo.donationalerts.com/connection/websocket', {
      websocket: WebSocket,
      onPrivateSubscribe: async (ctx, cb) => {
        const request = await fetch('https://www.donationalerts.com/api/v1/centrifuge/subscribe', {
          method: 'POST',
          body: JSON.stringify(ctx.data),
          headers: { Authorization: `Bearer ${this.accessToken}` },
        });

        const response = await request.json();
        if (!request.ok) {
          console.error(response);
          cb({ status: request.status, data: {} as any });
        }

        cb({ status: 200, data: { channels: response.channels } });
      },
    });

    this.#socket.setToken(this.socketConnectionToken);
    this.#socket.connect();

    const channel = this.#socket.subscribe(`$alerts:donation_${this.donationAlertsUserId}`);

    channel.on('publish', async ({ data }: { data: Message }) => {
      const event = await typeorm.getRepository(ChannelEvent).save({
        channelId: this.twitchUserId,
        type: EventType.DONATION,
      });
      await typeorm.getRepository(ChannelDonationEvent).save({
        event,
        amount: data.amount,
        currency: data.currency,
        toUserId: this.twitchUserId,
        message: data.message,
        username: data.username,
      });
      sendMessage({
        channelId: this.twitchUserId,
        message: `${data.username}: ${data.amount}${data.currency} ${data.message}`,
        channelName: '',
      });
    });
  }

  async destroy() {
    this.#socket.disconnect();
  }
}

export type Message = {
  id: number;
  name: string;
  username: string;
  message: string;
  message_type: 'text' | 'audio';
  payin_system: null | any;
  amount: number;
  currency: string;
  amount_in_user_currency: number;
  recipient_name: string;
  recipient: {
    user_id: number;
    code: string;
    name: string;
    avatar: string;
  };
  created_at: string;
  shown_at: null | any;
  reason: string;
};

export async function addDonationAlertsIntegration(integration: ChannelIntegration) {
  if (
    !integration.accessToken ||
    !integration.refreshToken ||
    !integration.integration ||
    !integration.integration.clientId ||
    !integration.integration.clientSecret
  ) {
    return;
  }

  const refresh = await fetch('https://www.donationalerts.com/oauth/token', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded',
    },
    body: new URLSearchParams({
      grant_type: 'refresh_token',
      refresh_token: integration.refreshToken,
      client_id: integration.integration.clientId,
      client_secret: integration.integration.clientSecret,
    }).toString(),
  });

  if (!refresh.ok) {
    console.error(await refresh.text());
    return;
  }

  const refreshResponse = await refresh.json();

  await typeorm
    .getRepository(ChannelIntegration)
    .update(
      { id: integration.id },
      { accessToken: refreshResponse.access_token, refreshToken: refreshResponse.refresh_token },
    );

  const request = await fetch('https://www.donationalerts.com/api/v1/user/oauth', {
    headers: {
      Authorization: `Bearer ${refreshResponse.access_token}`,
    },
  });

  if (!request.ok) {
    console.log(await request.text(), refreshResponse);
    return;
  }

  const { data } = await request.json();
  const { id, socket_connection_token } = data;
  const instance = new DonationAlerts(
    refreshResponse.access_token,
    id,
    socket_connection_token,
    integration.channelId,
  );
  await instance.init();

  return instance;
}

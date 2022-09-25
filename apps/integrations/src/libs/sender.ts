import * as Bots from '@tsuwari/nats/bots';

import { nats } from './nats.js';

export async function sendMessage(opts: Bots.SendMessage) {
  const data = Bots.SendMessage.toBinary(opts);

  nats.publish('bots.sendMessage', data);
}

import { Logger } from '@nestjs/common';
import { config } from '@tsuwari/config';
import { ParseResponseRequest, ParseResponseResponse } from '@tsuwari/nats/parser';
import { Queue } from '@tsuwari/shared';
import { ChannelTimer } from '@tsuwari/typeorm/entities/ChannelTimer';
import { HelixStreamData } from '@twurple/api/lib/index.js';

import { Bots, staticApi } from '../bots.js';
import { nats } from './nats.js';
import { redis } from './redis.js';
import { typeorm } from './typeorm.js';

const logger = new Logger('Timers');
const repository = typeorm.getRepository(ChannelTimer);
export const timersQueue = new Queue<ChannelTimer>(async function (taskId: string) {
  const timer = await repository.findOne({
    where: { id: taskId },
    relations: { channel: true },
  });

  if (!timer || !timer.enabled || !timer.channel) {
    return;
  }

  const rawStream = await redis.get(`streams:${timer.channelId}`);
  if (!rawStream) return;
  const stream = JSON.parse(rawStream) as HelixStreamData & { parsedMessages?: number };

  stream.parsedMessages = stream.parsedMessages ?? 0;

  if (
    timer.messageInterval > 0 &&
    timer.lastTriggerMessageNumber - stream.parsedMessages + timer.messageInterval > 0
  ) {
    return;
  }

  const responses = timer.responses as Array<string>;

  const bot = Bots.cache.get(timer.channel.botId);
  const user = await staticApi.users.getUserById(timer.channelId);

  const response = responses[timer.last];
  if (!response) return;

  if (!bot || !user) {
    return;
  }

  if (bot._authProvider) {
    const data = ParseResponseRequest.toBinary({
      channel: {
        id: timer.channelId,
        name: '',
      },
      message: {
        id: '',
        text: response,
      },
      sender: {
        badges: [],
        displayName: '',
        id: '',
        name: '',
      },
    });
    const request = await nats.request('parser.parseTextResponse', data);
    const recievedResponse = ParseResponseResponse.fromBinary(request.data);

    for (const r of recievedResponse.responses) {
      if (config.isProd) {
        bot.say(user.name, r);
      } else {
        logger.log(`${user.name} -> ${r}`);
      }
    }
  }

  await repository.update(
    { id: timer.id },
    {
      last: ++timer.last % (timer.responses as string[]).length,
      lastTriggerMessageNumber: stream.parsedMessages as number,
    },
  );
});

const getId = (t: ChannelTimer | string) => (typeof t === 'string' ? t : t.id);
export async function addTimerToQueue(timerOrId: ChannelTimer | string) {
  let timer: ChannelTimer | null;

  if (typeof timerOrId === 'string') {
    timer = await repository.findOneBy({ id: timerOrId });
    if (!timer?.enabled) return;
  } else {
    timer = timerOrId as ChannelTimer;
  }

  removeTimerFromQueue(timerOrId);
  if (timer) {
    timersQueue.addTimerToQueue(timer.id, timer, {
      interval: timer.timeInterval * (config.isDev ? 1000 : 60000),
    });
  }
}

export function removeTimerFromQueue(timer: ChannelTimer | string) {
  const id = getId(timer);

  timersQueue.removeTimerFromQueue(id);
}

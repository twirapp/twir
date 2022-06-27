import { config } from '@tsuwari/config';
import { Timer } from '@tsuwari/prisma';
import { Worker, Queue, QueueScheduler } from 'bullmq';
import Redis from 'ioredis';

import { Bots, staticApi } from '../bots.js';
import { nestApp } from '../nest/index.js';
import { ParserService } from '../nest/parser/parser.service.js';
import { prisma } from './prisma.js';
import { redis } from './redis.js';

const createConnection = () => new Redis(config.REDIS_URL, { maxRetriesPerRequest: null });

type Data = { id: string }

new QueueScheduler('timers', {
  connection: createConnection(),
});
export const timersQueue = new Queue<Data>('timers', {
  connection: createConnection(),
});
await timersQueue.drain(true);
await timersQueue.clean(0, Number.MAX_SAFE_INTEGER);

new Worker<Data>(
  'timers',
  async (job) => {
    const timer = await prisma.timer.findFirst({
      where: {
        id: job.data.id,
      },
      include: {
        channel: true,
      },
    });

    if (!timer) {
      job.discard();
      await job.remove();
      return;
    }

    const stream = await redis.get(`streams:${timer.channelId}`);
    if (!stream) return;

    const parsedStream = JSON.parse(stream);

    if (timer.messageInterval > 0 && timer.lastTriggerMessageNumber - parsedStream.parsedMessages + timer.messageInterval > 0) {
      return;
    }

    const responses = timer.responses as Array<string>;

    const bot = Bots.cache.get(timer.channel.botId);
    const user = await staticApi.users.getUserById(timer.channelId);

    const response = responses[timer.last];
    if (!response) return;

    if (!bot || !user) {
      await job.remove();
      throw new Error('Something very unexpected happend');
    }

    if (bot._authProvider) {
      const service = nestApp.get(ParserService);
      const parsedResponses = await service.parseResponse({
        channelId: timer.channelId,
        text: response,
      });

      if (parsedResponses) {
        for (const r of parsedResponses) {
          bot.say(
            user.name,
            r,
          );
        }
      }
    }

    await prisma.timer.update({
      where: {
        id: timer.id,
      },
      data: {
        last: ++timer.last % (timer.responses as string[]).length,
        lastTriggerMessageNumber: parsedStream.parsedMessages as number,
      },
    });
  },
  {
    connection: createConnection(),
  },
);

export async function initTimers() {
  const timers = await prisma.timer.findMany();
  for (const timer of timers.filter(t => t.enabled)) {
    addTimerToQueue(timer);
  }
}

export async function addTimerToQueue(timerOrId: Timer | string) {
  const id = getId(timerOrId);
  let timer: Timer | null;

  if (typeof id === 'string') {
    timer = await prisma.timer.findFirst({ where: { id: id as string } });
    if (!timer?.enabled) return;
  } else {
    timer = timerOrId as Timer;
  }

  if (timer) {
    await timersQueue.add(timer.id, { id: timer.id }, { repeat: { every: timer.timeInterval * 1000 } });
  }
}

const getId = (t: Timer | string) => typeof t === 'string' ? t : t.id;

export async function updateTimer(timer: Timer | string) {
  const id = getId(timer);

  await removeTimerFromQueue(timer);
  if (typeof id === 'string') {
    const entity = await prisma.timer.findFirst({ where: { id } });
    if (entity) await addTimerToQueue(entity);
  } else {
    addTimerToQueue(timer as Timer);
  }
}

export async function removeTimerFromQueue(timer: Timer | string) {
  const id = getId(timer);

  return await timersQueue.remove(id);
}

import { config } from '@tsuwari/config';
import { Client } from '@tsuwari/redis';
import Redis from 'ioredis';
import Redlock from 'redlock';


export const redis = new Redis(config.REDIS_URL);
export const redlock = new Redlock(
  // You should have one client for each independent redis node
  // or cluster.
  [redis as any],
  {
    // The expected clock drift; for more details see:
    // http://redis.io/topics/distlock
    driftFactor: 0.01, // multiplied by lock ttl to determine drift time

    // The max number of times Redlock will attempt to lock a resource
    // before erroring.
    retryCount: 10,

    // the time in ms between attempts
    retryDelay: 200, // time in ms

    // the max time in ms randomly added to retries
    // to improve performance under high contention
    // see https://www.awsarchitectureblog.com/2015/03/backoff.html
    retryJitter: 200, // time in ms

    // The minimum remaining time on a lock before an extension is automatically
    // attempted with the `using` API.
    automaticExtensionThreshold: 500, // time in ms
  },
);
export async function findOrSetObject(key: string, data: Record<any, any>) {
  if (!(await redis.exists(key))) {
    await redis.hmset(key, data);
  } else {
    const entity = await redis.hgetall(key);

    return entity;
  }
}

export async function findOrSetString(key: string, data: string) {
  const entity = await redis.get(key);

  if (entity) return entity;

  await redis.set(key, data);
}

export const redisOm = await new Client().open(config.REDIS_URL);
import Redis from 'ioredis';

export class BaseRepository<T> {
  constructor(private readonly prefix: string, private readonly redis: Redis) {}

  async read(key: string): Promise<T | null> {
    const entity = await this.redis.get(`${this.prefix}:${key}`);
    if (!entity) {
      return null;
    }

    return JSON.parse(entity) as T;
  }

  async write(key: string, data: T, expire = 0) {
    await this.redis.set(`${this.prefix}:${key}`, JSON.stringify(data), 'EX', expire);
  }

  async readMany(keys: string[], rawKey = false): Promise<T[]> {
    const itemsKeys = rawKey ? keys : keys.map((k) => `${this.prefix}:${k}`);
    const items = await this.redis.mget(itemsKeys);
    const result = [] as T[];

    for (const item of items) {
      result.push(JSON.parse(item));
    }

    return result;
  }

  async remove(key: string) {
    await this.redis.del(`${this.prefix}:${key}`);
  }
}

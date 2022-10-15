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
}

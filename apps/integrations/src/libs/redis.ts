import { config } from '@twir/config'
import { RedisClient } from 'bun'

export const client = new RedisClient(config.REDIS_URL)

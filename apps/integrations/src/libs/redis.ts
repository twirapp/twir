import { config } from '@twir/config'
import { createClient } from 'redis'

export const client = await createClient({
	url: config.REDIS_URL,
}).connect()

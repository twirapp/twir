import { dirname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'

import { config, readEnv } from '@twir/config'
import { createClient } from 'redis'

const __dirname = dirname(fileURLToPath(import.meta.url))

readEnv(resolve(__dirname, '..', '..', '..', '.env'))

export const client = await createClient({
	url: config.REDIS_URL,
}).connect()

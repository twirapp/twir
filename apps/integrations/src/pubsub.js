import { dirname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'

import { config, readEnv } from '@twir/config'
import { createPubSub } from '@twir/pubsub'

import { onDonation } from './utils/onDonation.js'

const __dirname = dirname(fileURLToPath(import.meta.url))

readEnv(resolve(__dirname, '..', '..', '..', '.env'))

const pubSub = await createPubSub(config.REDIS_URL)

pubSub.subscribe('donations:new', async (message) => {
	try {
		const data = JSON.parse(message)
		await onDonation(data)
	} catch (e) {
		console.log(message, e)
	}
})

import { config } from '@twir/config'
import { createPubSub } from '@twir/pubsub'

import { onDonation } from './utils/onDonation.js'

const pubSub = await createPubSub(config.REDIS_URL)

pubSub.subscribe('donations:new', async (message) => {
	try {
		const data = JSON.parse(message)
		await onDonation(data)
	} catch (e) {
		console.log(message, e)
	}
})

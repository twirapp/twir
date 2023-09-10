import { config } from '@twir/config';
import { createPubSub } from '@twir/pubsub';

import { Donate, onDonation } from './utils/onDonation.js';

const pubSub = await createPubSub(config.REDIS_URL);

pubSub.subscribe('donations:new', (message) => {
	try {
		const data = JSON.parse(message) as Donate;
		onDonation(data);
	} catch (e) {
		console.log(message, e);
	}
});

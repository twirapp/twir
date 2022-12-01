import { config } from '@tsuwari/config';
import { createPubSub } from '@tsuwari/pubsub';

export const pubsub = await createPubSub(config.REDIS_URL);

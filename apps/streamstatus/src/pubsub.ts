import { config } from '@tsuwari/config';
import { createPubSub } from '@tsuwari/pubsub';

export const pubSub = await createPubSub(config.REDIS_URL);

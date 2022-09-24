import { config } from '@tsuwari/config';
import { connect } from 'nats';

export const nats = await connect({
  servers: [config.NATS_URL],
});

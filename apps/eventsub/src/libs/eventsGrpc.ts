import { config } from '@tsuwari/config';
import { createEvents } from '@tsuwari/grpc/clients/events';

export const eventsGrpcClient = await createEvents(config.NODE_ENV);

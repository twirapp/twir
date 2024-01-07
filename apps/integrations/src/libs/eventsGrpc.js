import { config } from '@twir/config';
import { createEvents } from '@twir/grpc/clients/events';

export const eventsGrpcClient = await createEvents(config.NODE_ENV);

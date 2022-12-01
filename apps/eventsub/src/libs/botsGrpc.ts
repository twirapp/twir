import { config } from '@tsuwari/config';
import { createBots } from '@tsuwari/grpc/clients/bots';

export const botsGrpcClient = createBots(config.NODE_ENV);

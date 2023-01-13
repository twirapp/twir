import { config } from '@tsuwari/config';
import { createBots } from '@tsuwari/grpc/clients/bots';

export const botsGrpcClient = await createBots(config.NODE_ENV);

import { config } from '@twir/config';
import { createBots } from '@twir/grpc/clients/bots';

export const botsGrpcClient = await createBots(config.NODE_ENV);


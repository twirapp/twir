import { config } from '@tsuwari/config';
import { createWatched } from '@tsuwari/grpc/clients/watched';

export const watchedGrpcClient = createWatched(config.NODE_ENV);

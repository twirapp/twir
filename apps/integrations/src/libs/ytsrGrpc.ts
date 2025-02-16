import { config } from '@twir/config';
import { createYtsr } from '@twir/grpc/clients/ytsr';

export const ytsrGrpcClient = await createYtsr(config.NODE_ENV);

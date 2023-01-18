import { config } from '@tsuwari/config';
import { createTokens } from '@tsuwari/grpc/clients/tokens';

export const tokensGrpcClient = await createTokens(config.NODE_ENV);


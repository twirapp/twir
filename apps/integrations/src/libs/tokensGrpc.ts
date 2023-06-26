import { config } from '@twir/config';
import { createTokens } from '@twir/grpc/clients/tokens';

export const tokensGrpcClient = await createTokens(config.NODE_ENV);


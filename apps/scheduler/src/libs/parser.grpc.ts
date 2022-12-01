import { config } from '@tsuwari/config';
import { createParser } from '@tsuwari/grpc/clients/parser';

export const parserGrpcClient = createParser(config.NODE_ENV);

import { config } from '@twir/config'
import { createParser } from '@twir/grpc/clients/parser'

export const parserGrpcClient = await createParser(config.NODE_ENV)

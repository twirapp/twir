import { config } from '@twir/config'

const host = config.SITE_BASE_URL ?? 'localhost:3005'
const isDev = config.NODE_ENV === 'development'
const apiAddr = isDev ? `${host}/api` : 'api-gql:3009'
const protocol = config.USE_WSS ? 'https' : 'http'

export const gqlUrl = process.env.API_GQL_ADDR || `${protocol}://${apiAddr}/query`

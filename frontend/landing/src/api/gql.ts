import { config } from '@twir/config'

const host = config.SITE_BASE_URL ?? 'localhost:3005'
const isDev = config.NODE_ENV === 'development'
const apiAddr = isDev ? `${host}/api` : 'api-gql:3009'
const protocol = isDev && !config.USE_WSS ? 'http' : 'https'

export const gqlUrl = process.env.API_GQL_ADDR || `${protocol}://${apiAddr}/query`

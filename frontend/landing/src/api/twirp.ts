import { TwirpFetchTransport } from '@protobuf-ts/twirp-transport'
import { ProtectedClient, UnProtectedClient } from '@twir/api/api.client'
import { config } from '@twir/config'

const host = config.SITE_BASE_URL ?? 'localhost:3005'
const isDev = config.NODE_ENV === 'development'
const apiAddr = isDev ? `${host}/api-old` : 'api:3002'

const baseUrl = `http://${apiAddr}/v1`

const transport = new TwirpFetchTransport({
	baseUrl,
	sendJson: isDev
})

export const protectedClient = new ProtectedClient(transport)
export const unProtectedClient = new UnProtectedClient(transport)

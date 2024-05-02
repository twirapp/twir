import { TwirpFetchTransport } from '@protobuf-ts/twirp-transport'
import { ProtectedClient, UnProtectedClient } from '@twir/api/api.client'

const transport = new TwirpFetchTransport({
	baseUrl: `${location.origin}/api-old/v1`
})

export const browserUnProtectedClient = new UnProtectedClient(transport)
export const browserProtectedClient = new ProtectedClient(transport)

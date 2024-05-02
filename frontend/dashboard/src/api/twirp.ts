import { TwirpFetchTransport } from '@protobuf-ts/twirp-transport'
import { AdminClient, ProtectedClient, UnProtectedClient } from '@twir/api/api.client'

const transport = new TwirpFetchTransport({
	baseUrl: `${window.location.origin}/api-old/v1`,
	sendJson: import.meta.env.DEV
})
export const protectedApiClient = new ProtectedClient(transport)
export const unprotectedApiClient = new UnProtectedClient(transport)
export const adminApiClient = new AdminClient(transport)

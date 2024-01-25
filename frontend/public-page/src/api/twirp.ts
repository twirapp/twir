import { TwirpFetchTransport } from '@protobuf-ts/twirp-transport';
import { ProtectedClient, UnProtectedClient } from '@twir/api/api.client';

const transport = new TwirpFetchTransport({
	baseUrl: `${window.location.origin}/api/v1`,
	sendJson: import.meta.env.DEV,
});

export const unprotectedClient = new UnProtectedClient(transport);
export const protectedClient = new ProtectedClient(transport);

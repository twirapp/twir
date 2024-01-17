import { TwirpFetchTransport } from '@protobuf-ts/twirp-transport';
import { ProtectedClient } from '@twir/api/api.client';
import { UnProtectedClient } from '@twir/api/api.client';

const transport = new TwirpFetchTransport({
	baseUrl: `${window.location.origin}/api/v1`,
	sendJson: import.meta.env.DEV,
});
export const protectedApiClient = new ProtectedClient(transport);
export const unprotectedApiClient = new UnProtectedClient(transport);

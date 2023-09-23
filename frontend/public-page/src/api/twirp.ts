import { TwirpFetchTransport } from '@protobuf-ts/twirp-transport';
import { UnProtectedClient } from '@twir/grpc/generated/api/api.client';

const transport = new TwirpFetchTransport({
	baseUrl: `${window.location.origin}/api/v1`,
	sendJson: import.meta.env.DEV,
});

export const unprotectedClient = new UnProtectedClient(transport);

import { TwirpFetchTransport } from '@protobuf-ts/twirp-transport';
import { ProtectedClient } from '@twir/grpc/generated/api/api.client';
import { UnProtectedClient } from '@twir/grpc/generated/api/api.client';

const transport = new TwirpFetchTransport({
	baseUrl: `${window.location.origin}/api/v1`,
	sendJson: process.env.NODE_ENV == 'development',
});
export const protectedApiClient = new ProtectedClient(transport);
export const unprotectedApiClient = new UnProtectedClient(transport);

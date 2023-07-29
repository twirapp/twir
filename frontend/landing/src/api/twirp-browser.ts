import { TwirpFetchTransport } from '@protobuf-ts/twirp-transport';
import { ProtectedClient, UnProtectedClient } from '@twir/grpc/generated/api/api.client';

const transport = new TwirpFetchTransport({
	baseUrl: `${window.location.origin}/api/v1`,
});

export const browserUnProtectedClient = new UnProtectedClient(transport);
export const browserProtectedClient = new ProtectedClient(transport);

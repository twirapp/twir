import { TwirpFetchTransport } from '@protobuf-ts/twirp-transport';
import { ProtectedClient } from '@twir/grpc/generated/api/api.client';
import { UnProtectedClient } from '@twir/grpc/generated/api/api.client';

export let protectedApiClient: ProtectedClient;
export let unprotectedApiClient: UnProtectedClient;

if (typeof window != 'undefined') {
	const transport = new TwirpFetchTransport({
		baseUrl: `${window.location.origin}/api/v1`,
		sendJson: process.env.NODE_ENV == 'development',
	});

	protectedApiClient = new ProtectedClient(transport);
	unprotectedApiClient = new UnProtectedClient(transport);
}

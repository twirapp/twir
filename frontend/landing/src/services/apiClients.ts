import { TwirpFetchTransport } from '@protobuf-ts/twirp-transport';
import { ProtectedClient } from '@twir/grpc/generated/api/api.client';
import { UnProtectedClient } from '@twir/grpc/generated/api/api.client';

export let protectedApiClient: ProtectedClient;
export let unprotectedApiClient: UnProtectedClient;

if (!import.meta.env.SSR) {
	const transport = new TwirpFetchTransport({
		baseUrl: `${window.location.origin}/api/v1`,
		sendJson: import.meta.env.DEV,
	});

	protectedApiClient = new ProtectedClient(transport);
	unprotectedApiClient = new UnProtectedClient(transport);
}

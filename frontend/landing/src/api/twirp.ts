import { TwirpFetchTransport } from '@protobuf-ts/twirp-transport';
import { config } from '@twir/config';
import { ProtectedClient, UnProtectedClient } from '@twir/grpc/generated/api/api.client';

const host = config.SITE_BASE_URL ?? 'localhost:3005';
const isDev = config.NODE_ENV === 'development';
const apiAddr = isDev ? `${host}/api` : 'api:3002';

const baseUrl = `http://${apiAddr}/v1`;

console.info('BaseURL:', baseUrl);

const transport = new TwirpFetchTransport({
	baseUrl,
	sendJson: isDev,
});

export const protectedClient = new ProtectedClient(transport);
export const unProtectedClient = new UnProtectedClient(transport);

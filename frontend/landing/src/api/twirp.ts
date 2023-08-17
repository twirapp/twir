import { TwirpFetchTransport } from '@protobuf-ts/twirp-transport';
import { ProtectedClient, UnProtectedClient } from '@twir/grpc/generated/api/api.client';

const host = process.env.HOSTNAME;
const isDev = process.env.NODE_ENV === 'development';
const baseUrl = `${host?.startsWith('localhost') || !isDev ? 'http' : 'https'}://${isDev ? `${host}/api` : 'api:3002'}/v1`;

console.info('BaseURL:', baseUrl);

const transport = new TwirpFetchTransport({
	baseUrl,
	sendJson: isDev,
});

export const protectedClient = new ProtectedClient(transport);
export const unProtectedClient = new UnProtectedClient(transport);

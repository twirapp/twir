import { TwirpFetchTransport } from '@protobuf-ts/twirp-transport';
import { ProtectedClient, UnProtectedClient } from '@twir/grpc/generated/api/api.client';

const host = process.env.HOSTNAME;
const isDev = process.env.NODE_ENV === 'development';
const baseUrl = `${host?.startsWith('localhost') ? 'http' : 'https'}://${isDev ? host : 'api:3002'}/api/v1`;

const transport = new TwirpFetchTransport({
	baseUrl,
	sendJson: isDev,
});

export const protectedClient = new ProtectedClient(transport);
export const unProtectedClient = new UnProtectedClient(transport);

import { TwirpFetchTransport } from '@protobuf-ts/twirp-transport';
import { ProtectedClient, UnProtectedClient } from '@twir/grpc/generated/api/api.client';

const host = import.meta.env.HOST;
const baseUrl = `${host.startsWith('localhost') ? 'http' : 'https'}://${import.meta.env.HOST}/api/v1`;

const transport = new TwirpFetchTransport({
  baseUrl,
  sendJson: import.meta.env.DEV,
});

export const protectedClient = new ProtectedClient(transport);
export const unProtectedClient = new UnProtectedClient(transport);

import { TwirpFetchTransport } from '@protobuf-ts/twirp-transport';
import { ProtectedClient, UnProtectedClient } from '@twir/grpc/generated/api/api.client';


const host = process.env.HOSTNAME;
const baseUrl = `${host?.startsWith('localhost') ? 'http' : 'https'}://${host}/api/v1`;

const transport = new TwirpFetchTransport({
  baseUrl,
  sendJson: process.env.NODE_ENV === 'development',
});

export const protectedClient = new ProtectedClient(transport);
export const unProtectedClient = new UnProtectedClient(transport);

import * as Websocket from '@tsuwari/grpc/generated/websocket/websocket';
import { PORTS } from '@tsuwari/grpc/servers/constants';
import { createServer } from 'nice-grpc';

import { onAddRequest, onRemoveRequest } from '../namespaces/youtube.js';

export const grpcServer = createServer({
  'grpc.keepalive_time_ms': 1 * 60 * 1000,
});

const websocketService: Websocket.WebsocketServiceImplementation = {
  youtubeAddSongToQueue: onAddRequest,
  youtubeRemoveSongToQueue: onRemoveRequest,
};


grpcServer.add(Websocket.WebsocketDefinition, websocketService);

export const listen = () => grpcServer.listen(`0.0.0.0:${PORTS.WEBSOCKET_SERVER_PORT}`);


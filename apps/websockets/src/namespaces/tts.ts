import { Empty } from '@tsuwari/grpc/generated/websockets/google/protobuf/empty';
import { TTSMessage, TTSSkipMessage } from '@tsuwari/grpc/generated/websockets/websockets';
import SocketIo from 'socket.io';

import { authMiddleware, io } from '../libs/io.js';

const sockets: Map<string, SocketIo.Socket> = new Map();
export const obsNameSpace = io.of('tts');

obsNameSpace.use(authMiddleware);
obsNameSpace.on('connection', async (socket) => {
  const channelId = socket.data.channelId as string;
  sockets.set(channelId, socket);
});

export const onTtsSay = async (data: TTSMessage): Promise<Empty> => {
  const socket = sockets.get(data.channelId);
  if (!socket) return {};

  socket.emit('say', data);

  return {};
};


export const onTtsSkip = async (data: TTSSkipMessage): Promise<Empty> => {
  const socket = sockets.get(data.channelId);
  if (!socket) return {};

  socket.emit('skip', data);

  return {};
};
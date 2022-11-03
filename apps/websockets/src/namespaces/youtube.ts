import { RequestedSong } from '@tsuwari/typeorm/entities/RequestedSong';
import SocketIo from 'socket.io';

import { typeorm } from '../libs/typeorm.js';
import { authMiddleware } from '../middlewares/auth.js';

const sockets: Map<string, SocketIo.Socket> = new Map();

export const createYoutubeNameSpace = (io: SocketIo.Server) => {
  const nameSpace = io.of('youtube');
  nameSpace.use(authMiddleware);
  nameSpace.on('connection', async (socket) => {
    const channelId = socket.handshake.auth.channelId;
    sockets.set(channelId, socket);
    socket.on('disconnect', () => {
      sockets.delete(channelId);
    });

    const songs = await typeorm.getRepository(RequestedSong).findBy({
      channelId,
    });
    socket.emit('currentQueue', songs);
  });
};

export async function addSongToQueue(channelId: string, entityId: string) {
  const entity = await typeorm.getRepository(RequestedSong).findOneBy({
    channelId,
    id: entityId,
  });
  if (!entity) return;

  const socket = sockets.get(channelId);
  if (!socket) return;

  socket.emit('addSong', entity);
}

export async function removeSongFromQueue(channelId: string, entityId: string) {
  const entity = await typeorm.getRepository(RequestedSong).findOneBy({
    channelId,
    id: entityId,
  });
  if (!entity) return;

  const socket = sockets.get(channelId);
  if (!socket) return;

  socket.emit('removeSong', entity);
}

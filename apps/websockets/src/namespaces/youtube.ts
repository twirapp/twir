import * as Youtube from '@tsuwari/nats/youtube';
import { IsNull } from '@tsuwari/typeorm';
import { RequestedSong } from '@tsuwari/typeorm/entities/RequestedSong';
import SocketIo from 'socket.io';

import { nats } from '../libs/nats.js';
import { typeorm } from '../libs/typeorm.js';
import { authMiddleware } from '../middlewares/auth.js';

const sockets: Map<string, SocketIo.Socket> = new Map();
const repository = typeorm.getRepository(RequestedSong);

export const createYoutubeNameSpace = async (io: SocketIo.Server) => {
  const nameSpace = io.of('youtube');
  nameSpace.use(authMiddleware);
  nameSpace.on('connection', async (socket) => {
    const channelId = socket.handshake.auth.channelId;
    sockets.set(channelId, socket);
    socket.on('disconnect', () => {
      sockets.delete(channelId);
    });

    const songs = await repository.findBy({
      channelId,
      deletedAt: IsNull(),
    });
    socket.emit('currentQueue', songs);

    socket.on('skip', async (id) => {
      const entity = await repository.findOneBy({ id });
      if (entity) {
        await repository.softDelete({ id });
      }
    });
  });

  for await (const event of nats.subscribe(Youtube.SUBJECTS.ADD_SONG_TO_QUEUE)) {
    const data = Youtube.AddSongToQueue.fromBinary(event.data);

    const socket = sockets.get(data.channelId);
    if (!socket) return;

    const entity = await repository.findOneBy({
      id: data.entityId,
    });
    socket.emit('newTrack', entity);
  }
};

export async function addSongToQueue(channelId: string, entityId: string) {
  const entity = await repository.findOneBy({
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

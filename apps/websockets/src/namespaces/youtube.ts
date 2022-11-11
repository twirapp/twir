import * as Youtube from '@tsuwari/nats/youtube';
import { IsNull } from '@tsuwari/typeorm';
import { RequestedSong } from '@tsuwari/typeorm/entities/RequestedSong';
import SocketIo from 'socket.io';

import { nats } from '../libs/nats.js';
import { redis } from '../libs/redis.js';
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

    socket.on('currentQueue', async (cb) => {
      const songs = await repository.findBy({
        channelId,
        deletedAt: IsNull(),
      });

      cb(songs);
    });

    socket.on('skip', async (id) => {
      const entity = await repository.findOneBy({ id });
      if (entity) {
        await repository.softDelete({ id });
      }
    });

    socket.on('disconnect', () => {
      sockets.delete(channelId);
    });

    socket.on('play', async (data) => {
      console.log(`songrequests:youtube:${channelId}:currentPlaying`);
      const result = await redis.set(
        `songrequests:youtube:${channelId}:currentPlaying`,
        data.id,
        'PX',
        data.timeToEnd,
      );
      console.log(result);
      console.log('play', data);
    });

    socket.on('pause', () => {
      redis.del(`songrequests:youtube:${channelId}:currentPlaying`);
    });
  });

  const addSubscription = nats.subscribe(Youtube.SUBJECTS.ADD_SONG_TO_QUEUE);
  const removeSubscription = nats.subscribe(Youtube.SUBJECTS.REMOVE_SONG_FROM_QUEUE);

  (async () => {
    for await (const event of addSubscription) {
      const data = Youtube.AddSongToQueue.fromBinary(event.data);

      const socket = sockets.get(data.channelId);
      if (!socket) return;

      const entity = await repository.findOneBy({
        id: data.entityId,
      });
      socket.emit('newTrack', entity);
    }
  })();

  (async () => {
    for await (const event of removeSubscription) {
      const data = Youtube.RemoveSongFromQueue.fromBinary(event.data);
      console.log(data);

      const socket = sockets.get(data.channelId);
      if (!socket) return;

      const entity = await repository.findOne({
        where: { id: data.entityId },
        withDeleted: true,
      });
      socket.emit('removeTrack', entity);
    }
  })();
};

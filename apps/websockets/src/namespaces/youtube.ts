import { Empty } from '@tsuwari/grpc/generated/websockets/google/protobuf/empty';
import {
  YoutubeAddSongToQueueRequest,
  YoutubeRemoveSongFromQueueRequest,
} from '@tsuwari/grpc/generated/websockets/websockets';
import { IsNull } from '@tsuwari/typeorm';
import { ChannelModuleSettings } from '@tsuwari/typeorm/entities/ChannelModuleSettings';
import { RequestedSong } from '@tsuwari/typeorm/entities/RequestedSong';
import { type YouTubeSettings } from '@tsuwari/types/api';
import SocketIo from 'socket.io';


import { botsGrpcClient } from '../libs/botsGrpc.js';
import { authMiddleware, io } from '../libs/io.js';
import { redis } from '../libs/redis.js';
import { typeorm } from '../libs/typeorm.js';

const sockets: Map<string, SocketIo.Socket> = new Map();
const repository = typeorm.getRepository(RequestedSong);

export const youtubeNamespace = io.of('youtube');

youtubeNamespace.use(authMiddleware);
youtubeNamespace.on('connection', async (socket) => {
  const channelId = socket.handshake.auth.channelId;
  sockets.set(channelId, socket);

  socket.on('currentQueue', async (cb) => {
    const songs = await repository.find({
      where: {
        channelId,
        deletedAt: IsNull(),
      },
      order: {
        queuePosition: 'asc',
      },
    });

    cb(songs);
  });

  socket.on('skip', async (id: string | string[]) => {
    const ids = Array.isArray(id) ? id : [id];
    for (const video of ids) {
      const entity = await repository.findOneBy({ id: video });
      if (entity) {
        await repository.softDelete({ id: video });
      }
    }

    redis.del(`songrequests:youtube:${channelId}:currentPlaying`);
  });

  socket.on('disconnect', () => {
    sockets.delete(channelId);
  });

  socket.on('play', async (data) => {
    const key = `songrequests:youtube:${channelId}:currentPlaying`;
    const current = await redis.get(key);
    const song = await repository.findOneBy({
      id: data.id,
    });
    const settingsEntity = await typeorm.getRepository(ChannelModuleSettings).findOneBy({
      channelId: song?.channelId,
    });

    if (settingsEntity && song) {
      const settings = settingsEntity.settings as YouTubeSettings;
      const announcePlay = settings.announcePlay ?? true;
      const message = settings.translations.nowPlaying
        .replace('{{songTitle}}', song.title)
        .replace('{{songId}}', song.videoId)
        .replace('{{orderedByName}}', song.orderedByName)
        .replace('{{orderedByDisplayName}}', song.orderedByDisplayName ?? song.orderedByName);

      if (!current && song && announcePlay) {
        await botsGrpcClient.sendMessage({
          channelId: song.channelId,
          isAnnounce: true,
          message,
        });
      }
    }


    await redis.set(key, data.id);
    await redis.expire(key, data.duration);
  });

  socket.on('pause', async () => {
    // await redis.set(`songrequests:youtube:${channelId}:currentPlaying`, {
    //   paused: true,
    // });
  });

  socket.on('newOrder', async (videos: RequestedSong[]) => {
    const currentVideosCount = await repository.count({
      where: { channelId },
    });

    if (videos.some((v) => v.queuePosition > currentVideosCount + 2)) {
      return;
    }

    for (const video of videos) {
      await repository.update({ id: video.id }, { queuePosition: video.queuePosition });
    }
  });
});

export const onAddRequest = async (data: YoutubeAddSongToQueueRequest): Promise<Empty> => {
  const socket = sockets.get(data.channelId);
  if (!socket) return {};

  const entity = await repository.findOneBy({
    id: data.entityId,
  });
  socket.emit('newTrack', entity);

  return {};
};

export const onRemoveRequest = async (data: YoutubeRemoveSongFromQueueRequest): Promise<Empty> => {
  const socket = sockets.get(data.channelId);
  if (!socket) return {};

  const entity = await repository.findOne({
    where: { id: data.entityId },
    withDeleted: true,
  });
  socket.emit('removeTrack', entity);

  return {};
};

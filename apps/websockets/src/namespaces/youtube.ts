import * as Youtube from '@tsuwari/nats/youtube';
import { IsNull } from '@tsuwari/typeorm';
import { ChannelModuleSettings, ModuleType } from '@tsuwari/typeorm/entities/ChannelModuleSettings';
import { RequestedSong } from '@tsuwari/typeorm/entities/RequestedSong';
import { YoutubeSettings } from '@tsuwari/types/generated';
import SocketIo from 'socket.io';

import { nats } from '../libs/nats.js';
import { redis } from '../libs/redis.js';
import { typeorm } from '../libs/typeorm.js';
import { authMiddleware } from '../middlewares/auth.js';

const sockets: Map<string, SocketIo.Socket> = new Map();
const repository = typeorm.getRepository(RequestedSong);
const settingsRepository = typeorm.getRepository(ChannelModuleSettings);

const defaultSettings: YoutubeSettings = {
  acceptOnlyWhenOnline: true,
  channelPointsRewardName: '',
  maxRequests: 200,
  blacklist: {
    artistsNames: [],
    channels: [],
    songs: [],
    users: [],
  },
  song: { acceptedCategories: [] },
  user: {},
};

const findOrCreateSettings = async (channelId: string): Promise<ChannelModuleSettings & { settings: YoutubeSettings }> => {
  const settings = await settingsRepository.findOneBy({
    channelId,
    type: ModuleType.YOUTUBE_SONG_REQUESTS,
  });
  if (settings) return settings ;

  return settingsRepository.save({
    channelId,
    type: ModuleType.YOUTUBE_SONG_REQUESTS,
    settings: defaultSettings,
  });
};

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
      await redis.set(
        `songrequests:youtube:${channelId}:currentPlaying`,
        data.id,
        'PX',
        data.timeToEnd,
      );
    });


    socket.on('blacklist.user', async (data) => {
      const { userId, userName } = data;
      if (!userId || !userName) {
        return;
      }
      const settings = await findOrCreateSettings(channelId);
      settings.settings.blacklist?.users.push({ userId, userName });
      await settingsRepository.save(settings);
    });

    socket.on('blacklist.song', async (data) => {
      const { id, title, thumbNail } = data;
      if (!id || !title || !thumbNail) {
        return;
      }
      const settings = await findOrCreateSettings(channelId);
      settings.settings.blacklist?.songs.push({ id, title, thumbNail });
      await settingsRepository.save(settings);
    });

    socket.on('blacklist.channel', async (data) => {
      const { id, title, thumbNail } = data;
      if (!id || !title || !thumbNail) {
        return;
      }
      const settings = await findOrCreateSettings(channelId);
      settings.settings.blacklist?.channels.push({ id, title, thumbNail });
      await settingsRepository.save(settings);
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

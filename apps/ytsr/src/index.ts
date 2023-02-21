import * as YTSR from '@tsuwari/grpc/generated/ytsr/ytsr';
import { PORTS } from '@tsuwari/grpc/servers/constants';
import { createServer } from 'nice-grpc';
import ytdl from 'ytdl-core';
import ytsrLib, { Video } from 'ytsr';

import { durationToMilliseconds } from './utils/convertDuration.js';

export const grpcServer = createServer({
  'grpc.keepalive_time_ms': 1 * 60 * 1000,
});

const linkRegexp = /(?:https?:\/\/)?(?:www\.)?(?:youtube\.com|youtu\.be)\/(?:watch\?v=|embed\/|v\/|.+\?v=)?([\w-]{11})(?:\S+)?/g;


const ytsrService: YTSR.YtsrServiceImplementation = {
  async search(request: YTSR.SearchRequest, context): Promise<YTSR.DeepPartial<YTSR.SearchResponse>> {
    const videos: Array<YTSR.Song> = [];

    const linkMatches = [...request.search.matchAll(linkRegexp)];
    if (linkMatches.length) {
      await Promise.all(linkMatches.map(async (match) => {
        const song = await ytdl.getInfo(match[0]).catch(() => null);
        if (!song) return;
        videos.push({
          title: song.videoDetails.title,
          isLive: song.videoDetails.isLiveContent,
          duration: Number(song.videoDetails.lengthSeconds) * 1000,
          thumbnailUrl: song.videoDetails.thumbnails[0]?.url,
          author: {
            name: song.videoDetails.author.name,
            avatarUrl: song.videoDetails.author.thumbnails?.at(0)?.url,
            channelId: song.videoDetails.author.id,
          },
          id: song.videoDetails.videoId,
          views: Number(song.videoDetails.viewCount),
        });
      }));
    } else {
      const search = await ytsrLib(request.search, { limit: 1 });
      if (search.items.length && search.items?.at(0)?.type === 'video') {
        const item = search.items.at(0) as ytsrLib.Video;
        videos.push({
          title: item.title,
          isLive: item.isLive,
          views: item.views ?? 0,
          id: item.id,
          thumbnailUrl: item.bestThumbnail?.url ?? undefined,
          duration: durationToMilliseconds(item.duration ?? '0:00'),
          author: {
            name: item.author?.name || '',
            channelId: item.author?.channelID || '',
            avatarUrl: item.author?.bestAvatar?.url || '',
          },
        });
      }
    }

    return {
      songs: videos,
    };
  },
};

grpcServer.add(YTSR.YtsrDefinition, ytsrService);

grpcServer.listen(`0.0.0.0:${PORTS.YTSR_SERVER_PORT}`);

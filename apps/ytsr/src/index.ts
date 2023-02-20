import * as YTSR from '@tsuwari/grpc/generated/ytsr/ytsr';
import { PORTS } from '@tsuwari/grpc/servers/constants';
import { createServer } from 'nice-grpc';
import ytsrLib from 'ytsr';

import { durationToMilliseconds } from './utils/convertDuration.js';

export const grpcServer = createServer({
  'grpc.keepalive_time_ms': 1 * 60 * 1000,
});

const linkRegexp = /(?:https?:\/\/)?(?:www\.)?(?:youtube\.com|youtu\.be)\/(?:watch\?v=|embed\/|v\/|.+\?v=)?([\w-]{11})(?:\S+)?/g;

const ytsrService: YTSR.YtsrServiceImplementation = {
  async search(request: YTSR.SearchRequest, context): Promise<YTSR.DeepPartial<YTSR.SearchResponse>> {
    const videos: Array<ytsrLib.Video> = [];

    const searchAndPush = async (input: string) => {
      const search = await ytsrLib(input, { limit: 1 });
      if (search.items.length && search.items[0]?.type === 'video') {
        videos.push(search.items[0]);
      }
    };

    const linkMatches = [...request.search.matchAll(linkRegexp)];
    if (linkMatches.length) {
      await Promise.all(linkMatches.map(async (match) => {
        await searchAndPush(match[1]!);
      }));
    } else {
      await searchAndPush(request.search);
    }

    return {
      songs: videos.map((item) => {
        return {
          title: item.title,
          isLive: item.isLive,
          views: item.views ?? 0,
          id: item.id,
          author: item.author ? {
            name: item.author.name,
            avatarUrl: item.author.bestAvatar?.url ?? undefined,
            channelId: item.author.channelID,
          } : undefined,
          thumbnailUrl: item.bestThumbnail.url ?? undefined,
          duration: durationToMilliseconds(item.duration ?? '0:00'),
        };
      }),
    };
  },
};

grpcServer.add(YTSR.YtsrDefinition, ytsrService);

grpcServer.listen(`0.0.0.0:${PORTS.YTSR_SERVER_PORT}`);

import { createSpotifyAPI } from '@soundify/api';
import { ClientCredentials } from '@soundify/node-auth';
import { config } from '@tsuwari/config';
import * as YTSR from '@tsuwari/grpc/generated/ytsr/ytsr';
import { PORTS } from '@tsuwari/grpc/servers/constants';
import { createServer } from 'nice-grpc';
import ytsrLib, { Video } from 'ytsr';

import { durationToMilliseconds } from './utils/convertDuration.js';

export const grpcServer = createServer({
  'grpc.keepalive_time_ms': 1 * 60 * 1000,
});

const youtubeLinkRegexp =
  /(?:https?:\/\/)?(?:www\.)?(?:youtube\.com|youtu\.be)\/(?:watch\?v=|embed\/|v\/|.+\?v=)?([\w-]{11})(?:\S+)?/g;
const spotifyLinkRegexp =
  /^(https:\/\/open.spotify.com\/track\/|spotify:track:)([a-zA-Z0-9]+)(.*)$/g;

const spotifyClient = createSpotifyAPI(
  new ClientCredentials({
    client_id: config.SPOTIFY_CLIENT_ID,
    client_secret: config.SPOTIFY_CLIENT_SECRET,
  }).createAuthProvider(),
);

const ytsrService: YTSR.YtsrServiceImplementation = {
  async search(
    request: YTSR.SearchRequest,
    context,
  ): Promise<YTSR.DeepPartial<YTSR.SearchResponse>> {
    const videos: Array<YTSR.Song> = [];

    const tracksForSearch: string[] = [];

    const youtubeLinkMatches = [...request.search.matchAll(youtubeLinkRegexp)];
    const spotifyLinkMatches = [...request.search.matchAll(spotifyLinkRegexp)];

    await Promise.all(
      spotifyLinkMatches.map(async (match) => {
        const song = await spotifyClient.getTrack(match[2]!).catch(() => null);
        if (!song) return;

        tracksForSearch.push(`${song.artists.map((a) => a.name)} ${song.name}`);
      }),
    );

    for (const match of youtubeLinkMatches) {
      tracksForSearch.push(match[1]!);
    }

    if (!youtubeLinkMatches.length && !spotifyLinkMatches.length) {
      tracksForSearch.push(request.search);
    }

    await Promise.all(
      tracksForSearch.map(async (track) => {
        const search = await ytsrLib(track, { limit: 1 });
        const item = search.items.at(0) as Video;
        if (!item) return;

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
      }),
    );

    return {
      songs: videos,
    };
  },
};

grpcServer.add(YTSR.YtsrDefinition, ytsrService);

grpcServer.listen(`0.0.0.0:${PORTS.YTSR_SERVER_PORT}`);

process.on('SIGINT', () => grpcServer.shutdown()).on('SIGTERM', () => grpcServer.shutdown());

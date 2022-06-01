import { URLSearchParams } from 'url';

import { prisma } from '../libs/prisma.js';

export interface LastFMDbData {
  username: string;
}

interface ILastfmResponse {
  error?: number;
  message?: string;
  recenttracks?: {
    track?: Array<{
      artist: {
        '#text': string;
      };
      '@attr'?: {
        nowplaying: boolean;
      };
      album?: {
        '#text': string;
      };
      name: string;
    }>;
  };
}

class LastFmIntegrationClass {
  #apiKey: string;

  async init() {
    const service = await prisma.integration.findFirst({
      where: {
        service: 'LASTFM',
      },
    });

    if (service?.apiKey) {
      this.#apiKey = service.apiKey;
    }

    return this;
  }

  async fetchSong(userName: string) {
    const params = {
      method: 'user.getrecenttracks',
      user: userName,
      api_key: this.#apiKey,
      format: 'json',
      limit: '1',
    } as Record<string, string>;

    const request = await fetch(`http://ws.audioscrobbler.com/2.0?${new URLSearchParams(params)}`, {
      method: 'GET',
    });

    const response = await request.json() as ILastfmResponse;

    const song = response.recenttracks?.track?.find((t) => t['@attr']?.nowplaying);
    if (!song) return null;

    return `${song.artist['#text']} — ${song.name}`;
  }
}

export const LastFmIntegration = await new LastFmIntegrationClass().init();

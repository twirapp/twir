import { URLSearchParams } from 'url';

import { Injectable, OnApplicationBootstrap, OnModuleInit } from '@nestjs/common';
import { PrismaService } from '@tsuwari/prisma';

import { app } from '../index.js';

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

@Injectable()
export class LastFmIntegration implements OnModuleInit {
  #apiKey: string;

  constructor(private readonly prisma: PrismaService) { }

  async onModuleInit() {
    const service = await this.prisma.integration.findFirst({
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

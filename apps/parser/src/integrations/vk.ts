import { URLSearchParams } from 'url';

import { Injectable, OnApplicationBootstrap, OnModuleInit } from '@nestjs/common';
import { PrismaService } from '@tsuwari/prisma';


export interface VKDBData {
  userId: string;
}

interface IVkResponse {
  error?: unknown;
  response?: {
    text?: string;
    audio?: {
      artist: string;
      title: string;
    };
  };
}

@Injectable()
export class VkIntegration implements OnModuleInit {
  #accessToken: string;

  constructor(private readonly prisma: PrismaService) { }

  async onModuleInit() {
    const service = await this.prisma.integration.findFirst({
      where: {
        service: 'VK',
      },
    });

    if (service && service.accessToken) {
      this.#accessToken = service.accessToken;
    }

    return this;
  }

  async fetchSong(vkUserId: string) {
    const params = {
      access_token: this.#accessToken,
      uid: vkUserId,
      v: '5.131',
    };

    const request = await fetch(`https://api.vk.com/method/status.get?${new URLSearchParams(params)}`, {
      method: 'GET',
    });

    const response = await request.json() as IVkResponse;

    if (response.error) return null;
    if (!response.response || !response.response.audio) return null;

    return `${response.response.audio?.artist} — ${response.response.audio?.title}`;
  }
}

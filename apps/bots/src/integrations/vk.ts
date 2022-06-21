import { URLSearchParams } from 'url';

import { prisma } from '../libs/prisma.js';

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

class VkIntegrationClass {
  #accessToken: string;

  async init() {
    const service = await prisma.integration.findFirst({
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

export const VKIntegration = await new VkIntegrationClass().init();

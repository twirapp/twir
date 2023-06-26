import { Repository } from '@twir/typeorm';
import { ChannelIntegration } from '@twir/typeorm/entities/ChannelIntegration';
import { IntegrationService } from '@twir/typeorm/entities/Integration';
import axios, { AxiosError } from 'axios';

import { SpotifyIntegration } from './integration.js';

export class UserIntegration {
  #axios = axios.create({});
  #integration: ChannelIntegration;

  constructor(
    private readonly userId: string,
    private readonly spotify: SpotifyIntegration,
    private readonly repository: Repository<ChannelIntegration>,
  ) {
    this.#axios.interceptors.response.use(
      (response) => response,
      async (error: AxiosError & { config: { __retry: boolean } }) => {
        const response = error.response;

        if (response && error.config && response.status === 401 && !error.config.__retry) {
          error.config.__retry = true;

          const newAccessToken = await this.refreshToken();

          return axios({
            ...error.config,
            headers: {
              ...error.config.headers,
              authorization: `Bearer ${newAccessToken.accessToken}`,
            },
          });
        }

        return Promise.reject(error);
      },
    );
  }

  async init(int?: ChannelIntegration) {
    if (this.#integration) return this;
    if (int) {
      this.#integration = int;
      return this;
    }

    const integration = await this.repository.findOneBy({
      channelId: this.userId,
      integration: {
        service: IntegrationService.SPOTIFY,
      },
    });

    if (!integration?.refreshToken || !integration.accessToken) return null;

    this.#integration = integration;
    return this;
  }

  async refreshToken() {
    const service = await this.spotify.getSettings();
    if (!service) throw new Error('Service not setuped.');

    try {
      const token = Buffer.from(service.clientId + ':' + service.clientSecret).toString('base64');

      const request = await this.#axios.post(
        'https://accounts.spotify.com/api/token',
        new URLSearchParams({
          grant_type: 'refresh_token',
          refresh_token: this.#integration.refreshToken!,
        }),
        {
          headers: {
            Authorization: `Basic ${token}`,
            'Content-Type': 'application/x-www-form-urlencoded',
          },
        },
      );

      const { access_token: accessToken } = request.data;

      if (!accessToken) {
        throw new Error('fail');
      }

      await this.repository.update({ id: this.#integration.id }, { accessToken });
      this.#integration = await this.repository.findOneBy({ id: this.#integration.id });

      return { accessToken };
    } catch (e) {
      console.log(e);
      return null;
    }
  }

  async getCurrentSong() {
    try {
      const request = await this.#axios.get(
        `https://api.spotify.com/v1/me/player/currently-playing`,
        {
          headers: {
            Authorization: `Bearer ${this.#integration.accessToken}`,
          },
        },
      );

      const track = request.data?.item;

      if (!track) return null;

      return `${track.artists.map((a: { name: string }) => a.name).join(', ')} — ${track.name}`;
    } catch (error) {
      console.error(error);
      return null;
    }
  }

  async getProfile() {
    const request = await this.#axios.get(`https://api.spotify.com/v1/me`, {
      headers: {
        Authorization: `Bearer ${this.#integration.accessToken}`,
      },
    });

    return request.data;
  }
}

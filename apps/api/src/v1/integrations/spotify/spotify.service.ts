import { Injectable } from '@nestjs/common';
import { PrismaService } from '@tsuwari/prisma';
import { UserIntegration } from '@tsuwari/spotify';

import { SpotifyIntegrationService } from './integration.js';

@Injectable()
export class SpotifyService {
  constructor(private readonly prisma: PrismaService, private readonly spotify: SpotifyIntegrationService) { }

  async getTokens(userId: string, code: string) {
    const service = await this.spotify.getSettings();
    if (!service) throw new Error('Spotify not setuped.');

    const token = Buffer.from(service.clientId + ':' + service.clientSecret).toString('base64');

    const request = await fetch('https://accounts.spotify.com/api/token', {
      method: 'POST',
      body: new URLSearchParams({
        code,
        redirect_uri: service.redirectUrl,
        grant_type: 'authorization_code',
      }),
      headers: {
        authorization: `Basic ${token}`,
        'content-type': 'application/x-www-form-urlencoded',
      },
    });

    if (!request.ok) throw new Error(`Cannot get tokens`);

    const response = await request.json() as {
      'access_token': string,
      'token_type': string,
      'scope': string,
      'expires_in': number,
      'refresh_token': string
    };

    const { access_token: accessToken, refresh_token: refreshToken } = response;

    const currentIntegration = await this.prisma.channelIntegration.findFirst({
      where: { channelId: userId },
    });

    const data = { accessToken, refreshToken, enabled: true };
    if (currentIntegration) {
      await this.prisma.channelIntegration.update({
        where: { id: currentIntegration.id },
        data,
      });
    } else {
      await this.prisma.channelIntegration.create({
        data: {
          ...data,
          channel: { connect: { id: userId } },
          integration: { connect: { id: service.id } },
        },
      });
    }


    return { accessToken, refreshToken };
  }

  async getProfile(userId: string) {
    const integration = await this.prisma.channelIntegration.findFirst({
      where: {
        channelId: userId,
        integration: {
          service: 'SPOTIFY',
        },
      },
    });

    if (!integration) return null;

    const instance = new UserIntegration(userId, this.spotify, this.prisma);
    const spotifyIntegration = await instance.init(integration);
    const profile = spotifyIntegration.getProfile();

    return profile;
  }
}

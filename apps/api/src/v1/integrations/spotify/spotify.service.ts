import { HttpException, Injectable } from '@nestjs/common';
import { UserIntegration } from '@tsuwari/spotify';
import { ChannelIntegration } from '@tsuwari/typeorm/entities/ChannelIntegration';
import { Integration, IntegrationService } from '@tsuwari/typeorm/entities/Integration';

import { typeorm } from '../../../index.js';
import { UpdateSpotifyIntegrationDto } from './dto/patch.js';
import { SpotifyIntegrationService } from './integration.js';

@Injectable()
export class SpotifyService {
  constructor(private readonly spotify: SpotifyIntegrationService) {}

  async getAuthLink() {
    const integration = await typeorm.getRepository(Integration).findOneBy({
      service: IntegrationService.SPOTIFY,
    });

    if (!integration)
      throw new HttpException(
        'Spotify service not enabled on our side. Please, contact site owner.',
        404,
      );

    return (
      `https://accounts.spotify.com/authorize?` +
      new URLSearchParams({
        response_type: 'code',
        client_id: integration.clientId!,
        scope: 'user-read-currently-playing',
        redirect_uri: integration.redirectUrl!,
      })
    );
  }

  async getTokens(userId: string, code: string) {
    const service = await this.spotify.getSettings();
    if (!service)
      throw new HttpException(
        'Spotify service not enabled on our side. Please, contact site owner.',
        404,
      );

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

    if (!request.ok) throw new HttpException(`Cannot get tokens`, 404);

    const response = (await request.json()) as {
      access_token: string;
      token_type: string;
      scope: string;
      expires_in: number;
      refresh_token: string;
    };

    const { access_token: accessToken, refresh_token: refreshToken } = response;

    const repository = typeorm.getRepository(ChannelIntegration);
    const currentIntegration = await repository.findOneBy({
      channelId: userId,
      integration: {
        service: IntegrationService.SPOTIFY,
      },
    });

    const data = { accessToken, refreshToken, enabled: true };
    if (currentIntegration) {
      await repository.update({ id: currentIntegration.id, enabled: true }, data);
    } else {
      await repository.save({
        ...data,
        channelId: userId,
        enabled: true,
        integrationId: service.id,
      });
    }

    return { accessToken, refreshToken };
  }

  async getProfile(userId: string) {
    const repository = typeorm.getRepository(ChannelIntegration);
    const integration = await repository.findOneBy({
      channelId: userId,
      integration: {
        service: IntegrationService.SPOTIFY,
      },
    });

    if (!integration) return null;

    const instance = new UserIntegration(userId, this.spotify, repository);
    const spotifyIntegration = await instance.init(integration);
    const profile = spotifyIntegration.getProfile();

    return profile;
  }

  async getIntegration(channelId: string) {
    return typeorm.getRepository(ChannelIntegration).findOneBy({
      channelId,
      integration: {
        service: IntegrationService.SPOTIFY,
      },
    });
  }

  async updateIntegration(channelId: string, body: UpdateSpotifyIntegrationDto) {
    const integration = await this.getIntegration(channelId);

    if (!integration) {
      throw new HttpException('You need to authorize first.', 404);
    }

    const repository = typeorm.getRepository(ChannelIntegration);
    await repository.update({ id: integration.id }, body);
    return this.getIntegration(channelId);
  }
}

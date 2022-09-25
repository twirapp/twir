import { HttpException, Injectable } from '@nestjs/common';
import { ChannelIntegration } from '@tsuwari/typeorm/entities/ChannelIntegration';
import { Integration, IntegrationService } from '@tsuwari/typeorm/entities/Integration';

import { typeorm } from '../../../index.js';
import { UpdateStreamlabsIntegrationDto } from './dto/patch.js';

@Injectable()
export class StreamlabsService {
  async getAuthLink() {
    const integration = await typeorm.getRepository(Integration).findOneBy({
      service: IntegrationService.STREAMLABS,
    });

    if (!integration)
      throw new HttpException(
        'Streamlabs service not enabled on our side. Please, contact site owner.',
        404,
      );

    return (
      `https://www.streamlabs.com/api/v1.0/authorize?` +
      new URLSearchParams({
        response_type: 'code',
        client_id: integration.clientId!,
        scope: 'socket.token donations.read',
        redirect_uri: integration.redirectUrl!,
      })
    );
  }

  async getTokens(userId: string, code: string) {
    const service = await typeorm.getRepository(Integration).findOneBy({
      service: IntegrationService.STREAMLABS,
    });
    if (!service)
      throw new HttpException(
        'Streamlabs service not enabled on our side. Please, contact site owner.',
        404,
      );

    const body = new URLSearchParams({
      grant_type: 'authorization_code',
      client_id: service.clientId!,
      client_secret: service.clientSecret!,
      redirect_uri: service.redirectUrl!,
      code,
    });

    const request = await fetch('https://streamlabs.com/api/v1.0/token', {
      method: 'POST',
      body: body.toString(),
      headers: {
        'content-type': 'application/x-www-form-urlencoded',
      },
    });

    if (!request.ok) {
      console.log(await request.text());
      throw new HttpException(`Cannot get tokens`, 404);
    }

    const response = (await request.json()) as {
      access_token: string;
      token_type: string;
      expires_in: number;
      refresh_token: string;
    };

    const { access_token: accessToken, refresh_token: refreshToken } = response;

    const repository = typeorm.getRepository(ChannelIntegration);
    const currentIntegration = await repository.findOneBy({
      channelId: userId,
      integration: {
        service: IntegrationService.STREAMLABS,
      },
    });

    const profileRequest = await fetch(
      `https://streamlabs.com/api/v1.0/user?access_token=${accessToken}`,
    );

    if (!profileRequest.ok) {
      throw new HttpException('Cannot fetch your streamlabs profile', 404);
    }

    const profile = await profileRequest.json();

    const integrationData = {
      accessToken,
      refreshToken,
      enabled: true,
      data: profile.streamlabs,
    };
    if (currentIntegration) {
      await repository.update({ id: currentIntegration.id, enabled: true }, integrationData);
    } else {
      await repository.save({
        ...integrationData,
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
        service: IntegrationService.STREAMLABS,
      },
    });

    if (!integration) return null;

    return integration.data;
  }

  async getIntegration(channelId: string) {
    return typeorm.getRepository(ChannelIntegration).findOneBy({
      channelId,
      integration: {
        service: IntegrationService.STREAMLABS,
      },
    });
  }

  async updateIntegration(channelId: string, body: UpdateStreamlabsIntegrationDto) {
    const integration = await typeorm.getRepository(ChannelIntegration).findOneBy({
      channelId,
      integration: {
        service: IntegrationService.STREAMLABS,
      },
    });

    if (!integration) {
      throw new HttpException('You need to authorize first.', 404);
    }

    const repository = typeorm.getRepository(ChannelIntegration);
    await repository.update({ id: integration.id }, body);
    return this.getIntegration(channelId);
  }
}

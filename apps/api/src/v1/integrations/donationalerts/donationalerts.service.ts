import { HttpException, Injectable } from '@nestjs/common';
import * as NatsIntegration from '@tsuwari/nats/integrations';
import { ChannelIntegration } from '@tsuwari/typeorm/entities/ChannelIntegration';
import { Integration, IntegrationService } from '@tsuwari/typeorm/entities/Integration';

import { typeorm } from '../../../index.js';
import { nats } from '../../../libs/nats.js';
import { UpdateDonationAlertsIntegrationDto } from './dto/patch.js';

@Injectable()
export class DonationAlertsService {
  async getAuthLink() {
    const integration = await typeorm.getRepository(Integration).findOneBy({
      service: IntegrationService.DONATIONALERTS,
    });

    if (!integration)
      throw new HttpException(
        'DonationAlerts service not enabled on our side. Please, contact site owner.',
        404,
      );

    return (
      `https://www.donationalerts.com/oauth/authorize?` +
      new URLSearchParams({
        response_type: 'code',
        client_id: integration.clientId!,
        scope: 'oauth-user-show oauth-donation-subscribe',
        redirect_uri: integration.redirectUrl!,
      })
    );
  }

  async getTokens(userId: string, code: string) {
    const service = await typeorm.getRepository(Integration).findOneBy({
      service: IntegrationService.DONATIONALERTS,
    });
    if (!service)
      throw new HttpException(
        'DonationAlerts service not enabled on our side. Please, contact site owner.',
        404,
      );

    const request = await fetch('https://www.donationalerts.com/oauth/token', {
      method: 'POST',
      body: new URLSearchParams({
        grant_type: 'authorization_code',
        client_id: service.clientId!,
        client_secret: service.clientSecret!,
        redirect_uri: service.redirectUrl!,
        code,
      }).toString(),
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
        service: IntegrationService.DONATIONALERTS,
      },
    });

    const profileRequest = await fetch('https://www.donationalerts.com/api/v1/user/oauth', {
      headers: {
        Authorization: `Bearer ${accessToken}`,
      },
    });

    if (!profileRequest.ok) {
      throw new HttpException('Cannot fetch your donationalerts profile', 404);
    }

    const profile = await profileRequest.json();

    const integrationData = {
      accessToken,
      refreshToken,
      enabled: true,
      data: {
        name: profile.data.name,
        code: profile.data.code,
        avatar: profile.data.avatar,
      },
    };
    let integrationId = currentIntegration?.id;
    if (currentIntegration) {
      await repository.update({ id: currentIntegration.id, enabled: true }, integrationData);
    } else {
      const { id } = await repository.save({
        ...integrationData,
        channelId: userId,
        enabled: true,
        integrationId: service.id,
      });
      integrationId = id;
    }

    nats.publish(
      'integrations.add',
      NatsIntegration.AddIntegration.toBinary({ id: integrationId! }),
    );

    return { accessToken, refreshToken };
  }

  async getProfile(userId: string) {
    const repository = typeorm.getRepository(ChannelIntegration);
    const integration = await repository.findOneBy({
      channelId: userId,
      integration: {
        service: IntegrationService.DONATIONALERTS,
      },
    });

    if (!integration) return null;

    return integration.data;
  }

  async getIntegration(channelId: string) {
    return typeorm.getRepository(ChannelIntegration).findOneBy({
      channelId,
      integration: {
        service: IntegrationService.DONATIONALERTS,
      },
    });
  }

  async updateIntegration(channelId: string, body: UpdateDonationAlertsIntegrationDto) {
    const integration = await typeorm.getRepository(ChannelIntegration).findOneBy({
      channelId,
      integration: {
        service: IntegrationService.DONATIONALERTS,
      },
    });

    if (!integration) {
      throw new HttpException('You need to authorize first.', 404);
    }

    if (body.enabled) {
      nats.publish(
        'integrations.add',
        NatsIntegration.AddIntegration.toBinary({ id: integration.id }),
      );
    } else {
      nats.publish(
        'integrations.remove',
        NatsIntegration.RemoveIntegration.toBinary({ id: integration.id }),
      );
    }

    const repository = typeorm.getRepository(ChannelIntegration);
    await repository.update({ id: integration.id }, body);
    return this.getIntegration(channelId);
  }
}

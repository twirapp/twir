import { HttpException, Injectable } from '@nestjs/common';
import { ChannelIntegration } from '@tsuwari/typeorm/entities/ChannelIntegration';
import { Integration, IntegrationService } from '@tsuwari/typeorm/entities/Integration';

import { typeorm } from '../../../index.js';
import { LastfmUpdateDto } from './dto/update.js';

@Injectable()
export class LastfmService {
  async getIntegration(channelId: string) {
    const integration = await typeorm.getRepository(ChannelIntegration).findOneBy({
      channelId,
      integration: {
        service: IntegrationService.LASTFM,
      },
    });

    return integration;
  }

  async updateIntegration(channelId: string, body: LastfmUpdateDto) {
    const integrationService = await typeorm.getRepository(Integration).findOne({
      where: {
        service: IntegrationService.LASTFM,
      },
    });

    if (!integrationService)
      throw new HttpException(
        `LastFM not enabled on our backed. Please, make patience or contact us`,
        404,
      );

    const repository = typeorm.getRepository(ChannelIntegration);
    const integration = await this.getIntegration(channelId);
    let integrationId = integration?.id;
    if (!integration) {
      const newIntegration = await repository.save({
        channelId,
        enabled: body.enabled,
        data: body.data as any,
        integrationId: integrationService.id,
      });
      integrationId = newIntegration.id;
    } else {
      await repository.update(
        { id: integrationId },
        { enabled: body.enabled, data: body.data as any },
      );
    }

    return repository.findOneBy({ id: integrationId });
  }
}

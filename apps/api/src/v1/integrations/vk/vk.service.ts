import { HttpException, Injectable } from '@nestjs/common';
import { ChannelIntegration } from '@tsuwari/typeorm/entities/ChannelIntegration';
import { Integration, IntegrationService } from '@tsuwari/typeorm/entities/Integration';

import { typeorm } from '../../../index.js';
import { VkUpdateDto } from './dto/update.js';

@Injectable()
export class VkService {
  getIntegration(channelId: string) {
    return typeorm.getRepository(ChannelIntegration).findOneBy({
      channelId,
      integration: {
        service: IntegrationService.VK,
      },
    });
  }

  async updateIntegration(channelId: string, data: VkUpdateDto) {
    const integrationService = await typeorm.getRepository(Integration).findOneBy({
      service: IntegrationService.VK,
    });

    if (!integrationService)
      throw new HttpException(
        `Vk not enabled on our backed. Please, make patience or contact us`,
        404,
      );

    const repository = typeorm.getRepository(ChannelIntegration);
    const integration = await this.getIntegration(channelId);

    if (!integration) {
      await repository.save({
        channelId,
        enabled: data.enabled,
        data: data.data as any,
        integrationId: integrationService.id,
      });
    } else {
      await repository.update({ id: integration.id }, { ...data, data: data.data as any });
    }

    return this.getIntegration(channelId);
  }
}

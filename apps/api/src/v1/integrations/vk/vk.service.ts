import { HttpException, Injectable } from '@nestjs/common';
import { Prisma, PrismaService } from '@tsuwari/prisma';

import { VkUpdateDto } from './dto/update.js';

@Injectable()
export class VkService {
  constructor(private readonly prisma: PrismaService) { }

  async getIntegration(channelId: string) {
    const integration = await this.prisma.channelIntegration.findFirst({
      where: {
        channelId,
        integration: {
          service: 'VK',
        },
      },
    });

    return integration;
  }

  async updateIntegration(channelId: string, data: VkUpdateDto) {
    const integrationService = await this.prisma.integration.findFirst({
      where: {
        service: 'VK',
      },
    });

    if (!integrationService) throw new HttpException(`Vk not enabled on our backed. Please, make patience or contact us`, 404);

    let integration = await this.getIntegration(channelId);

    if (!integration) {
      integration = await this.prisma.channelIntegration.create({
        data: {
          channelId,
          data: {
            ...data,
            data: { ...data.data } as unknown as Prisma.InputJsonObject,
          },
          integrationId: integrationService.id,
        },
      });
    } else {
      integration = await this.prisma.channelIntegration.update({
        where: {
          id: integration.id,
        },
        data: {
          ...data,
          data: { ...data.data } as unknown as Prisma.InputJsonObject,
        },
      });
    }

    return integration;
  }
}

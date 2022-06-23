import { HttpException, Injectable } from '@nestjs/common';
import { Prisma, PrismaService } from '@tsuwari/prisma';

import { FaceitUpdateDto } from './dto/update.js';

@Injectable()
export class FaceitService {
  constructor(private readonly prisma: PrismaService) { }

  async getIntegration(channelId: string) {
    const integration = await this.prisma.channelIntegration.findFirst({
      where: {
        channelId,
        integration: {
          service: 'FACEIT',
        },
      },
    });

    return integration;
  }

  async updateIntegration(channelId: string, body: FaceitUpdateDto) {
    const integrationService = await this.prisma.integration.findFirst({
      where: {
        service: 'FACEIT',
      },
    });

    if (!integrationService) throw new HttpException(`Faceit not enabled on our backed. Please, make patience or contact us`, 404);

    body.data.game = body.data.game ?? 'csgo';

    let integration = await this.getIntegration(channelId);
    if (!integration) {
      integration = await this.prisma.channelIntegration.create({
        data: {
          channelId,
          enabled: body.enabled,
          data: { ...body.data } as unknown as Prisma.InputJsonObject,
          integrationId: integrationService.id,
        },
      });
    } else {
      integration = await this.prisma.channelIntegration.update({
        where: {
          id: integration.id,
        },
        data: {
          enabled: body.enabled,
          data: { ...body.data } as unknown as Prisma.InputJsonObject,
        },
      });
    }

    return integration;
  }
}

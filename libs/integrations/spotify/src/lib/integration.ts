import { PrismaClient } from '@tsuwari/prisma';

export class SpotifyIntegration {
  constructor(private readonly prisma: PrismaClient) {}

  async getSettings() {
    const service = await this.prisma.integration.findFirst({
      where: {
        service: 'SPOTIFY',
      },
    });

    if (service && service.clientId && service.clientSecret && service.redirectUrl) {
      return {
        id: service.id,
        clientId: service.clientId,
        clientSecret: service.clientSecret,
        redirectUrl: service.redirectUrl,
      };
    }

    return null;
  }
}
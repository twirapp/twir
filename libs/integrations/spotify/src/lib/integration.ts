import { Repository } from '@tsuwari/typeorm';
import { Integration, IntegrationService } from '@tsuwari/typeorm/entities/Integration';

export class SpotifyIntegration {
  constructor(private readonly repository: Repository<Integration>) {}

  async getSettings() {
    const service = await this.repository.findOneBy({
      service: IntegrationService.SPOTIFY,
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

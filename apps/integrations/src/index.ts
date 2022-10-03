import { AppDataSource, In } from '@tsuwari/typeorm';
import { ChannelIntegration } from '@tsuwari/typeorm/entities/ChannelIntegration';
import { IntegrationService } from '@tsuwari/typeorm/entities/Integration';

import { addDonationAlertsIntegration, DonationAlerts } from './services/donationAlerts.js';
import { addStreamlabsIntegration, StreamLabs } from './services/streamLabs.js';

export const typeorm = await AppDataSource.initialize();
export const donationAlertsStore: Map<string, DonationAlerts> = new Map();
export const streamlabsStore: Map<string, StreamLabs> = new Map();

const integrations = await typeorm.getRepository(ChannelIntegration).find({
  where: {
    integration: {
      service: In([IntegrationService.DONATIONALERTS, IntegrationService.STREAMLABS]),
    },
    enabled: true,
    channel: {
      isEnabled: true,
      isBanned: false,
      isTwitchBanned: false,
    },
  },
  relations: {
    integration: true,
  },
});

for (const integration of integrations) {
  if (integration.integration?.service === IntegrationService.DONATIONALERTS) {
    addDonationAlertsIntegration(integration).then((r) => {
      if (r) {
        donationAlertsStore.set(integration.channelId, r);
      }
    });
  }

  if (integration.integration?.service === IntegrationService.STREAMLABS) {
    addStreamlabsIntegration(integration).then((r) => {
      if (r) {
        streamlabsStore.set(integration.channelId, r);
      }
    });
  }
}

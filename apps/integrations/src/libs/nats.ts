import { config } from '@tsuwari/config';
import * as NatsIntegration from '@tsuwari/nats/integrations';
import { ChannelIntegration } from '@tsuwari/typeorm/entities/ChannelIntegration';
import { IntegrationService } from '@tsuwari/typeorm/entities/Integration';
import { connect } from 'nats';

import { donationAlertsStore, streamlabsStore, typeorm } from '../index.js';
import { addDonationAlertsIntegration } from '../services/donationAlerts.js';
import { addStreamlabsIntegration } from '../services/streamLabs.js';

export const nats = await connect({
  servers: [config.NATS_URL],
});

async function subscribeToAdd() {
  for await (const msg of nats.subscribe('integrations.add')) {
    const data = NatsIntegration.AddIntegration.fromBinary(msg.data);
    const integration = await typeorm.getRepository(ChannelIntegration).findOne({
      where: { id: data.id },
      relations: { integration: true },
    });

    if (!integration) {
      msg.respond(new Uint8Array());
      continue;
    }
    await removeIntegration(integration);
    console.info(`Adding ${integration.id} connection`);
    if (integration.integration?.service === IntegrationService.DONATIONALERTS) {
      const instance = await addDonationAlertsIntegration(integration);
      if (instance) {
        donationAlertsStore.set(integration.channelId, instance);
      }
    }
    if (integration.integration?.service === IntegrationService.STREAMLABS) {
      const instance = await addStreamlabsIntegration(integration);
      if (instance) {
        streamlabsStore.set(integration.channelId, instance);
      }
    }
    msg.respond(new Uint8Array());
  }
}
async function subscribeToRemove() {
  for await (const msg of nats.subscribe('integrations.remove')) {
    const data = NatsIntegration.RemoveIntegration.fromBinary(msg.data);

    const integration = await typeorm.getRepository(ChannelIntegration).findOne({
      where: { id: data.id },
      relations: { integration: true },
    });

    if (!integration) {
      msg.respond(new Uint8Array());
      continue;
    }
    console.info(`Destroying ${integration.id} connection`);
    await removeIntegration(integration);
    msg.respond(new Uint8Array());
  }
}
subscribeToAdd();
subscribeToRemove();

async function removeIntegration(integration: ChannelIntegration) {
  if (integration.integration?.service === IntegrationService.STREAMLABS) {
    const existed = streamlabsStore.get(integration.channelId);
    if (!existed) return;
    await existed.destroy();
    streamlabsStore.delete(integration.channelId);
  }

  if (integration.integration?.service === IntegrationService.DONATIONALERTS) {
    const existed = donationAlertsStore.get(integration.channelId);
    if (!existed) return;
    await existed.destroy();
    donationAlertsStore.delete(integration.channelId);
  }
}

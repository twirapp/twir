import { AppDataSource, In } from '@tsuwari/typeorm';
import { ChannelIntegration } from '@tsuwari/typeorm/entities/ChannelIntegration';
import { IntegrationService } from '@tsuwari/typeorm/entities/Integration';

import { DonationAlerts } from './services/donationAlerts.js';

export const typeorm = await AppDataSource.initialize();
const donationAlertsStore: Map<string, DonationAlerts> = new Map();

const integrations = await typeorm.getRepository(ChannelIntegration).find({
  where: {
    integration: {
      service: In([IntegrationService.DONATIONALERTS, IntegrationService.STREAMLABS]),
    },
    enabled: true,
  },
  relations: {
    integration: true,
  },
});

for (const integration of integrations) {
  if (integration.integration?.service === IntegrationService.DONATIONALERTS) {
    addDonationAlertsIntegration(integration);
  }
}

async function addDonationAlertsIntegration(integration: ChannelIntegration) {
  if (
    !integration.accessToken ||
    !integration.refreshToken ||
    !integration.integration ||
    !integration.integration.clientId ||
    !integration.integration.clientSecret
  ) {
    return;
  }

  const refresh = await fetch('https://www.donationalerts.com/oauth/token', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded',
    },
    body: new URLSearchParams({
      grant_type: 'refresh_token',
      refresh_token: integration.refreshToken,
      client_id: integration.integration.clientId,
      client_secret: integration.integration.clientSecret,
    }).toString(),
  });

  if (!refresh.ok) {
    console.error(await refresh.text());
    return;
  }

  const refreshResponse = await refresh.json();

  await typeorm
    .getRepository(ChannelIntegration)
    .update(
      { id: integration.id },
      { accessToken: refreshResponse.access_token, refreshToken: refreshResponse.refresh_token },
    );

  const request = await fetch('https://www.donationalerts.com/api/v1/user/oauth', {
    headers: {
      Authorization: `Bearer ${refreshResponse.access_token}`,
    },
  });

  if (!request.ok) {
    console.log(await request.text(), refreshResponse);
    return;
  }

  const { data } = await request.json();
  const { id, socket_connection_token } = data;
  const instance = new DonationAlerts(
    refreshResponse.access_token,
    id,
    socket_connection_token,
    integration.channelId,
  );
  await instance.init();

  donationAlertsStore.set(integration.channelId, instance);
}

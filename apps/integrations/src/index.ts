import { PORTS } from '@twir/grpc/constants/constants';
import * as Integrations from '@twir/grpc/generated/integrations/integrations';
import { createServer } from 'nice-grpc';

import { DONATEPAY, DONATIONALERTS, getIntegrations, STREAMLABS } from './libs/db.js';
import { addDonatePayIntegration, DonatePay } from './services/donatepay.js';
import { addDonationAlertsIntegration, DonationAlerts } from './services/donationAlerts.js';
import { addStreamlabsIntegration, StreamLabs } from './services/streamLabs.js';
import { Integration } from './types.js';
import './pubsub.js';

export const donationAlertsStore: Map<string, DonationAlerts> = new Map();
export const streamlabsStore: Map<string, StreamLabs> = new Map();
export const donatePayStore: Map<string, DonatePay> = new Map();

const integrations = await getIntegrations();

for (const integration of integrations) {
	if (integration.integration.service === DONATIONALERTS) {
		addDonationAlertsIntegration(integration).then((r) => {
			if (r) {
				donationAlertsStore.set(integration.channelId, r);
			}
		});
	}

	if (integration.integration.service === STREAMLABS) {
		addStreamlabsIntegration(integration).then((r) => {
			if (r) {
				streamlabsStore.set(integration.channelId, r);
			}
		});
	}

	if (integration.integration.service === DONATEPAY) {
		addDonatePayIntegration(integration).then((r) => {
			if (r) {
				donatePayStore.set(integration.channelId, r);
			}
		});
	}
}

const integrationsServer: Integrations.IntegrationsServiceImplementation = {
	async addIntegration(data: Integrations.Request) {
		const integration = await getIntegrations(data.id);

		if (!integration) {
			return {};
		}
		await removeIntegration(integration);
		console.info(`Adding ${integration.id} connection`);
		if (integration.integration.service === DONATIONALERTS) {
			const instance = await addDonationAlertsIntegration(integration);
			if (instance) {
				donationAlertsStore.set(integration.channelId, instance);
			}
		}
		if (integration.integration.service === STREAMLABS) {
			const instance = await addStreamlabsIntegration(integration);
			if (instance) {
				streamlabsStore.set(integration.channelId, instance);
			}
		}
		if (integration.integration.service === DONATEPAY) {
			const instance = await addDonatePayIntegration(integration);
			if (instance) {
				donatePayStore.set(integration.channelId, instance);
			}
		}

		return {};
	},
	async removeIntegration(data: Integrations.Request) {
		const integration = await getIntegrations(data.id);

		if (!integration) {
			return {};
		}

		console.info(`Destroying ${integration.id} connection`);
		await removeIntegration(integration);
		return {};
	},
};

const server = createServer({
	'grpc.keepalive_time_ms': 1 * 60 * 1000,
});

server.add(Integrations.IntegrationsDefinition, integrationsServer);

await server.listen(`0.0.0.0:${PORTS.INTEGRATIONS_SERVER_PORT}`);
console.info('Integrations started');

export async function removeIntegration(integration: Integration) {
	if (integration.integration.service === STREAMLABS) {
		const existed = streamlabsStore.get(integration.channelId);
		if (!existed) return;
		await existed.destroy();
		streamlabsStore.delete(integration.channelId);
	}

	if (integration.integration.service === DONATIONALERTS) {
		const existed = donationAlertsStore.get(integration.channelId);
		if (!existed) return;
		await existed.destroy();
		donationAlertsStore.delete(integration.channelId);
	}

	if (integration.integration.service === DONATEPAY) {
		const existed = donatePayStore.get(integration.channelId);
		if (!existed) return;
		await existed.disconnect();
		donatePayStore.delete(integration.channelId);
	}
}

process.on('uncaughtException', console.error);
process.on('unhandledRejection', console.error);

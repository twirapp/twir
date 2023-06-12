import * as Integrations from '@tsuwari/grpc/generated/integrations/integrations';
import { PORTS } from '@tsuwari/grpc/servers/constants';
import { AppDataSource, In } from '@tsuwari/typeorm';
import { ChannelIntegration } from '@tsuwari/typeorm/entities/ChannelIntegration';
import { IntegrationService } from '@tsuwari/typeorm/entities/Integration';
import { createServer } from 'nice-grpc';

import { addDonatePayIntegration, DonatePay } from './services/donatepay.js';
import { addDonationAlertsIntegration, DonationAlerts } from './services/donationAlerts.js';
import { addStreamlabsIntegration, StreamLabs } from './services/streamLabs.js';

export const typeorm = await AppDataSource.initialize();
export const donationAlertsStore: Map<string, DonationAlerts> = new Map();
export const streamlabsStore: Map<string, StreamLabs> = new Map();
export const donatePayStore: Map<string, DonatePay> = new Map();


const integrations = await typeorm.getRepository(ChannelIntegration).find({
	where: {
		integration: {
			service: In([
				IntegrationService.DONATIONALERTS,
				IntegrationService.STREAMLABS,
				IntegrationService.DONATEPAY,
				IntegrationService.DONATELLO,
			]),
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

	if (integration.integration?.service === IntegrationService.DONATEPAY) {
		addDonatePayIntegration(integration).then((r) => {
			if (r) {
				donatePayStore.set(integration.channelId, r);
			}
		});
	}
}

const integrationsServer: Integrations.IntegrationsServiceImplementation = {
	async addIntegration(data: Integrations.Request) {
		const integration = await typeorm.getRepository(ChannelIntegration).findOne({
			where: {
				id: data.id,
				channel: {
					isBanned: false,
					isEnabled: true,
					isTwitchBanned: false,
				},
			},
			relations: { integration: true },
		});

		if (!integration) {
			return {};
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
		if (integration.integration?.service === IntegrationService.DONATEPAY) {
			const instance = await addDonatePayIntegration(integration);
			if (instance) {
				donatePayStore.set(integration.channelId, instance);
			}
		}

		return {};
	},
	async removeIntegration(data: Integrations.Request) {
		const integration = await typeorm.getRepository(ChannelIntegration).findOne({
			where: { id: data.id },
			relations: { integration: true },
		});

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

export async function removeIntegration(integration: ChannelIntegration) {
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

	if (integration.integration?.service === IntegrationService.DONATEPAY) {
		const existed = donatePayStore.get(integration.channelId);
		if (!existed) return;
		await existed.disconnect();
		donatePayStore.delete(integration.channelId);
	}
}

process.on('uncaughtException', console.error);
process.on('unhandledRejection', console.error);

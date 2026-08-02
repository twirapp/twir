import { IntegrationService } from '@twir/bus-core';
import { config } from '@twir/config';
import process from 'node:process';

import {
	getDonationAlertsIntegrations,
	getDonationPayIntegrations,
	getStreamElementsIntegrations,
	getStreamlabsIntegrations,
	providerTokenStores,
} from './libs/db';
import { StreamElementsClient } from './libs/streamelements-client.ts';
import { twirBus } from './libs/twirbus.ts';
import { StreamElementsConnection } from './services/streamElements.ts';
import {
	addIntegration as addDonatePayIntegration,
	removeIntegration as removeDonatePayIntegration,
} from './store/donatePay';
import {
	addIntegration as addDonationAlertsIntegration,
	removeIntegration as removeDonationAlertsIntegration,
} from './store/donationAlerts';
import {
	addIntegration as addStreamlabsIntegration,
	removeIntegration as removeStreamlabsIntegration,
} from './store/streamlabs';
import { createStreamElementsStore, runLifecycleOperation } from './store/streamelements.ts';
import './pubsub';

const streamElementsStore = createStreamElementsStore({
	loadIntegrationByID: (integrationID) => getStreamElementsIntegrations({ id: integrationID }),
	createConnection(integration) {
		if (!config.STREAM_ELEMENTS_CLIENT_ID || !config.STREAM_ELEMENTS_CLIENT_SECRET) {
			throw new Error('StreamElements OAuth credentials are not configured');
		}
		if (!integration.accessToken || !integration.refreshToken) {
			throw new Error('StreamElements integration tokens are missing');
		}
		const client = new StreamElementsClient({
			channelID: integration.channelId,
			tokens: {
				accessToken: integration.accessToken,
				refreshToken: integration.refreshToken,
			},
			tokenStore: providerTokenStores.streamElements,
			clientID: config.STREAM_ELEMENTS_CLIENT_ID,
			clientSecret: config.STREAM_ELEMENTS_CLIENT_SECRET,
		});
		const connection = new StreamElementsConnection({
			channelID: integration.channelId,
			client,
		});
		connection.connect();
		return connection;
	},
});

for (const donatePayIntegration of await getDonationPayIntegrations()) {
	addDonatePayIntegration(donatePayIntegration);
}

for (const integration of await getDonationAlertsIntegrations()) {
	addDonationAlertsIntegration(integration);
}

for (const integration of await getStreamlabsIntegrations()) {
	await addStreamlabsIntegration(integration);
}

for (const integration of await getStreamElementsIntegrations()) {
	await streamElementsStore.addIntegration(integration);
}

interface IntegrationLifecycleRequest {
	readonly id: string
	readonly service: IntegrationService
}

async function handleIntegrationAdd(data: IntegrationLifecycleRequest): Promise<void> {
	console.info(`Adding ${data.id} (${data.service}) connection`);

	if (data.service === IntegrationService.DONATEPAY) {
		const integration = await getDonationPayIntegrations({ id: data.id });
		if (!integration) {
			console.error(`Integration with id ${data.id} not found for DonatePay`);
			return;
		}
		await addDonatePayIntegration(integration);
		return;
	}

	if (data.service === IntegrationService.DONATIONALERTS) {
		const integration = await getDonationAlertsIntegrations({ id: Number(data.id) });
		if (!integration) {
			console.error(`Integration with id ${data.id} not found for DonateAlerts`);
			return;
		}
		await addDonationAlertsIntegration(integration);
		return;
	}

	if (data.service === IntegrationService.STREAMLABS) {
		const integration = await getStreamlabsIntegrations({ id: data.id });
		if (!integration) {
			console.error(`Integration with id ${data.id} not found for Streamlabs`);
			return;
		}
		await addStreamlabsIntegration(integration);
		return;
	}

	if (data.service === IntegrationService.STREAMELEMENTS) {
		await streamElementsStore.addIntegrationByID(data.id);
		return;
	}
}

async function handleIntegrationRemove(data: IntegrationLifecycleRequest): Promise<void> {
	console.info(`Destroying ${data.id} (${data.service}) connection`);

	if (data.service === IntegrationService.DONATEPAY) {
		await removeDonatePayIntegration(data.id); // channelId
		return;
	}

	if (data.service === IntegrationService.DONATIONALERTS) {
		await removeDonationAlertsIntegration(data.id); // channelId
		return;
	}

	if (data.service === IntegrationService.STREAMLABS) {
		await removeStreamlabsIntegration(data.id); // channelId
		return;
	}

	if (data.service === IntegrationService.STREAMELEMENTS) {
		await streamElementsStore.removeIntegration(data.id); // channelId
		return;
	}
}

twirBus.Integrations.Add.subscribe((data) => {
	runLifecycleOperation(() => handleIntegrationAdd(data));
	return null;
});

twirBus.Integrations.Remove.subscribe((data) => {
	runLifecycleOperation(() => handleIntegrationRemove(data));
	return null;
});

console.info('Integrations started');

process.on('uncaughtException', console.error);
process.on('unhandledRejection', console.error);

let shuttingDown = false;
const shutdown = async () => {
	if (shuttingDown) return;
	shuttingDown = true;
	await streamElementsStore.closeAll();
	process.exit(0);
};

process.on('SIGTERM', () => { void shutdown(); });
process.on('SIGINT', () => { void shutdown(); });

import process from 'node:process'

import { PORTS } from '@twir/grpc/constants/constants'
import * as Integrations from '@twir/grpc/integrations/integrations'
import { createServer } from 'nice-grpc'

import { Service, getIntegrations } from './libs/db'
import {
	addIntegration as addDonatePayIntegration,
	removeIntegration as removeDonatePayIntegration,
} from './store/donatePay'
import {
	addIntegration as addDonationAlertsIntegration,
	removeIntegration as removeDonationAlertsIntegration,
} from './store/donationAlerts'
import {
	addIntegration as addStreamlabsIntegration,
	removeIntegration as removeStreamlabsIntegration,
} from './store/streamlabs'

import type { Integration } from './libs/db'
import type { IntegrationsServiceImplementation } from '@twir/grpc/integrations/integrations'

import './pubsub'

const integrations = await getIntegrations()

for (const integration of integrations) {
	if (integration.integration.service === Service.DONATIONALERTS) {
		await addDonationAlertsIntegration(integration)
	}

	if (integration.integration.service === Service.STREAMLABS) {
		await addStreamlabsIntegration(integration)
	}

	if (integration.integration.service === Service.DONATEPAY) {
		await addDonatePayIntegration(integration)
	}
}

const integrationsServer: IntegrationsServiceImplementation = {
	async addIntegration(data) {
		const integration = await getIntegrations(data.id)

		if (!integration) {
			return {}
		}

		console.info(`Adding ${integration.id} connection`)

		if (integration.integration.service === Service.DONATIONALERTS) {
			await addDonationAlertsIntegration(integration)
		}
		if (integration.integration.service === Service.STREAMLABS) {
			await addStreamlabsIntegration(integration)
		}
		if (integration.integration.service === Service.DONATEPAY) {
			await addDonatePayIntegration(integration)
		}

		return {}
	},

	async removeIntegration(data) {
		const integration = await getIntegrations(data.id)

		if (!integration) {
			return {}
		}

		console.info(`Destroying ${integration.id} connection`)
		await removeIntegration(integration)
		return {}
	},
}

const server = createServer({
	'grpc.keepalive_time_ms': 1 * 60 * 1000,
})

server.add(Integrations.IntegrationsDefinition, integrationsServer)

await server.listen(`0.0.0.0:${PORTS.INTEGRATIONS_SERVER_PORT}`)
console.info('Integrations started')

export async function removeIntegration(integration: Integration) {
	if (integration.integration.service === Service.STREAMLABS) {
		await removeStreamlabsIntegration(integration.channelId)
	}

	if (integration.integration.service === Service.DONATIONALERTS) {
		await removeDonationAlertsIntegration(integration.channelId)
	}

	if (integration.integration.service === Service.DONATEPAY) {
		await removeDonatePayIntegration(integration.channelId)
	}
}

process.on('uncaughtException', console.error)
process.on('unhandledRejection', console.error)

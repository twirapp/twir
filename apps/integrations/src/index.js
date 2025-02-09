import process from 'node:process'

import { PORTS } from '@twir/grpc/constants/constants'
import * as Integrations from '@twir/grpc/integrations/integrations'
import { createServer } from 'nice-grpc'

import { Services, getIntegrations } from './libs/db.js'
import {
	addIntegration as addDonatePayIntegration,
	removeIntegration as removeDonatePayIntegration,
} from './store/donatePay.js'
import {
	addIntegration as addDonationAlertsIntegration,
	removeIntegration as removeDonationAlertsIntegration,
} from './store/donationAlerts.js'
import {
	addIntegration as addStreamlabsIntegration,
	removeIntegration as removeStreamlabsIntegration,
} from './store/streamlabs.js'

import './pubsub.js'

const integrations = await getIntegrations()

for (const integration of integrations) {
	if (integration.integration.service === Services.DONATIONALERTS) {
		await addDonationAlertsIntegration(integration)
	}

	if (integration.integration.service === Services.STREAMLABS) {
		await addStreamlabsIntegration(integration)
	}

	if (integration.integration.service === Services.DONATEPAY) {
		await addDonatePayIntegration(integration)
	}
}

/**
 * @type {import('@twir/grpc/integrations/integrations').IntegrationsServiceImplementation}
 */
const integrationsServer = {
	async addIntegration(data) {
		const integration = await getIntegrations(data.id)

		if (!integration) {
			return {}
		}

		console.info(`Adding ${integration.id} connection`)

		if (integration.integration.service === Services.DONATIONALERTS) {
			addDonationAlertsIntegration(integration)
		}
		if (integration.integration.service === Services.STREAMLABS) {
			addStreamlabsIntegration(integration)
		}
		if (integration.integration.service === Services.DONATEPAY) {
			addDonatePayIntegration(integration)
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

/**
 * @param {Integration} integration  Options object for each OS, and global options.
 */
export async function removeIntegration(integration) {
	if (integration.integration.service === Services.STREAMLABS) {
		removeStreamlabsIntegration(integration.channelId)
	}

	if (integration.integration.service === Services.DONATIONALERTS) {
		removeDonationAlertsIntegration(integration.channelId)
	}

	if (integration.integration.service === Services.DONATEPAY) {
		removeDonatePayIntegration(integration.channelId)
	}
}

process.on('uncaughtException', console.error)
process.on('unhandledRejection', console.error)

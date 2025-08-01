import process from 'node:process'

import { Service, getDonationPayIntegrations, getIntegrations } from './libs/db'
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

import './pubsub'
import { twirBus } from './libs/twirbus.ts'

const integrations = await getIntegrations()

for (const integration of integrations) {
	if (integration.integration.service === Service.DONATIONALERTS) {
		addDonationAlertsIntegration(integration)
	}

	if (integration.integration.service === Service.STREAMLABS) {
		addStreamlabsIntegration(integration)
	}
}

for (const donatePayIntegration of await getDonationPayIntegrations()) {
	addDonatePayIntegration(donatePayIntegration)
}

twirBus.Integrations.Add.subscribe(async (data) => {
	console.info(`Adding ${data.id} connection`)
	const integration = await getIntegrations(data.id)

	if (!integration) {
		return null
	}

	if (integration.integration.service === Service.DONATIONALERTS) {
		await addDonationAlertsIntegration(integration)
	}
	if (integration.integration.service === Service.STREAMLABS) {
		await addStreamlabsIntegration(integration)
	}
	if (integration.integration.service === Service.DONATEPAY) {
		await addDonatePayIntegration(integration)
	}

	return null
})

twirBus.Integrations.Remove.subscribe(async (data) => {
	console.info(`Destroying ${data.id} connection`)
	const integration = await getIntegrations(data.id)

	if (!integration) {
		return null
	}

	await removeIntegration(integration)
	return null
})

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

console.info('Integrations started')

process.on('uncaughtException', console.error)
process.on('unhandledRejection', console.error)

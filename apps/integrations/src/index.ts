import process from 'node:process'

import { IntegrationService } from '@twir/bus-core'

import { Service, getDonationPayIntegrations, getIntegrations } from './libs/db'
import { twirBus } from './libs/twirbus.ts'
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
	await addDonatePayIntegration(donatePayIntegration)
}

twirBus.Integrations.Add.subscribe(async (data) => {
	console.info(`Adding ${data.id} (${data.service}) connection`)

	if (data.service === IntegrationService.DONATEPAY) {
		const integration = await getDonationPayIntegrations({ id: data.id })
		if (!integration) {
			console.error(`Integration with id ${data.id} not found for DonatePay`)
			return null
		}
		await addDonatePayIntegration(integration)
	}

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

	return null
})

twirBus.Integrations.Remove.subscribe(async (data) => {
	console.info(`Destroying ${data.id} (${data.service}) connection`)

	if (data.service === IntegrationService.DONATEPAY) {
		const integration = await getDonationPayIntegrations({ id: data.id })
		if (!integration) {
			console.error(`Integration with id ${data.id} not found for DonatePay`)
			return null
		}
		await removeDonatePayIntegration(integration.channel_id)
		return null
	}

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
}

console.info('Integrations started')

process.on('uncaughtException', console.error)
process.on('unhandledRejection', console.error)

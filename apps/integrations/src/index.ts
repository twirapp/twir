import { IntegrationService } from '@twir/bus-core'
import process from 'node:process'

import type { Integration } from './libs/db'

import {
	Service,
	getDonationAlertsIntegrations,
	getDonationPayIntegrations,
	getIntegrations,
	getStreamlabsIntegrations,
} from './libs/db'
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
import './pubsub'

for (const donatePayIntegration of await getDonationPayIntegrations()) {
	addDonatePayIntegration(donatePayIntegration)
}

for (const integration of await getDonationAlertsIntegrations()) {
	addDonationAlertsIntegration(integration)
}

for (const integration of await getStreamlabsIntegrations()) {
	addStreamlabsIntegration(integration)
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
		return
	}

	if (data.service === IntegrationService.DONATIONALERTS) {
		const integration = await getDonationAlertsIntegrations({ id: Number(data.id) })
		if (!integration) {
			console.error(`Integration with id ${data.id} not found for DonateAlerts`)
			return null
		}
		await addDonationAlertsIntegration(integration)
		return
	}

	if (data.service === IntegrationService.STREAMLABS) {
		const integration = await getStreamlabsIntegrations({ id: data.id })
		if (!integration) {
			console.error(`Integration with id ${data.id} not found for Streamlabs`)
			return null
		}
		await addStreamlabsIntegration(integration)
		return
	}

	return null
})

twirBus.Integrations.Remove.subscribe(async (data) => {
	console.info(`Destroying ${data.id} (${data.service}) connection`)

	if (data.service === IntegrationService.DONATEPAY) {
		await removeDonatePayIntegration(data.id) // channelId
		return null
	}

	if (data.service === IntegrationService.DONATIONALERTS) {
		await removeDonationAlertsIntegration(data.id) // channelId
		return null
	}

	if (data.service === IntegrationService.STREAMLABS) {
		await removeStreamlabsIntegration(data.id) // channelId
		return null
	}

	return null
})

console.info('Integrations started')

process.on('uncaughtException', console.error)
process.on('unhandledRejection', console.error)

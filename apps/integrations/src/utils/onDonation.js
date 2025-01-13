import { randomUUID } from 'node:crypto'

import { db } from '../libs/db.js'
import { eventsGrpcClient } from '../libs/eventsGrpc.js'

/**
 * @param {Donate} donate
 */
export async function onDonation(donate) {
	const userName = donate.userName ?? 'Anonymous'

	await db.insert({
		id: randomUUID(),
		channel_id: donate.twitchUserId,
		type: 'DONATION',
		data: {
			donationAmount: donate.amount.toString(),
			donationCurrency: donate.currency,
			donationMessage: donate.message,
			donationUsername: userName,
		},
	}).into('channels_events_list')

	const msg = donate.message || ''

	await eventsGrpcClient.donate({
		amount: donate.amount.toString(),
		message: msg,
		currency: donate.currency,
		baseInfo: { channelId: donate.twitchUserId },
		userName,
	})
}

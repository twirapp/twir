import { insertDonation } from '../libs/db.ts'
import { twirBus } from '../libs/twirbus.ts'

export interface Donate {
	twitchUserId: string
	amount: number | string
	currency: string
	message?: string | null
	userName?: string | null
}

export async function onDonation(donate: Donate) {
	const userName = donate.userName ?? 'Anonymous'

	await insertDonation(donate)

	const msg = donate.message || ''

	await twirBus.Events.Donate.publish({
		base_info: { channel_id: donate.twitchUserId, channel_name: '' },
		user_name: userName,
		amount: donate.amount.toString(),
		currency: donate.currency,
		message: msg,
	})
}

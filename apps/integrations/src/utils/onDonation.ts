import { insertDonation } from '../libs/db.ts'
import { eventsGrpcClient } from '../libs/eventsGrpc.ts'

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

	await eventsGrpcClient.donate({
		amount: donate.amount.toString(),
		message: msg,
		currency: donate.currency,
		baseInfo: { channelId: donate.twitchUserId },
		userName,
	})
}

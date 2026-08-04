const DEDUPE_TTL_SECONDS = 86_400

export interface DonationDedupeRedis {
	send(command: string, args: readonly string[]): Promise<unknown>
}

export async function claimDonation(
	provider: string,
	eventID: string,
	redis?: DonationDedupeRedis
): Promise<boolean> {
	if (!eventID) {
		return true
	}

	const commands = redis ?? (await import('./redis.ts')).client
	const result = await commands.send('SET', [
		`twir:donation:${provider}:${eventID}`,
		'1',
		'NX',
		'EX',
		String(DEDUPE_TTL_SECONDS),
	])
	return result === 'OK'
}

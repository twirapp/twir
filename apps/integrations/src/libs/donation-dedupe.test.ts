import { expect, test } from 'bun:test'

import { type DonationDedupeRedis, claimDonation } from './donation-dedupe.ts'

test('claims a provider event for 24 hours with Redis NX', async () => {
	const calls: Array<{ command: string; args: readonly string[] }> = []
	const redis: DonationDedupeRedis = {
		async send(command, args) {
			calls.push({ command, args })
			return 'OK'
		},
	}

	expect(await claimDonation('streamelements', 'event-1', redis)).toBe(true)
	expect(calls).toEqual([{
		command: 'SET',
		args: ['twir:donation:streamelements:event-1', '1', 'NX', 'EX', '86400'],
	}])
})

test('rejects an event already claimed by another process', async () => {
	const redis: DonationDedupeRedis = {
		async send() {
			return null
		},
	}

	expect(await claimDonation('streamlabs', 'event-2', redis)).toBe(false)
})

test('passes through missing event ids without Redis access', async () => {
	let called = false
	const redis: DonationDedupeRedis = {
		async send() {
			called = true
			return null
		},
	}

	expect(await claimDonation('streamlabs', '', redis)).toBe(true)
	expect(called).toBe(false)
})

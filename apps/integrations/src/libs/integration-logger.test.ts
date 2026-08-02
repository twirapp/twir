import { expect, test } from 'bun:test'

import { StreamLabsClientError } from './streamlabs-client.ts'
import {
	type IntegrationErrorDetails,
	logIntegrationError,
} from './integration-logger.ts'

test('logs allowlisted typed errors with provider and channel context', () => {
	let details: IntegrationErrorDetails | undefined
	logIntegrationError(
		{ provider: 'streamlabs', operation: 'hydrate', channelID: 'channel-1' },
		new StreamLabsClientError(),
		(_message, value) => { details = value },
	)

	expect(details).toEqual({
		provider: 'streamlabs',
		operation: 'hydrate',
		channelID: 'channel-1',
		category: 'StreamLabsClientError',
		message: 'Streamlabs provider request failed',
	})
})

test('redacts messages from unknown errors', () => {
	let serialized = ''
	logIntegrationError(
		{ provider: 'streamlabs', operation: 'add', integrationID: 'integration-1' },
		new Error('provider body contains client_secret=secret'),
		(message, details) => { serialized = JSON.stringify({ message, details }) },
	)

	expect(serialized).not.toContain('client_secret')
	expect(serialized).not.toContain('secret')
	expect(serialized).toContain('UnexpectedError')
})

test('does not trust a spoofed allowlisted error name', () => {
	let serialized = ''
	const spoofed = new Error('secret provider response')
	spoofed.name = 'StreamLabsClientError'
	logIntegrationError(
		{ provider: 'streamlabs', operation: 'hydrate', channelID: 'channel-1' },
		spoofed,
		(message, details) => { serialized = JSON.stringify({ message, details }) },
	)

	expect(serialized).not.toContain('secret provider response')
	expect(serialized).toContain('UnexpectedError')
})

import {
	OAuthRefreshLockLostError,
	OAuthRefreshLockUnavailableError,
} from './oauth-lock.ts'
import {
	InvalidProviderTokensError,
	ProviderTokensNotFoundError,
} from './provider-token-store.ts'
import { StreamElementsRefreshError } from './streamelements-client.ts'
import { StreamLabsClientError } from './streamlabs-client.ts'

export type IntegrationOperation = 'add' | 'remove' | 'hydrate' | 'donation' | 'process'

export interface IntegrationLogContext {
	readonly provider: string
	readonly operation: IntegrationOperation
	readonly channelID?: string
	readonly integrationID?: string
}

export interface IntegrationErrorDetails extends IntegrationLogContext {
	readonly category: string
	readonly message: string
}

export type IntegrationLogWriter = (
	message: string,
	details: IntegrationErrorDetails
) => void

function isSafeTypedError(error: unknown): error is Error {
	return error instanceof InvalidProviderTokensError
		|| error instanceof OAuthRefreshLockLostError
		|| error instanceof OAuthRefreshLockUnavailableError
		|| error instanceof ProviderTokensNotFoundError
		|| error instanceof StreamElementsRefreshError
		|| error instanceof StreamLabsClientError
}

export function logIntegrationError(
	context: IntegrationLogContext,
	error: unknown,
	write: IntegrationLogWriter = (message, details) => { console.error(message, details) }
): void {
	const safeError = isSafeTypedError(error)
	write('Integration operation failed', {
		...context,
		category: safeError ? error.name : 'UnexpectedError',
		message: safeError ? error.message : 'Integration operation failed',
	})
}

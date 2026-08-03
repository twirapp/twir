import { describe, expect, it, vi } from 'vitest'

import { parseMcpConsentAttempt } from '~/utils/mcp-consent.js'

import { createMcpConsentApi } from './api.js'

function getAttempt(value: string) {
	const attempt = parseMcpConsentAttempt(value)
	if (attempt === null) throw new Error('expected valid attempt fixture')
	return attempt
}

describe('MCP consent API', () => {
	it('loads and parses a consent request for a valid opaque attempt', async () => {
		const fetchMock = vi.fn().mockResolvedValue({
			client: { id: 'client-1', name: 'Desktop Agent', uri: 'https://agent.example' },
			channel_id: 'channel-1',
			requested_scopes: ['read', 'write'],
			access_levels: ['read', 'write'],
			csrf_token: 'csrf-token',
		})
		const api = createMcpConsentApi(fetchMock)
		const result = await api.getMcpConsent(getAttempt('opaque-attempt-123'))

		expect(fetchMock).toHaveBeenCalledWith('/api/oauth/consent', {
			query: { attempt: 'opaque-attempt-123' },
		})
		expect(result).toEqual({
			kind: 'success',
			data: {
				client: { id: 'client-1', name: 'Desktop Agent', uri: 'https://agent.example' },
				channel_id: 'channel-1',
				requested_scopes: ['read', 'write'],
				access_levels: ['read', 'write'],
				csrf_token: 'csrf-token',
			},
		})
	})

	it('classifies an expired attempt without exposing a backend error', async () => {
		const api = createMcpConsentApi(
			vi.fn().mockRejectedValue(Object.assign(new Error('expired attempt'), { status: 410 })),
		)

		const result = await api.getMcpConsent(getAttempt('opaque-attempt-123'))

		expect(result).toEqual({ kind: 'expired' })
	})

	it('requires the consent page to reload when its dashboard changed', async () => {
		const api = createMcpConsentApi(
			vi.fn().mockRejectedValue(Object.assign(new Error('dashboard changed'), { status: 409 })),
		)

		const result = await api.getMcpConsent(getAttempt('opaque-attempt-123'))

		expect(result).toEqual({ kind: 'expired' })
	})

	it('rejects a write-only response before it can initialize an invalid read default', async () => {
		const api = createMcpConsentApi(
			vi.fn().mockResolvedValue({
				client: { id: 'client-1', name: 'Desktop Agent' },
				channel_id: 'channel-1',
				requested_scopes: ['write'],
				access_levels: ['write'],
				csrf_token: 'csrf-token',
			}),
		)

		const result = await api.getMcpConsent(getAttempt('opaque-attempt-123'))

		expect(result).toEqual({ kind: 'network' })
	})

	it('posts an approval with the consent response channel and backend-only redirect target', async () => {
		const fetchMock = vi.fn().mockResolvedValue({ redirect_to: 'https://agent.example/callback' })
		const api = createMcpConsentApi(fetchMock)
		const result = await api.submitMcpConsent({
			attempt: getAttempt('opaque-attempt-123'),
			csrf_token: 'csrf-token',
			channel_id: 'channel-1',
			decision: 'approve',
			access_level: 'read',
		})

		expect(fetchMock).toHaveBeenCalledWith('/api/oauth/consent', {
			method: 'POST',
			body: {
				attempt: 'opaque-attempt-123',
				csrf_token: 'csrf-token',
				channel_id: 'channel-1',
				decision: 'approve',
				access_level: 'read',
			},
		})
		expect(result).toEqual({
			kind: 'success',
			data: { redirectTo: 'https://agent.example/callback' },
		})
	})

	it('posts a denial with the consent response channel', async () => {
		const fetchMock = vi.fn().mockResolvedValue({ redirect_to: 'https://agent.example/callback' })
		const api = createMcpConsentApi(fetchMock)
		const result = await api.submitMcpConsent({
			attempt: getAttempt('opaque-attempt-123'),
			csrf_token: 'csrf-token',
			channel_id: 'channel-1',
			decision: 'deny',
		})

		expect(fetchMock).toHaveBeenCalledWith('/api/oauth/consent', {
			method: 'POST',
			body: {
				attempt: 'opaque-attempt-123',
				csrf_token: 'csrf-token',
				channel_id: 'channel-1',
				decision: 'deny',
			},
		})
		expect(result).toEqual({
			kind: 'success',
			data: { redirectTo: 'https://agent.example/callback' },
		})
	})
})

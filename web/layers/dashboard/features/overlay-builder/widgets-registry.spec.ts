import { describe, expect, it, vi } from 'vitest'

import { overlayWidgetRegistry } from './widgets-registry'

vi.mock('./components/widget-settings/ChatWidgetSettings.vue', () => ({ default: {} }))
vi.mock('./components/widget-settings/DudesWidgetSettings.vue', () => ({ default: {} }))
vi.mock('./components/widget-settings/FaceitStatsWidgetSettings.vue', () => ({ default: {} }))
vi.mock('./components/widget-settings/NowPlayingWidgetSettings.vue', () => ({ default: {} }))
vi.mock('./components/widget-settings/ValorantStatsWidgetSettings.vue', () => ({ default: {} }))

const context = {
	origin: 'https://twir.localhost',
	apiKey: 'channel-key',
}

function getWidget(key: string) {
	const widget = overlayWidgetRegistry.find((entry) => entry.key === key)
	if (!widget) throw new Error(`Missing widget registry entry: ${key}`)
	return widget
}

describe('overlay widget registry', () => {
	it('builds the six requested runtime URLs', () => {
		expect(getWidget('now-playing').buildUrl(context)).toBe(
			'https://twir.localhost/overlays/channel-key/now-playing',
		)
		expect(getWidget('kappagen').buildUrl(context)).toBe('https://twir.localhost/overlays/channel-key/kappagen')
		expect(getWidget('dudes').buildUrl({ ...context, params: { id: 'dudes-preset' } })).toBe(
			'https://twir.localhost/overlays/channel-key/dudes?id=dudes-preset',
		)
		expect(getWidget('afk').buildUrl(context)).toBe('https://twir.localhost/overlays/channel-key/brb')
		expect(getWidget('faceit-stats').buildUrl({ ...context, params: { nickname: 'qa player' } })).toContain(
		'https://twir.localhost/overlays/faceit-stats?nickname=qa+player',
		)
		expect(getWidget('valorant-stats').buildUrl(context)).toContain(
		'https://twir.localhost/o/channel-key/valorant-stats?',
		)
	})

	it('preserves chat preset query parameters', () => {
		expect(getWidget('chat').buildUrl({ ...context, params: { id: 'chat-preset' } })).toBe(
			'https://twir.localhost/overlays/channel-key/chat?id=chat-preset',
		)
	})
})

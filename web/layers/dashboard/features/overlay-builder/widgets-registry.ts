import type { Component } from 'vue'

import ChatWidgetSettings from './components/widget-settings/ChatWidgetSettings.vue'

export interface OverlayWidgetRegistryEntry {
	readonly key: string
	readonly name: string
	readonly description: string
	readonly icon: string
	readonly buildUrl: (ctx: { readonly origin: string; readonly apiKey: string }) => string
	readonly settingsComponent?: Component
}

export const overlayWidgetRegistry = [
	{
		key: 'chat',
		name: 'Чат',
		description: 'Чат канала Twir',
		icon: 'lucide:messages-square',
		buildUrl: ({ origin, apiKey }) => `${origin}/overlays/${apiKey}/chat`,
		settingsComponent: ChatWidgetSettings,
	},
] satisfies readonly OverlayWidgetRegistryEntry[]

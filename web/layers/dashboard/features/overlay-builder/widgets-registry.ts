import type { Component } from 'vue'

import ChatWidgetSettings from './components/widget-settings/ChatWidgetSettings.vue'

export interface OverlayWidgetRegistryEntry {
	readonly key: string
	readonly nameKey: string
	readonly descriptionKey: string
	readonly icon: string
	readonly buildUrl: (ctx: {
		readonly origin: string
		readonly apiKey: string
		readonly params?: Record<string, string>
	}) => string
	readonly settingsComponent?: Component
}

export const overlayWidgetRegistry = [
	{
		key: 'chat',
		nameKey: 'overlayBuilder.widgets.chat.name',
		descriptionKey: 'overlayBuilder.widgets.chat.description',
		icon: 'lucide:messages-square',
		buildUrl: ({ origin, apiKey, params }) => {
			const base = `${origin}/overlays/${apiKey}/chat`
			return params?.id ? `${base}?id=${params.id}` : base
		},
		settingsComponent: ChatWidgetSettings,
	},
] satisfies readonly OverlayWidgetRegistryEntry[]

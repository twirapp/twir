import type { Layer } from '../types'
import { overlayWidgetRegistry } from '../widgets-registry'

interface WidgetUrlContext {
	readonly origin: string
	readonly apiKey: string
	readonly chatPresetIds: readonly string[]
	readonly chatPresetsReady: boolean
}

export function buildWidgetUrl(widgetKey: string, context: WidgetUrlContext, params?: Record<string, string>): string {
	const widget = overlayWidgetRegistry.find((entry) => entry.key === widgetKey)
	if (!widget) return ''

	return widget.buildUrl({
		origin: context.origin,
		apiKey: context.apiKey,
		params,
	})
}

function getPresetId(url: string): string | null {
	try {
		return new URL(url).searchParams.get('id')
	} catch {
		return null
	}
}

export function resolveWidgetLayerUrl(layer: Layer, context: WidgetUrlContext): string {
	if (!layer.settings.widgetKey || !context.apiKey) return layer.settings.iframeUrl

	const storedUrl = layer.settings.iframeUrl
	if (layer.settings.widgetKey === 'chat') {
		const storedPresetId = getPresetId(storedUrl)
		if (storedPresetId && (!context.chatPresetsReady || context.chatPresetIds.includes(storedPresetId))) return storedUrl
		if (!context.chatPresetsReady) return storedUrl

		const firstPresetId = context.chatPresetIds[0]
		return buildWidgetUrl(
			layer.settings.widgetKey,
			context,
			firstPresetId ? { id: firstPresetId } : undefined,
		)
	}

	return buildWidgetUrl(layer.settings.widgetKey, context)
}

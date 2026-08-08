import { type Ref, watch } from 'vue'

import type { Layer, LayerSettings } from '../types'
import { getWidgetUrlParams } from './widget-url'

interface WidgetEntityResolverOptions<T extends { id: string }> {
	widgetKey: string
	entities: Ref<readonly T[] | undefined>
	refetch: () => Promise<readonly T[]>
	create: () => Promise<unknown>
}

interface WidgetEntityIdOptions<T extends { id: string }> extends WidgetEntityResolverOptions<T> {
	layer: Ref<Layer>
	buildUrl: (id: string) => string
	updateSettings: (updates: Partial<LayerSettings>) => void
}

// Dedupes concurrent auto-creates across widget layer editors, keyed by widget.
const createPromises = new Map<string, Promise<readonly { id: string }[]>>()

export function createWidgetEntityResolver<T extends { id: string }>(
	options: WidgetEntityResolverOptions<T>
) {
	async function ensureEntities(): Promise<readonly T[]> {
		const current = options.entities.value
		if (current === undefined) return []
		if (current.length > 0) return current

		if (!createPromises.has(options.widgetKey)) {
			const promise = options
				.create()
				.then(() => options.refetch())
				.finally(() => createPromises.delete(options.widgetKey))
			createPromises.set(options.widgetKey, promise)
		}

		return (await createPromises.get(options.widgetKey)) as readonly T[]
	}

	async function resolveId(currentId?: string): Promise<string | undefined> {
		const current = options.entities.value
		if (current === undefined) return undefined
		if (currentId && current.some((entity) => entity.id === currentId)) return currentId

		const ensured = await ensureEntities()
		return ensured[0]?.id
	}

	return { resolveId, ensureEntities }
}

export function useWidgetEntityId<T extends { id: string }>(options: WidgetEntityIdOptions<T>) {
	const resolver = createWidgetEntityResolver(options)

	watch(
		[
			() => options.layer.value.settings.widgetKey,
			() => options.layer.value.settings.iframeUrl,
			options.entities,
		],
		async ([widgetKey, iframeUrl]) => {
			if (widgetKey !== options.widgetKey) return

			const currentId = getWidgetUrlParams(iframeUrl).id
			const id = await resolver.resolveId(currentId)
			if (!id || id === currentId) return

			const url = options.buildUrl(id)
			if (url && url !== iframeUrl) options.updateSettings({ iframeUrl: url })
		},
		{ immediate: true }
	)
}

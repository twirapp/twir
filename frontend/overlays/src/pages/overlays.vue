<script lang="ts" setup>
import { type MaybeRef, onMounted, ref, toRaw, toValue, watch } from 'vue'
import { useRoute } from 'vue-router'

import htmlLayer from '@/components/html-layer.vue'
import imageLayer from '@/components/image-layer.vue'
import { type Layer, useOverlays } from '@/composables/overlays/use-overlays.js'

const route = useRoute()

const { layers, parsedLayersData, connectToOverlays, requestLayerData } = useOverlays()

// Храним предыдущее состояние слоёв
const previousLayers = ref<Map<string, Layer>>(new Map())
// Храним интервалы для каждого слоя
const layerIntervals = ref<Map<string, NodeJS.Timeout>>(new Map())

onMounted(() => {
	const apiKey = route.params.apiKey as string
	const overlayId = route.params.overlayId as string
	connectToOverlays(apiKey, overlayId)
})

// Функция проверки, изменились ли только позиционные свойства
const positionKeys = ['posX', 'posY', 'rotation', 'width', 'height']
const nonPositionKeys = [
	'htmlOverlayDataPollSecondsInterval',
	'htmlOverlayHtml',
	'htmlOverlayCss',
	'htmlOverlayJs',
	'imageUrl',
]
function onlyPositionChanged(prev: MaybeRef<Layer>, current: MaybeRef<Layer>): boolean {
	if (!prev) return false

	const prevRaw = toRaw(toValue(prev))
	const currentRaw = toRaw(toValue(current))

	for (const key of nonPositionKeys) {
		// @ts-expect-error
		if (prevRaw.settings[key] !== currentRaw.settings[key]) {
			return false
		}
	}

	for (const key of positionKeys) {
		// @ts-expect-error
		if (prevRaw.settings[key] !== currentRaw.settings[key]) {
			continue
		}
	}

	return true
}

watch(layers, (newLayers) => {
	if (!newLayers.length) return

	for (const layer of newLayers) {
		if (layer.type !== 'HTML') continue

		const prevLayer = previousLayers.value.get(layer.id)

		// Если слоя не было или изменились не только позиционные свойства
		if (!prevLayer || !onlyPositionChanged(prevLayer, layer)) {
			// Очищаем старый интервал, если был
			const existingInterval = layerIntervals.value.get(layer.id)
			if (existingInterval) {
				clearInterval(existingInterval)
			}

			// Запрашиваем данные
			requestLayerData(layer.id)

			// Создаём новый интервал
			const interval = setInterval(
				() => requestLayerData(layer.id),
				layer.settings.htmlOverlayDataPollSecondsInterval * 1000
			)
			layerIntervals.value.set(layer.id, interval)
		}

		// Обновляем сохранённое состояние
		previousLayers.value.set(layer.id, { ...layer })
	}

	// Очищаем интервалы для удалённых слоёв
	const currentLayerIds = new Set(newLayers.map(l => l.id))
	for (const [layerId, interval] of layerIntervals.value.entries()) {
		if (!currentLayerIds.has(layerId)) {
			clearInterval(interval)
			layerIntervals.value.delete(layerId)
			previousLayers.value.delete(layerId)
		}
	}
}, { immediate: true })
</script>

<template>
	<div class="container mx-auto">
		<template v-for="layer of layers" :key="layer.id">
			<htmlLayer v-if="layer.type === 'HTML'" :layer="layer" :parsedData="parsedLayersData[layer.id]" />
			<imageLayer v-else-if="layer.type === 'IMAGE'" :layer="layer" />
		</template>
	</div>
</template>

<style scoped>
.container {
	width: 100%;
	height: 100%;
}
</style>

<script lang="ts" setup>
import { onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'

import htmlLayer from '@/components/html-layer.vue'
import { useOverlays } from '@/composables/overlays/use-overlays.js'

const route = useRoute()

const {
	layers,
	parsedLayersData,
	connectToOverlays,
	requestLayerData,
} = useOverlays()

onMounted(() => {
	const apiKey = route.params.apiKey as string
	const overlayId = route.params.overlayId as string
	connectToOverlays(apiKey, overlayId)
})

watch(layers, (layers) => {
	if (!layers.length) return

	for (const layer of layers) {
		if (layer.type === 'HTML') {
			requestLayerData(layer.id)

			setInterval(
				() => requestLayerData(layer.id),
				layer.settings.htmlOverlayDataPollSecondsInterval * 1000,
			)
		}
	}
})
</script>

<template>
	<div class="container">
		<template v-for="layer of layers" :key="layer.id">
			<htmlLayer :layer="layer" :parsedData="parsedLayersData[layer.id]" />
		</template>
	</div>
</template>

<style scoped>
.container {
	width: 100%;
	height: 100%;
}
</style>

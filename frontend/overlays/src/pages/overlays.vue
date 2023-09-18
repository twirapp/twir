<script lang="ts" setup>
import { watch } from 'vue';
import { useRoute } from 'vue-router';

import htmlLayer from '../components/htmlLayer.vue';
import { useOverlays } from '../sockets/overlays';

const route = useRoute();

const { layers, parsedLayersData, requestLayerData } = useOverlays(
	route.params.apiKey as string,
	route.params.overlayId as string,
);

watch(layers, (l) => {
	if (!l.length) return;

	for (const layer of l) {
		if (layer.type === 'HTML') {
			requestLayerData(layer.id);

			setInterval(
				() => requestLayerData(layer.id),
				layer.settings.htmlOverlayDataPollSecondsInterval * 1000,
			);
		}
	}
});
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

<script setup lang="ts">
import type { OverlayLayer} from '@twir/api/messages/overlays/overlays';
import { OverlayLayerType } from '@twir/api/messages/overlays/overlays';
import { NGrid, NGridItem } from 'naive-ui';
import { useI18n } from 'vue-i18n';

import Card from '@/components/card/card.vue';

defineEmits<{
	select: [OverlayLayer]
}>();

const { t } = useI18n();
</script>

<template>
	<n-grid responsive="screen" cols="1 s:2 m:3 l:4">
		<n-grid-item :span="1">
			<card
				class="cursor-pointer"
				title="HTML"
				@click="() => {
					$emit('select', {
						id: '',
						posX: 0,
						posY: 0,
						width: 200,
						height: 200,
						settings: {
							htmlOverlayCss: '.text { color: red }',
							htmlOverlayHtml: `<span class='text'>$(stream.uptime)</span>`,
							htmlOverlayHtmlDataPollSecondsInterval: 5,
							htmlOverlayJs: `
// will be triggered, when new overlay data comes from backend
function onDataUpdate() {
	console.log('updated')
}
							`
						},
						createdAt: '',
						overlayId: '',
						type: OverlayLayerType.HTML,
						updatedAt: '',
						periodicallyRefetchData: true,
					})
				}"
			>
				<template #content>
					{{ t('overlaysRegistry.html.description') }}
				</template>
			</card>
		</n-grid-item>
	</n-grid>
</template>

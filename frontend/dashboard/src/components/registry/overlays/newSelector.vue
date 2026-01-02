<script setup lang="ts">
import { ChannelOverlayLayerType } from '@/gql/graphql'
import { NGrid, NGridItem } from 'naive-ui'
import { useI18n } from 'vue-i18n'

import Card from '@/components/card/card.vue'

interface OverlayLayerForm {
	type: ChannelOverlayLayerType
	posX: number
	posY: number
	width: number
	height: number
	periodicallyRefetchData: boolean
	settings: {
		htmlOverlayHtml: string
		htmlOverlayCss: string
		htmlOverlayJs: string
		htmlOverlayDataPollSecondsInterval: number
	}
}

defineEmits<{
	select: [OverlayLayerForm]
}>()

const { t } = useI18n()
</script>

<template>
	<n-grid responsive="screen" cols="1 s:2 m:3 l:4">
		<n-grid-item :span="1">
			<card
				class="cursor-pointer"
				title="HTML"
				@click="
					() => {
						$emit('select', {
							posX: 0,
							posY: 0,
							width: 200,
							height: 200,
							settings: {
								htmlOverlayCss: '.text { color: red }',
								htmlOverlayHtml: `<span class='text'>$(stream.uptime)</span>`,
								htmlOverlayDataPollSecondsInterval: 5,
								htmlOverlayJs: `
// will be triggered, when new overlay data comes from backend
function onDataUpdate() {
	console.log('updated')
}
							`,
							},
							type: ChannelOverlayLayerType.Html,
							periodicallyRefetchData: true,
						})
					}
				"
			>
				<template #content>
					{{ t('overlaysRegistry.html.description') }}
				</template>
			</card>
		</n-grid-item>
	</n-grid>
</template>

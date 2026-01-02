<script lang="ts" setup>
import { Settings, Trash2 } from 'lucide-vue-next'
import type { ChannelOverlayLayerType } from '@/gql/graphql'
import { NButton, NCard, useThemeVars } from 'naive-ui'
import { computed } from 'vue'

import { convertOverlayLayerTypeToText } from '../helpers.js'

const theme = useThemeVars()
const activeLayourCardColor = computed(() => theme.value.infoColor)

export interface LayerProps {
	isFocused: boolean
	layerIndex: number
	type: ChannelOverlayLayerType
}

defineProps<LayerProps>()

defineEmits<{
	focus: [index: number]
	remove: [index: number]
	openSettings: []
}>()
</script>

<template>
	<n-card
		:title="convertOverlayLayerTypeToText(type)"
		class="cursor-pointer"
		:style="{
			border: isFocused ? `1px solid ${activeLayourCardColor}` : undefined,
		}"
		@click="$emit('focus', layerIndex)"
	>
		<slot />

		<div class="flex gap-3 w-full">
			<n-button secondary class="flex-1" @click="$emit('openSettings')">
				<Settings class="size-4" />
			</n-button>
			<n-button secondary class="flex-1" type="error" @click="$emit('remove', layerIndex)">
				<Trash2 class="size-4" />
			</n-button>
		</div>
	</n-card>
</template>

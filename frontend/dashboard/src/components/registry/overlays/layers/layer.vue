<script lang="ts" setup>
import { IconSettings, IconTrash } from '@tabler/icons-vue';
import { type OverlayLayerType } from '@twir/grpc/generated/api/api/overlays';
import { NCard, NButton, useThemeVars } from 'naive-ui';
import { computed } from 'vue';

import { convertOverlayLayerTypeToText } from '../helpers.js';

const theme = useThemeVars();
const activeLayourCardColor = computed(() => theme.value.infoColor);

export type LayerProps = {
	isFocused: boolean
	layerIndex: number
	type: OverlayLayerType
}

defineProps<LayerProps>();

defineEmits<{
	focus: [index: number]
	remove: [index: number]
	openSettings: []
}>();
</script>

<template>
	<n-card
		:title="convertOverlayLayerTypeToText(type)"
		style="cursor: pointer;"
		:style="{
			border: isFocused ? `1px solid ${activeLayourCardColor}` : undefined
		}"
		@click="$emit('focus', layerIndex)"
	>
		<slot />

		<div style="display: flex; gap: 12px; width: 100%">
			<n-button
				style="flex: 1" secondary
				@click="$emit('openSettings')"
			>
				<IconSettings />
			</n-button>
			<n-button
				style="flex: 1" secondary type="error"
				@click="$emit('remove', layerIndex)"
			>
				<IconTrash />
			</n-button>
		</div>
	</n-card>
</template>

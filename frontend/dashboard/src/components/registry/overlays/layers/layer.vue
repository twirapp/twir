<script lang="ts" setup>
import { IconSettings, IconTrash } from '@tabler/icons-vue';
import { type OverlayLayerType } from '@twir/api/messages/overlays/overlays';
import { NButton, NCard, useThemeVars } from 'naive-ui';
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
		class="cursor-pointer"
		:style="{
			border: isFocused ? `1px solid ${activeLayourCardColor}` : undefined
		}"
		@click="$emit('focus', layerIndex)"
	>
		<slot />

		<div class="flex gap-3 w-full">
			<n-button
				secondary
				class="flex-1"
				@click="$emit('openSettings')"
			>
				<IconSettings />
			</n-button>
			<n-button
				secondary
				class="flex-1"
				type="error"
				@click="$emit('remove', layerIndex)"
			>
				<IconTrash />
			</n-button>
		</div>
	</n-card>
</template>

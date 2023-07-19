<script setup lang="ts">
import { IconSettings, IconCopy } from '@tabler/icons-vue';
import { NCard, NButton, useThemeVars, useMessage, NTooltip } from 'naive-ui';
import { FunctionalComponent, computed } from 'vue';

const themeVars = useThemeVars();
const titleColor = computed(() => themeVars.value.textColor1);

const props = defineProps<{
	description: string;
	title: string;
	overlayLink?: string;
	icon: FunctionalComponent;
}>();

defineEmits<{
	openSettings: [];
}>();

const messages = useMessage();
const copyOverlayLink = () => {
	if (!props.overlayLink) return;

	navigator.clipboard.writeText(props.overlayLink);
	messages.success('Copied link url, paste it in obs as browser source');
};
</script>

<template>
	<n-card>
		<div style="display: flex; flex-direction: column">
			<component
				:is="icon"
				style="width: 48px; height: 48px; stroke-width: 2px; stroke: #61e8bb; margin-bottom: 16px"
			/>
			<div>
				<h2 class="card-title">
					{{ title }}
				</h2>
				<span :style="{ color: themeVars.textColor3 }">
					{{ description }}
				</span>
			</div>
			<div class="card-buttons">
				<n-button secondary size="large" class="card-button" @click="$emit('openSettings')">
					<span>Settings</span>
					<IconSettings />
				</n-button>
				<n-tooltip :disabled="!!overlayLink">
					<template #trigger>
						<n-button
							size="large"
							:disabled="!overlayLink"
							class="card-button"
							@click="copyOverlayLink"
						>
							<span>Copy overlay link</span>
							<IconCopy />
						</n-button>
					</template>
					You should configure overlay first
				</n-tooltip>
			</div>
		</div>
	</n-card>
</template>

<style scoped>
.card-button span {
	font-size: 14px;
	line-height: 20px
}

.card-button svg {
	height: 20px;
	width: 20px;
	margin-left: 8px
}

@media (max-width: 568px) {
	.card-button {
		width: 100%;
	}
}

.card-title {
	color: v-bind(titleColor);
	margin: 0 0 12px 0;
	font-size: 20px;
	line-height: 24px;
}

.card-buttons {
	display: flex;
	gap: 8px;
	margin-top: 20px;
	flex-wrap: wrap;
}
</style>

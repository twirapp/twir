<script setup lang="ts">
import { IconSettings, IconCopy } from '@tabler/icons-vue';
import { NCard, NButton, useThemeVars, useMessage, NTooltip } from 'naive-ui';
import { FunctionalComponent } from 'vue';

const themeVars = useThemeVars();

const props =
	defineProps<{
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
				<h2
					:style="{
						color: themeVars.textColor1,
						margin: '0 0 12px 0',
						fontSize: '20px',
						lineHeight: '24px',
					}"
				>
					{{ title }}
				</h2>
				<span :style="{ color: themeVars.textColor3 }">
					{{ description }}
				</span>
			</div>
			<div style="display: flex; gap: 8px; margin-top: 20px; flex-wrap: wrap">
				<n-button secondary size="large" @click="$emit('openSettings')" class="card-button">
					<span style="font-size: 14px; line-height: 20px">Settings</span>
					<IconSettings style="height: 20px; width: 20px; margin-left: 8px" />
				</n-button>
				<n-tooltip :disabled="!!overlayLink">
					<template #trigger>
						<n-button
							size="large"
							:disabled="!overlayLink"
							@click="copyOverlayLink"
							class="card-button"
						>
							<span style="font-size: 14px; line-height: 20px">Copy overlay link</span>
							<IconCopy style="height: 20px; margin-left: 8px" />
						</n-button>
					</template>
					You should configure overlay first
				</n-tooltip>
			</div>
		</div>
	</n-card>
</template>

<style scoped>
@media (max-width: 568px) {
	.card-button {
		width: 100%;
	}
}
</style>

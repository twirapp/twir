<script setup lang="ts">
import { IconSettings, IconCopy } from '@tabler/icons-vue';
import { NButton, useMessage, NTooltip } from 'naive-ui';
import { FunctionalComponent } from 'vue';

import Card from '@/components/card/card.vue';

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
	<card :title="title" :icon="icon">
		<template #content>
			{{ description }}
		</template>

		<template #footer>
			<n-button secondary size="large" @click="$emit('openSettings')">
				<span>Settings</span>
				<IconSettings />
			</n-button>
			<n-tooltip :disabled="!!overlayLink">
				<template #trigger>
					<n-button
						size="large"
						:disabled="!overlayLink"
						@click="copyOverlayLink"
					>
						<span>Copy overlay link</span>
						<IconCopy />
					</n-button>
				</template>
				You should configure overlay first
			</n-tooltip>
		</template>
	</card>
</template>

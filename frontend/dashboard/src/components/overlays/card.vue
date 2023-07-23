<script setup lang="ts">
import { IconSettings, IconCopy } from '@tabler/icons-vue';
import { NButton, useMessage, NTooltip } from 'naive-ui';
import { FunctionalComponent } from 'vue';
import { useI18n } from 'vue-i18n';

import { useProfile, useUserAccessFlagChecker } from '@/api/index.js';
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

const { t } = useI18n();
const messages = useMessage();
const copyOverlayLink = () => {
	if (!props.overlayLink) return;

	navigator.clipboard.writeText(props.overlayLink);
	messages.success(t('overlays.copied'));
};

const userCanEditOverlays = useUserAccessFlagChecker('MANAGE_OVERLAYS');

const { data: profile } = useProfile();

</script>

<template>
	<card :title="title" :icon="icon" style="height: 100%;">
		<template #content>
			{{ description }}
		</template>

		<template #footer>
			<n-button :disabled="!userCanEditOverlays" secondary size="large" @click="$emit('openSettings')">
				<span>{{ t('sharedButtons.settings') }}</span>
				<IconSettings />
			</n-button>
			<n-tooltip :disabled="!!overlayLink || profile?.id != profile?.selectedDashboardId">
				<template #trigger>
					<n-button
						size="large"
						:disabled="!overlayLink || profile?.id != profile?.selectedDashboardId"
						@click="copyOverlayLink"
					>
						<span>{{ t('overlays.copyOverlayLink') }}</span>
						<IconCopy />
					</n-button>
				</template>
				<span v-if="profile?.id != profile?.selectedDashboardId">{{ t('overlays.noAccess') }}</span>
				<span v-else>{{ t('overlays.uncongirured') }}</span>
			</n-tooltip>
		</template>
	</card>
</template>

<script setup lang="ts">
import { IconSettings, IconCopy } from '@tabler/icons-vue';
import { NButton, useMessage, NTooltip } from 'naive-ui';
import { FunctionalComponent, computed } from 'vue';
import { useI18n } from 'vue-i18n';

import { useProfile, useUserAccessFlagChecker } from '@/api/index.js';
import Card from '@/components/card/card.vue';

const props = withDefaults(defineProps<{
	description: string;
	title: string;
	overlayPath: string;
	icon: FunctionalComponent;
	showSettings: boolean
	copyDisabled?: boolean,
}>(), { showSettings: true, copyDisabled: false });

defineEmits<{
	openSettings: [];
}>();

const { t } = useI18n();
const messages = useMessage();
const { data: profile } = useProfile();

const overlayLink = computed(() => {
	return `${window.location.origin}/overlays/${profile.value?.apiKey}/${props.overlayPath}`;
});

const copyOverlayLink = () => {
	if (!props.overlayPath) return;

	navigator.clipboard.writeText(overlayLink.value);
	messages.success(t('overlays.copied'));
};

const userCanEditOverlays = useUserAccessFlagChecker('MANAGE_OVERLAYS');


</script>

<template>
	<card :title="title" :icon="icon" style="height: 100%;">
		<template #content>
			{{ description }}
		</template>

		<template #footer>
			<n-button v-if="showSettings" :disabled="!userCanEditOverlays" secondary size="large" @click="$emit('openSettings')">
				<span>{{ t('sharedButtons.settings') }}</span>
				<IconSettings />
			</n-button>
			<n-tooltip :disabled="!!overlayLink || profile?.id != profile?.selectedDashboardId">
				<template #trigger>
					<n-button
						size="large"
						:disabled="copyDisabled || !overlayLink || profile?.id != profile?.selectedDashboardId"
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

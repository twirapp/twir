<script setup lang="ts">
import { IconSettings, IconCopy } from '@tabler/icons-vue';
import { NButton, NTooltip } from 'naive-ui';
import { FunctionalComponent } from 'vue';
import { useI18n } from 'vue-i18n';

import { useCopyOverlayLink } from './copyOverlayLink.js';

import { useProfile, useUserAccessFlagChecker } from '@/api/index.js';
import Card from '@/components/card/card.vue';

const props = withDefaults(defineProps<{
	description: string;
	title: string;
	overlayPath: string;
	icon: FunctionalComponent;
	showSettings?: boolean
	copyDisabled?: boolean,
}>(), { showSettings: true, copyDisabled: false });

defineEmits<{
	openSettings: [];
}>();

const { t } = useI18n();
const { data: profile } = useProfile();

const { copyOverlayLink } = useCopyOverlayLink(props.overlayPath);

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
			<n-tooltip :disabled="profile?.id !== profile?.selectedDashboardId">
				<template #trigger>
					<n-button
						size="large"
						:disabled="copyDisabled || profile?.id != profile?.selectedDashboardId"
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

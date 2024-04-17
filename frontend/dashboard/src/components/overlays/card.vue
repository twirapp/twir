<script setup lang="ts">
import { IconSettings, IconCopy } from '@tabler/icons-vue';
import { NButton, NTooltip } from 'naive-ui';
import { FunctionalComponent } from 'vue';
import { useI18n } from 'vue-i18n';

import { useCopyOverlayLink } from './copyOverlayLink.js';

import { useProfile, useUserAccessFlagChecker } from '@/api/index.js';
import Card from '@/components/card/card.vue';
import { storeToRefs } from 'pinia';

const props = withDefaults(defineProps<{
	description: string;
	title: string;
	overlayPath?: string;
	icon: FunctionalComponent;
	iconStroke?: number;
	showSettings?: boolean;
	copyDisabled?: boolean;
	showCopy?: boolean;
}>(), {
	showSettings: true,
	copyDisabled: false,
	showCopy: true,
	iconStroke: 1,
	overlayPath: '',
});

defineEmits<{
	openSettings: [];
}>();

const { t } = useI18n();
const { data: profile } = storeToRefs(useProfile());

const { copyOverlayLink } = useCopyOverlayLink(props.overlayPath);

const userCanEditOverlays = useUserAccessFlagChecker('MANAGE_OVERLAYS');
</script>

<template>
	<card :title="title" :icon="icon" :icon-stroke="iconStroke" class="h-full">
		<template #content>
			{{ description }}
		</template>

		<template #footer>
			<n-button
				v-if="showSettings" :disabled="!userCanEditOverlays" secondary size="large"
				@click="$emit('openSettings')"
			>
				<div class="flex gap-1">
					<span>{{ t('sharedButtons.settings') }}</span>
					<IconSettings />
				</div>
			</n-button>
			<n-tooltip v-if="showCopy" :disabled="profile?.id !== profile?.selectedDashboardId">
				<template #trigger>
					<n-button
						size="large"
						:disabled="copyDisabled || profile?.id != profile?.selectedDashboardId"
						@click="copyOverlayLink()"
					>
						<div class="flex gap-1">
							<span>{{ t('overlays.copyOverlayLink') }}</span>
							<IconCopy />
						</div>
					</n-button>
				</template>
				<span v-if="profile?.id != profile?.selectedDashboardId">{{ t('overlays.noAccess') }}</span>
				<span v-else>{{ t('overlays.uncongirured') }}</span>
			</n-tooltip>
		</template>
	</card>
</template>

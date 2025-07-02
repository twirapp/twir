<script setup lang="ts">
import { IconSettings } from '@tabler/icons-vue';
import { useThrottleFn } from '@vueuse/core';
import { NButton, NPopconfirm, NRow, NSpace, NSwitch, NTag, NText, NTooltip } from 'naive-ui';
import { useI18n } from 'vue-i18n';

import { flatEvents, getEventName, getOperation } from './helpers.js';
import { EditableEvent } from './types.js';

import { useEventsManager, useUserAccessFlagChecker } from '@/api/index.js';
import Card from '@/components/card/card.vue';
import { ChannelRolePermissionEnum } from '@/gql/graphql';

const props = defineProps<{
	event: EditableEvent
}>();

defineEmits<{
	openSettings: [id: string]
}>();

const eventsManager = useEventsManager();
const eventsPatcher = eventsManager.patch!;
const eventsDeleter = eventsManager.deleteOne;

const throttledSwitchState = useThrottleFn((v: boolean) => {
	eventsPatcher.mutate({ id: props.event.id!, enabled: v });
}, 500);

const { t } = useI18n();

const userCanEditEvents = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageEvents);
</script>

<template>
	<card class="h-full" :icon="flatEvents[event.type]?.icon" :title="getEventName(event.type)">
		<template #headerExtra>
			<n-switch
				:value="event.enabled"
				:disabled="!userCanEditEvents"
				@update-value="(v) => throttledSwitchState(v)"
			/>
		</template>

		<template #content>
			<n-space vertical>
				<n-text>{{ event.description }}</n-text>
				<n-row class="gap-2">
					<n-tooltip v-for="(operation, index) of event.operations" :key="index">
						<template #trigger>
							<n-tag
								:disabled="!operation.enabled"
								:bordered="false"
								:type="getOperation(operation.type)?.color ?? 'info'"
							>
								{{ getOperation(operation.type)?.name ?? 'Unknown event' }}
							</n-tag>
						</template>
						<n-space vertical>
							<n-text>{{ operation.input }}</n-text>
							<n-text>{{ t('events.delay') }}: {{ operation.delay }} | {{ t('events.repeat') }}: {{ operation.repeat }}</n-text>
						</n-space>
					</n-tooltip>
				</n-row>
			</n-space>
		</template>

		<template #footer>
			<n-button secondary size="large" :disabled="!userCanEditEvents || !event.id" @click="$emit('openSettings', event.id!)">
				<span>{{ t('sharedButtons.settings') }}</span>
				<IconSettings />
			</n-button>
			<n-popconfirm
				:positive-text="t('deleteConfirmation.confirm')"
				:negative-text="t('deleteConfirmation.cancel')"
				@positive-click="eventsDeleter.mutateAsync({ id: event.id! })"
			>
				<template #trigger>
					<n-button secondary type="error" size="large" :disabled="!userCanEditEvents">
						<span>{{ t('sharedButtons.delete') }}</span>
					</n-button>
				</template>
				{{ t('deleteConfirmation.text') }}
			</n-popconfirm>
		</template>
	</card>
</template>

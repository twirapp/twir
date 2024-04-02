<!-- eslint-disable no-undef -->
<script setup lang="ts" generic="Message extends { count?: number, text: string }">
import { IconTrash } from '@tabler/icons-vue';
import { NInput, NInputNumber, NInputGroup, NInputGroupLabel, NButton, NSwitch, NAlert } from 'naive-ui';
import { onMounted, nextTick } from 'vue';
import { useI18n } from 'vue-i18n';

import { useUserAccessFlagChecker } from '@/api/index.js';

const props = defineProps<{
	alertMessage?: string
	withCount?: boolean
	countLabel?: string
	maxMessages: number
	defaultMessageText: string
	minCount?: number
	minCooldown: number
}>();

const hasAccessToManageAlerts = useUserAccessFlagChecker('MANAGE_ALERTS');
const enabled = defineModel<boolean>('enabled', { default: true });
const messages = defineModel<Array<Message>>('messages');
const cooldown = defineModel<number>('cooldown');

onMounted(async () => {
	await nextTick();
	if (!messages.value?.length) {
		createMessage();
	}
});

function createMessage() {
	if (!hasAccessToManageAlerts) return;
	if (props.withCount) {
		const latest = messages.value?.at(-1);
		messages.value?.push({
			count: latest?.count ? latest.count + 1 : 1,
			text: props.defaultMessageText,
		} as any);
	} else {
		messages.value?.push({ text: props.defaultMessageText } as any);
	}
}

function removeMessage(index: number) {
	if (!hasAccessToManageAlerts) return;
	messages.value = messages.value?.filter((_, i) => i !== index);
}

const { t } = useI18n();

defineSlots<{
	header?: any
}>();
</script>

<template>
	<section class="flex gap-2 flex-col">
		<div class="flex gap-1">
			<span>{{ t('sharedTexts.enabled') }}</span>
			<n-switch v-model:value="enabled" />
		</div>

		<div class="flex gap-1 flex-col">
			<span>{{ t('chatAlerts.cooldown') }}</span>
			<n-input-number
				v-model:value="cooldown"
				:min="minCooldown"
				:max="9999"
				class="w-[10%] min-w-[100px]"
			/>
		</div>

		<slot name="header" />
	</section>

	<n-alert v-if="alertMessage" type="info" title="Info" class="mt-3">
		<span v-html="alertMessage" />
	</n-alert>

	<ul class="flex flex-col gap-3.5 p-0 mx-0 my-3.5">
		<li
			v-for="(m, index) of messages"
			:key="index"
			class="flex justify-between gap-3.5"
		>
			<div class="message-item-content">
				<n-input-group v-if="withCount && countLabel" class="w-auto min-w-[200px]">
					<n-input-group-label>{{ countLabel }} >=</n-input-group-label>
					<n-input-number
						v-model:value="m.count"
						:min="minCount ?? 1"
						:max="9999999"
						class="flex-1"
					/>
				</n-input-group>

				<n-input v-model:value="m.text" :title="m.text" />
			</div>

			<n-button :disabled="!hasAccessToManageAlerts" secondary type="error" @click="removeMessage(index)">
				<IconTrash />
			</n-button>
		</li>
	</ul>

	<n-button
		block
		secondary
		type="success"
		:disabled="(messages?.length === maxMessages) || !hasAccessToManageAlerts"
		@click="createMessage"
	>
		<span v-if="messages?.length">{{ t('sharedButtons.create') }} ({{ messages.length }} / {{ maxMessages }})</span>
		<span v-else>{{ t('sharedButtons.create') }}</span>
	</n-button>
</template>

<style scoped>
.message-item-content {
  @apply flex w-full gap-x-3.5;
}

@media screen and (max-width: 768px) {
  .message-item-content {
    @apply flex-col gap-x-0 gap-y-1.5;
  }
}
</style>

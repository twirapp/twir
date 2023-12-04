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
	<section style="display: flex; gap: 8px; flex-direction: column;">
		<div style="display: flex; gap: 4px;">
			<span>{{ t('sharedTexts.enabled') }}</span>
			<n-switch v-model:value="enabled" />
		</div>

		<div style="display: flex; gap: 4px; flex-direction: column;">
			<span>{{ t('chatAlerts.cooldown') }}</span>
			<n-input-number
				v-model:value="cooldown"
				:min="minCooldown"
				:max="9999"
				style="width: 10%; min-width: 100px;"
			/>
		</div>

		<slot name="header" />
	</section>

	<n-alert v-if="alertMessage" type="info" title="Info" style="margin-top: 14px;">
		<span v-html="alertMessage" />
	</n-alert>

	<ul class="messages">
		<li
			v-for="(m, index) of messages"
			:key="index"
			class="messageItem"
		>
			<div class="messageItemContent">
				<n-input-group v-if="withCount && countLabel" style="width: auto; min-width: 200px">
					<n-input-group-label>{{ countLabel }} >=</n-input-group-label>
					<n-input-number
						v-model:value="m.count"
						:min="minCount ?? 1"
						:max="9999999"
						style="flex: 1"
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
.messages {
	display: flex;
	gap: 14px;
	flex-direction: column;
	padding: 0;
	margin: 14px 0;
}

.messageItem {
	display: flex;
	gap: 14px;
	justify-content: space-between;
}

.messageItemContent {
	width: 100%;
	display: flex;
	column-gap: 14px;
}

@media screen and (max-width: 768px){
	.messageItemContent {
		flex-direction: column;
		column-gap: 0;
		row-gap: 6px;
	}
}
</style>

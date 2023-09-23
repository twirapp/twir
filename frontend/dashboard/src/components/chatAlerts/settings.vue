<!-- eslint-disable no-undef -->
<script setup lang="ts" generic="Message extends { count?: number, text: string }">
import { IconTrash } from '@tabler/icons-vue';
import { NInput, NInputNumber, NInputGroup, NInputGroupLabel, NButton, NSwitch, NAlert } from 'naive-ui';
import { onMounted, nextTick } from 'vue';
import { useI18n } from 'vue-i18n';

const props = defineProps<{
	alertMessage?: string
	withCount?: boolean
	countLabel?: string
	maxMessages: number
	defaultMessageText: string
}>();

const enabled = defineModel<boolean>('enabled', { default: true });
const messages = defineModel<Array<Message>>('messages');

onMounted(async () => {
	await nextTick();
	if (!messages.value?.length) {
		createMessage();
	}
});

function createMessage() {
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
	messages.value = messages.value?.filter((_, i) => i !== index);
}

const { t } = useI18n();
</script>

<template>
	<div style="display: flex; gap: 4px; margin-top: 10px;">
		<span>{{ t('sharedTexts.enabled') }}</span>
		<n-switch v-model:value="enabled" />
	</div>

	<n-alert v-if="alertMessage" type="info" title="Info" style="margin-top: 14px;">
		{{ alertMessage }}
	</n-alert>

	<div class="messages">
		<div
			v-for="(m, index) of messages"
			:key="index"
			style="display: flex; gap: 14px"
		>
			<n-input-group v-if="withCount && countLabel" style="width: auto;">
				<n-input-group-label>{{ countLabel }}</n-input-group-label>
				<n-input-number v-if="withCount && countLabel" v-model:value="m.count" :min="1" :max="9999999" />
			</n-input-group>

			<n-input v-model:value="m.text" />
			<n-button secondary type="error" @click="removeMessage(index)">
				<IconTrash />
			</n-button>
		</div>
	</div>

	<n-button block secondary type="success" :disabled="messages?.length === maxMessages" @click="createMessage">
		{{ t('sharedButtons.create') }} ({{ messages?.length }} / {{ maxMessages }})
	</n-button>
</template>

<style scoped>
.messages {
	display: flex;
	gap: 14px;
	flex-direction: column;
	gap: 14px;
	margin-bottom: 14px;
	margin-top: 14px;
}
</style>
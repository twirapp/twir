<script setup lang="ts">
import { IconMessageCircleQuestion, IconTrash } from '@tabler/icons-vue';
import { NModal, NInput, NButton, NSwitch, NDivider, useMessage } from 'naive-ui';
import { ref, watch, toRaw } from 'vue';
import { useI18n } from 'vue-i18n';

import Card from './card.vue';

import { use8ballSettings, use8ballUpdateSettings } from '@/api/index.js';
import CommandButton from '@/features/commands/components/command-button.vue';

const isModalOpened = ref(false);

const maxAnswers = 25;

const { data: settings } = use8ballSettings();
const updater = use8ballUpdateSettings();

const formValue = ref({
	enabled: false,
	answers: ['Yes', 'No'],
});

watch(settings, (v) => {
	if (!v) return;

	const raw = toRaw(v);
	formValue.value.answers = raw.answers;
	formValue.value.enabled = raw.enabled;
}, { immediate: true });

const { t } = useI18n();

const notifications = useMessage();

async function save() {
	await updater.mutateAsync({
		answers: formValue.value.answers,
		enabled: formValue.value.enabled,
	});
	notifications.success(t('sharedTexts.saved'));
}
</script>

<template>
	<card
		title="8ball"
		:icon="IconMessageCircleQuestion"
		:icon-stroke="1"
		:description="t('games.8ball.description')"
		@open-settings="isModalOpened = true"
	/>

	<n-modal
		v-model:show="isModalOpened"
		:mask-closable="false"
		:segmented="true"
		preset="card"
		title="8ball"
		content-style="padding: 10px; width: 100%"
		style="width: 500px; max-width: calc(100vw - 40px)"
	>
		<div class="flex gap-6">
			<div class="flex flex-col gap-1 items-center">
				<span>{{ t('sharedTexts.enabled') }}</span>
				<n-switch v-model:value="formValue.enabled"></n-switch>
			</div>

			<command-button name="8ball" />
		</div>

		<n-divider />

		<h3>{{ t('games.8ball.answers') }} ({{ formValue.answers.length }}/{{ maxAnswers }})</h3>

		<div class="flex flex-col gap-2">
			<div
				v-for="(_, index) of formValue.answers"
				:key="index"
				class="flex gap-1"
			>
				<n-input
					v-model:value="formValue.answers[index]"
					placeholder="Yes"
				/>

				<n-button
					secondary
					type="error"
					@click="() => {
						formValue.answers = formValue.answers.filter((_, i) => i != index)
					}"
				>
					<IconTrash />
				</n-button>
			</div>

			<n-button
				secondary
				type="info"
				block
				:disabled="formValue.answers.length >= maxAnswers"
				@click="() => formValue.answers.push('')"
			>
				{{ t('sharedButtons.create') }}
			</n-button>
		</div>

		<n-divider />

		<n-button block secondary type="success" @click="save">
			{{ t('sharedButtons.save') }}
		</n-button>
	</n-modal>
</template>

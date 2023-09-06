<script setup lang="ts">
import { IconBomb } from '@tabler/icons-vue';
import { NModal, NInput, NInputNumber, NFormItem, NButton, NSwitch, NDivider, useMessage } from 'naive-ui';
import { ref, watch, toRaw } from 'vue';
import { useI18n } from 'vue-i18n';

import Card from './card.vue';
import Command from './command.vue';

import { useRussianRouletteSettings, useRussianRouletteUpdateSettings } from '@/api/index.js';

const isModalOpened = ref(false);

const { data: settings } = useRussianRouletteSettings();
const updater = useRussianRouletteUpdateSettings();

const formValue = ref({
	enabled: false,
	canBeUsedByModerator: false,
	timeoutSeconds: 60,
	decisionSeconds: 2,
	chargedBullets: 1,
	initMessage: '{sender} has initiated a game of roulette. Is luck on their side?',
	surviveMessage: '{sender} survives the game of roulette! Luck smiles upon them.',
	deathMessage: `{sender} couldn't make it through the game of roulette. Unfortunately, luck wasn't on their side this time.`,
});

watch(settings, (v) => {
	if (!v) return;

	const raw = toRaw(v);

	formValue.value.enabled = raw.enabled;
	formValue.value.canBeUsedByModerator = raw.canBeUsedByModerator;
	formValue.value.timeoutSeconds = raw.timeoutSeconds;
	formValue.value.decisionSeconds = raw.decisionSeconds;
	formValue.value.initMessage = raw.initMessage;
	formValue.value.surviveMessage = raw.surviveMessage;
	formValue.value.deathMessage = raw.deathMessage;
	formValue.value.chargedBullets = raw.chargedBullets;
});

const { t } = useI18n();

const notifications = useMessage();

async function save() {
	const values = formValue.value;
	await updater.mutateAsync(values);
	notifications.success(t('sharedTexts.saved'));
}
</script>

<template>
	<card
		title="Russian Roulette"
		:icon="IconBomb"
		:description="t('games.russianRoulette.description')"
		@open-settings="isModalOpened = true"
	/>

	<n-modal
		v-model:show="isModalOpened"
		:mask-closable="false"
		:segmented="true"
		preset="card"
		title="Russian Roulette"
		content-style="padding: 10px; width: 100%"
		style="width: 500px; max-width: calc(100vw - 40px)"
	>
		<div style="display: flex; flex-direction: column; gap: 4px; align-items: start;">
			<span>{{ t('sharedTexts.enabled') }}</span>
			<n-switch v-model:value="formValue.enabled" />
		</div>

		<Command name="roulette" />

		<div style="display: flex; flex-direction: column; gap: 8px; margin-top: 10px">
			<n-form-item :label="t('games.russianRoulette.canBeUsedByModerator')">
				<n-switch v-model:value="formValue.canBeUsedByModerator" />
			</n-form-item>

			<n-form-item :label="t('games.russianRoulette.chargedBullets')">
				<n-input-number v-model:value="formValue.chargedBullets" :min="1" :max="6" />
			</n-form-item>

			<n-form-item :label="t('games.russianRoulette.timeoutSeconds')">
				<n-input-number v-model:value="formValue.timeoutSeconds" :max="1209600" />
			</n-form-item>

			<n-form-item :label="t('games.russianRoulette.decisionSeconds')">
				<n-input-number v-model:value="formValue.decisionSeconds" :max="60" />
			</n-form-item>

			<n-divider />

			<n-form-item :label="t('games.russianRoulette.initMessage')">
				<n-input v-model:value="formValue.initMessage" :maxlength="450" type="textarea" autosize :rows="1" />
			</n-form-item>

			<n-form-item :label="t('games.russianRoulette.surviveMessage')">
				<n-input v-model:value="formValue.surviveMessage" :maxlength="450" type="textarea" autosize :rows="1" />
			</n-form-item>

			<n-form-item :label="t('games.russianRoulette.deathMessage')">
				<n-input v-model:value="formValue.deathMessage" :maxlength="450" type="textarea" autosize :rows="1" />
			</n-form-item>
		</div>

		<n-divider />

		<n-button block secondary type="success" @click="save">
			{{ t('sharedButtons.save') }}
		</n-button>
	</n-modal>
</template>

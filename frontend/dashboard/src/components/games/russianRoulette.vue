<script setup lang="ts">
import { IconBomb } from '@tabler/icons-vue';
import type { UpdateRussianRouletteSettings } from '@twir/api/messages/games/games';
import {
	NModal,
	NInput,
	NInputNumber,
	NFormItem,
	NButton,
	NSwitch,
	NDivider,
	NSpace,
} from 'naive-ui';
import { ref, watch, toRaw } from 'vue';
import { useI18n } from 'vue-i18n';

import Card from './card.vue';
import Command from '../commandButton.vue';

import { useRussianRouletteSettings, useRussianRouletteUpdateSettings } from '@/api/index.js';
import { useNaiveDiscrete } from '@/composables/use-naive-discrete';

const isModalOpened = ref(false);

const { data: settings } = useRussianRouletteSettings();
const updater = useRussianRouletteUpdateSettings();

const initialSettings: UpdateRussianRouletteSettings = {
	enabled: false,
	canBeUsedByModerator: false,
	timeoutSeconds: 60,
	decisionSeconds: 2,
	chargedBullets: 1,
	initMessage: '{sender} has initiated a game of roulette. Is luck on their side?',
	surviveMessage: '{sender} survives the game of roulette! Luck smiles upon them.',
	deathMessage: `{sender} couldn't make it through the game of roulette. Unfortunately, luck wasn't on their side this time.`,
	tumberSize: 6,
};

const formValue = ref<UpdateRussianRouletteSettings>({ ...initialSettings });

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
	formValue.value.tumberSize = raw.tumberSize;
}, { immediate: true });

const { t } = useI18n();

const { dialog, notification } = useNaiveDiscrete();

async function save() {
	const values = formValue.value;
	await updater.mutateAsync(values);
	notification.success({
		title: t('sharedTexts.saved'),
		duration: 2500,
	});
}

function resetSettings() {
	dialog.create({
		type: 'warning',
		title: t('sharedTexts.dangerZone'),
		content: t('sharedTexts.setDefaultSettings'),
		positiveText: t('sharedButtons.confirm'),
		negativeText: t('sharedButtons.close'),
		onPositiveClick: () => {
			formValue.value = initialSettings;
			save();
		},
	});
}
</script>

<template>
	<card
		title="Russian Roulette"
		:icon="IconBomb"
		:icon-stroke="1"
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
		<div class="flex gap-6">
			<div class="flex flex-col gap-1 items-start">
				<span>{{ t('sharedTexts.enabled') }}</span>
				<n-switch v-model:value="formValue.enabled" />
			</div>

			<Command name="roulette" />
		</div>

		<n-divider />

		<div class="flex flex-col gap-2 mt-[10px]">
			<n-form-item :label="t('games.russianRoulette.canBeUsedByModerator')">
				<n-switch v-model:value="formValue.canBeUsedByModerator" />
			</n-form-item>

			<n-form-item :label="t('games.russianRoulette.tumberSize')">
				<n-input-number v-model:value="formValue.tumberSize" :min="2" :max="100" />
			</n-form-item>

			<n-form-item
				:label="t('games.russianRoulette.chargedBullets', { tumberSize: formValue.tumberSize })"
			>
				<n-input-number
					v-model:value="formValue.chargedBullets" :min="1"
					:max="formValue.tumberSize-1"
				/>
			</n-form-item>

			<n-form-item :label="t('games.russianRoulette.timeoutSeconds')">
				<n-input-number v-model:value="formValue.timeoutSeconds" :max="1209600" />
			</n-form-item>

			<n-form-item :label="t('games.russianRoulette.decisionSeconds')">
				<n-input-number v-model:value="formValue.decisionSeconds" :max="60" />
			</n-form-item>

			<n-divider />

			<n-form-item :label="t('games.russianRoulette.initMessage')">
				<n-input
					v-model:value="formValue.initMessage" :maxlength="450" type="textarea" autosize
					:rows="1"
				/>
			</n-form-item>

			<n-form-item :label="t('games.russianRoulette.surviveMessage')">
				<n-input
					v-model:value="formValue.surviveMessage" :maxlength="450" type="textarea" autosize
					:rows="1"
				/>
			</n-form-item>

			<n-form-item :label="t('games.russianRoulette.deathMessage')">
				<n-input
					v-model:value="formValue.deathMessage" :maxlength="450" type="textarea" autosize
					:rows="1"
				/>
			</n-form-item>
		</div>

		<n-divider />

		<n-space vertical>
			<n-button block secondary type="warning" @click="resetSettings">
				{{ t('sharedButtons.setDefaultSettings') }}
			</n-button>

			<n-button block secondary type="success" @click="save">
				{{ t('sharedButtons.save') }}
			</n-button>
		</n-space>
	</n-modal>
</template>

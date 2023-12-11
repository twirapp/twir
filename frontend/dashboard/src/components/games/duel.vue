<script setup lang="ts">
import type { DuelSettingsResponse } from '@twir/grpc/generated/api/api/games';
import {
	NModal,
	useNotification,
	NButton,
	NSwitch,
	NInput,
	NFormItem,
	NInputNumber,
	NDivider,
} from 'naive-ui';
import { ref, toRaw, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import Card from './card.vue';

import { useDuelGame } from '@/api/games/duel';
import IconDuel from '@/assets/icons/games/duel.svg?component';

const manager = useDuelGame();
const { data: settings } = manager.useSettings();
const updater = manager.useUpdate();

const formValue = ref<DuelSettingsResponse>({
	enabled: false,
	userCooldown: 0,
	globalCooldown: 0,
	resultMessage: '',
	secondsToAccept: 10,
	startMessage: '',
	timeoutSeconds: 600,
	pointsPerWin: 0,
	pointsPerLose: 0,
	bothDiePercent: 0,
	bothDieMessage: '',
});

watch(settings, (v) => {
	if (!v) return;
	formValue.value = toRaw(v);
});

const isModalOpened = ref(false);

const message = useNotification();
const { t } = useI18n();

async function save() {
	if (!formValue.value) return;
	await updater.mutateAsync(formValue.value);
	message.success({
		title: t('sharedTexts.success'),
		duration: 2500,
	});
}
</script>

<template>
	<card
		title="Duel"
		description="Duel game"
		:icon="IconDuel"
		icon-fill="#63e2b7"
		@open-settings="isModalOpened = true"
	/>

	<n-modal
		v-model:show="isModalOpened"
		:mask-closable="false"
		:segmented="true"
		preset="card"
		title="Duel"
		content-style="padding: 10px; width: 100%"
		style="width: 500px; max-width: calc(100vw - 40px)"
	>
		<div class="form">
			<n-form-item label="Enabled" label-placement="left" :show-feedback="false">
				<n-switch
					v-model:value="formValue.enabled"
				/>
			</n-form-item>

			<n-form-item label="startmessage" :show-feedback="false">
				<n-input
					v-model:value="formValue.startMessage"
					type="textarea"
					:autosize="{ minRows: 2 }"
					:maxlength="400"
				/>
			</n-form-item>

			<n-form-item label="global cooldown" :show-feedback="false">
				<n-input-number
					v-model:value="formValue.globalCooldown"
					:max="84000"
				/>
			</n-form-item>
		</div>

		<n-divider />

		<n-button
			secondary
			block
			type="success"
			@click="save"
		>
			{{ t('sharedButtons.save') }}
		</n-button>
	</n-modal>
</template>

<style scoped>
.form {
	display: flex;
	flex-direction: column;
	gap: 8px;
}
</style>

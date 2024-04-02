<script setup lang="ts">
import type { DuelSettingsResponse } from '@twir/api/messages/games/games';
import {
	NModal,
	NButton,
	NSwitch,
	NInput,
	NFormItem,
	NInputNumber,
	NDivider,
	useThemeVars,
	NSpace,
} from 'naive-ui';
import { ref, toRaw, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import Card from './card.vue';
import CommandButton from '../commandButton.vue';

import { useDuelGame } from '@/api/games/duel';
import IconDuel from '@/assets/games/duel.svg?use';
import { useNaiveDiscrete } from '@/composables/use-naive-discrete';

const manager = useDuelGame();
const { data: settings } = manager.useSettings();
const updater = manager.useUpdate();

const initialSettings: DuelSettingsResponse = {
	enabled: false,
	startMessage: '@{target}, @{initiator} challenges you to a fight. Use {duelAcceptCommandName} for next {acceptSeconds} seconds to accept the challenge.',
	resultMessage: `Sadly, @{loser} couldn't find a way to dodge the bullet and falls apart into eternal slumber.`,
	bothDieMessage: 'Unexpectedly @{initiator} and @{target} shoot each other. Only the time knows why this happened...',
	userCooldown: 0,
	globalCooldown: 0,
	secondsToAccept: 60,
	timeoutSeconds: 600,
	pointsPerWin: 0,
	pointsPerLose: 0,
	bothDiePercent: 0,
};

const formValue = ref<DuelSettingsResponse>({ ...initialSettings });

watch(settings, (v) => {
	if (!v) return;
	formValue.value = toRaw(v);
}, { immediate: true });

const isModalOpened = ref(false);

const themeVars = useThemeVars();
const { dialog, notification } = useNaiveDiscrete();
const { t } = useI18n();

async function save() {
	if (!formValue.value) return;
	await updater.mutateAsync(formValue.value);
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
		:title="t('games.duel.title')"
		:description="t('games.duel.description')"
		:icon="IconDuel"
		:icon-stroke="1"
		icon-fill="#63e2b7"
		@open-settings="isModalOpened = true"
	/>

	<n-modal
		v-model:show="isModalOpened"
		:mask-closable="false"
		:segmented="true"
		preset="card"
		:title="t('games.duel.title')"
		content-style="padding: 10px; width: 100%"
		:style="{
			width: '40vw',
			maxWidth: 'calc(100vw - 40px)',
			'--card-background': themeVars.actionColor,
			'--title-border': `1px solid ${themeVars.borderColor}`
		}"
	>
		<div class="flex flex-col gap-2">
			<n-form-item label="Enabled" label-placement="left" :show-feedback="false">
				<n-switch
					v-model:value="formValue.enabled"
				/>
			</n-form-item>

			<div class="card">
				<div class="content">
					<div class="title">
						{{ t('games.duel.commands.title') }}
					</div>
					<div class="form-item">
						<command-button name="duel" :title="t('games.duel.commands.duel')" />
						<command-button name="duel accept" :title="t('games.duel.commands.accept')" />
						<command-button name="duel stats" :title="t('games.duel.commands.stats')" />
					</div>
				</div>
			</div>


			<div class="card">
				<div class="content">
					<div class="title">
						{{ t('games.duel.cooldown.title') }}
					</div>
					<div class="form-item">
						<n-form-item
							:label="t('games.duel.cooldown.user')" :show-feedback="false"
							style="width: 45%"
						>
							<n-input-number
								v-model:value="formValue.userCooldown"
								:max="84000"
								style="width: 100%"
							/>
						</n-form-item>

						<n-form-item
							:label="t('games.duel.cooldown.global')" :show-feedback="false"
							style="width: 45%"
						>
							<n-input-number
								v-model:value="formValue.globalCooldown"
								:max="84000"
								style="width: 100%"
							/>
						</n-form-item>
					</div>
				</div>
			</div>

			<div class="card">
				<div class="content">
					<div class="title">
						{{ t('games.duel.messages.title') }}
					</div>
					<div class="form-item flex-col">
						<n-form-item
							:label="t('games.duel.messages.start.title')"
							:feedback="t('games.duel.messages.start.description', {}, {escapeParameter: false})"
						>
							<n-input
								v-model:value="formValue.startMessage"
								type="textarea"
								:autosize="{ minRows: 2 }"
								:maxlength="400"
							/>
						</n-form-item>

						<n-form-item
							:label="t('games.duel.messages.result.title')"
							:feedback="t('games.duel.messages.result.description')"
						>
							<n-input
								v-model:value="formValue.resultMessage"
								type="textarea"
								:autosize="{ minRows: 2 }"
								:maxlength="400"
							/>
						</n-form-item>

						<n-form-item
							:label="t('games.duel.messages.bothDie.title')"
							:feedback="t('games.duel.messages.bothDie.description')"
						>
							<n-input
								v-model:value="formValue.bothDieMessage"
								type="textarea"
								:autosize="{ minRows: 2 }"
								:maxlength="400"
							/>
						</n-form-item>
					</div>
				</div>
			</div>


			<div class="card">
				<div class="content">
					<div class="title">
						{{ t('games.duel.settings.title') }}
					</div>

					<div class="form-item">
						<n-form-item :label="t('games.duel.settings.secondsToAccept')" :show-feedback="false">
							<n-input-number
								v-model:value="formValue.secondsToAccept"
								:max="600"
							/>
						</n-form-item>
						<n-form-item :label="t('games.duel.settings.timeoutTime')" :show-feedback="false">
							<n-input-number
								v-model:value="formValue.timeoutSeconds"
								:max="84000"
							/>
						</n-form-item>
						<n-form-item :label="t('games.duel.settings.bothDiePercent')" :show-feedback="false">
							<n-input-number
								v-model:value="formValue.bothDiePercent"
								:max="100"
							/>
						</n-form-item>
						<n-form-item :label="t('games.duel.settings.pointsPerWin')" :show-feedback="false">
							<n-input-number
								v-model:value="formValue.pointsPerWin"
								:max="99999999"
							/>
						</n-form-item>
						<n-form-item :label="t('games.duel.settings.pointsPerLose')" :show-feedback="false">
							<n-input-number
								v-model:value="formValue.pointsPerLose"
								:max="99999999"
							/>
						</n-form-item>
					</div>
				</div>
			</div>
		</div>

		<n-divider />

		<n-space vertical>
			<n-button
				block
				secondary
				type="warning"
				@click="resetSettings"
			>
				{{ t('sharedButtons.setDefaultSettings') }}
			</n-button>

			<n-button
				secondary
				block
				type="success"
				@click="save"
			>
				{{ t('sharedButtons.save') }}
			</n-button>
		</n-space>
	</n-modal>
</template>

<style scoped>
.card {
	@apply flex flex-col gap-2 h-full rounded bg-[color:var(--card-background)];
}

.card .content {
	@apply p-1;
}

.card .content .settings {
	@apply flex flex-col gap-2 pt-1;
}

.card .title {
	@apply flex justify-between w-full pb-1 border-b-[length:var(--title-border)];
}

.card .form-item {
	@apply flex flex-wrap gap-3 p-2 w-full;
}
</style>

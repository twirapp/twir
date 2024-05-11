<script setup lang="ts">
import {
	NButton,
	NDivider,
	NFormItem,
	NInput,
	NInputNumber,
	NModal,
	NSpace,
	NSwitch,
	useThemeVars,
} from 'naive-ui'
import { ref, toRaw, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import Card from './card.vue'

import type { GamesQuery } from '@/gql/graphql'

import { useGamesApi } from '@/api/games/games.js'
import IconDuel from '@/assets/games/duel.svg?use'
import { useNaiveDiscrete } from '@/composables/use-naive-discrete'
import CommandButton from '@/features/commands/ui/command-button.vue'

const gamesApi = useGamesApi()
const { data: settings } = gamesApi.useGamesQuery()
const updater = gamesApi.useDuelMutation()

const initialSettings: GamesQuery['gamesDuel'] = {
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
}

const formValue = ref<GamesQuery['gamesDuel']>({ ...initialSettings })

watch(settings, (v) => {
	if (!v) return

	formValue.value = toRaw(v.gamesDuel)
}, { immediate: true })

const isModalOpened = ref(false)

const themeVars = useThemeVars()
const { dialog, notification } = useNaiveDiscrete()
const { t } = useI18n()

async function save() {
	if (!formValue.value) return
	await updater.executeMutation({
		opts: formValue.value,
	})
	notification.success({
		title: t('sharedTexts.saved'),
		duration: 2500,
	})
}

function resetSettings() {
	dialog.create({
		type: 'warning',
		title: t('sharedTexts.dangerZone'),
		content: t('sharedTexts.setDefaultSettings'),
		positiveText: t('sharedButtons.confirm'),
		negativeText: t('sharedButtons.close'),
		onPositiveClick: () => {
			formValue.value = initialSettings
			save()
		},
	})
}
</script>

<template>
	<Card
		:title="t('games.duel.title')"
		:description="t('games.duel.description')"
		:icon="IconDuel"
		:icon-stroke="1"
		icon-fill="#63e2b7"
		@open-settings="isModalOpened = true"
	/>

	<NModal
		v-model:show="isModalOpened"
		:mask-closable="false"
		:segmented="true"
		preset="card"
		:title="t('games.duel.title')"
		content-style="padding: 10px; width: 100%"
		:style="{
			'width': '40vw',
			'maxWidth': 'calc(100vw - 40px)',
			'--card-background': themeVars.actionColor,
			'--title-border': `1px solid ${themeVars.borderColor}`,
		}"
	>
		<div class="flex flex-col gap-2">
			<NFormItem label="Enabled" label-placement="left" :show-feedback="false">
				<NSwitch
					v-model:value="formValue.enabled"
				/>
			</NFormItem>

			<div class="card">
				<div class="content">
					<div class="title">
						{{ t('games.duel.commands.title') }}
					</div>
					<div class="form-item">
						<CommandButton name="duel" :title="t('games.duel.commands.duel')" />
						<CommandButton name="duel accept" :title="t('games.duel.commands.accept')" />
						<CommandButton name="duel stats" :title="t('games.duel.commands.stats')" />
					</div>
				</div>
			</div>

			<div class="card">
				<div class="content">
					<div class="title">
						{{ t('games.duel.cooldown.title') }}
					</div>
					<div class="form-item">
						<NFormItem
							:label="t('games.duel.cooldown.user')" :show-feedback="false"
							style="width: 45%"
						>
							<NInputNumber
								v-model:value="formValue.userCooldown"
								:max="84000"
								style="width: 100%"
							/>
						</NFormItem>

						<NFormItem
							:label="t('games.duel.cooldown.global')" :show-feedback="false"
							style="width: 45%"
						>
							<NInputNumber
								v-model:value="formValue.globalCooldown"
								:max="84000"
								style="width: 100%"
							/>
						</NFormItem>
					</div>
				</div>
			</div>

			<div class="card">
				<div class="content">
					<div class="title">
						{{ t('games.duel.messages.title') }}
					</div>
					<div class="form-item flex-col">
						<NFormItem
							:label="t('games.duel.messages.start.title')"
							:feedback="t('games.duel.messages.start.description', {}, { escapeParameter: false })"
						>
							<NInput
								v-model:value="formValue.startMessage"
								type="textarea"
								:autosize="{ minRows: 2 }"
								:maxlength="400"
							/>
						</NFormItem>

						<NFormItem
							:label="t('games.duel.messages.result.title')"
							:feedback="t('games.duel.messages.result.description')"
						>
							<NInput
								v-model:value="formValue.resultMessage"
								type="textarea"
								:autosize="{ minRows: 2 }"
								:maxlength="400"
							/>
						</NFormItem>

						<NFormItem
							:label="t('games.duel.messages.bothDie.title')"
							:feedback="t('games.duel.messages.bothDie.description')"
						>
							<NInput
								v-model:value="formValue.bothDieMessage"
								type="textarea"
								:autosize="{ minRows: 2 }"
								:maxlength="400"
							/>
						</NFormItem>
					</div>
				</div>
			</div>

			<div class="card">
				<div class="content">
					<div class="title">
						{{ t('games.duel.settings.title') }}
					</div>

					<div class="form-item">
						<NFormItem :label="t('games.duel.settings.secondsToAccept')" :show-feedback="false">
							<NInputNumber
								v-model:value="formValue.secondsToAccept"
								:max="600"
							/>
						</NFormItem>
						<NFormItem :label="t('games.duel.settings.timeoutTime')" :show-feedback="false">
							<NInputNumber
								v-model:value="formValue.timeoutSeconds"
								:max="84000"
							/>
						</NFormItem>
						<NFormItem :label="t('games.duel.settings.bothDiePercent')" :show-feedback="false">
							<NInputNumber
								v-model:value="formValue.bothDiePercent"
								:max="100"
							/>
						</NFormItem>
						<NFormItem :label="t('games.duel.settings.pointsPerWin')" :show-feedback="false">
							<NInputNumber
								v-model:value="formValue.pointsPerWin"
								:max="99999999"
							/>
						</NFormItem>
						<NFormItem :label="t('games.duel.settings.pointsPerLose')" :show-feedback="false">
							<NInputNumber
								v-model:value="formValue.pointsPerLose"
								:max="99999999"
							/>
						</NFormItem>
					</div>
				</div>
			</div>
		</div>

		<NDivider />

		<NSpace vertical>
			<NButton
				block
				secondary
				type="warning"
				@click="resetSettings"
			>
				{{ t('sharedButtons.setDefaultSettings') }}
			</NButton>

			<NButton
				secondary
				block
				type="success"
				@click="save"
			>
				{{ t('sharedButtons.save') }}
			</NButton>
		</NSpace>
	</NModal>
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

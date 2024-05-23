<script setup lang="ts">
import { IconBomb } from '@tabler/icons-vue'
import {
	NButton,
	NDivider,
	NFormItem,
	NInput,
	NInputNumber,
	NModal,
	NSpace,
	NSwitch,
} from 'naive-ui'
import { ref, toRaw, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import Card from './card.vue'

import type { GamesQuery } from '@/gql/graphql'

import { useGamesApi } from '@/api/games/games'
import { useNaiveDiscrete } from '@/composables/use-naive-discrete'
import CommandButton from '@/features/commands/components/command-button.vue'

const isModalOpened = ref(false)

const gamesManager = useGamesApi()
const { data: settings } = gamesManager.useGamesQuery()
const updater = gamesManager.useRussianRouletteMutation()

const initialSettings: GamesQuery['gamesRussianRoulette'] = {
	enabled: false,
	canBeUsedByModerator: false,
	timeoutSeconds: 60,
	decisionSeconds: 2,
	chargedBullets: 1,
	initMessage: '{sender} has initiated a game of roulette. Is luck on their side?',
	surviveMessage: '{sender} survives the game of roulette! Luck smiles upon them.',
	deathMessage: `{sender} couldn't make it through the game of roulette. Unfortunately, luck wasn't on their side this time.`,
	tumberSize: 6,
}

const formValue = ref<GamesQuery['gamesRussianRoulette']>({ ...initialSettings })

watch(settings, (v) => {
	if (!v) return

	const raw = toRaw(v)

	formValue.value = structuredClone(raw.gamesRussianRoulette)
}, { immediate: true })

const { t } = useI18n()

const { dialog, notification } = useNaiveDiscrete()

async function save() {
	await updater.executeMutation({ opts: formValue.value })
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
		title="Russian Roulette"
		:icon="IconBomb"
		:icon-stroke="1"
		:description="t('games.russianRoulette.description')"
		@open-settings="isModalOpened = true"
	/>

	<NModal
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
				<NSwitch v-model:value="formValue.enabled" />
			</div>

			<CommandButton name="roulette" />
		</div>

		<NDivider />

		<div class="flex flex-col gap-2 mt-[10px]">
			<NFormItem :label="t('games.russianRoulette.canBeUsedByModerator')">
				<NSwitch v-model:value="formValue.canBeUsedByModerator" />
			</NFormItem>

			<NFormItem :label="t('games.russianRoulette.tumberSize')">
				<NInputNumber v-model:value="formValue.tumberSize" :min="2" :max="100" />
			</NFormItem>

			<NFormItem
				:label="t('games.russianRoulette.chargedBullets', { tumberSize: formValue.tumberSize })"
			>
				<NInputNumber
					v-model:value="formValue.chargedBullets" :min="1"
					:max="formValue.tumberSize - 1"
				/>
			</NFormItem>

			<NFormItem :label="t('games.russianRoulette.timeoutSeconds')">
				<NInputNumber v-model:value="formValue.timeoutSeconds" :max="86400" />
			</NFormItem>

			<NFormItem :label="t('games.russianRoulette.decisionSeconds')">
				<NInputNumber v-model:value="formValue.decisionSeconds" :max="60" />
			</NFormItem>

			<NDivider />

			<NFormItem :label="t('games.russianRoulette.initMessage')">
				<NInput
					v-model:value="formValue.initMessage" :maxlength="450" type="textarea" autosize
					:rows="1"
				/>
			</NFormItem>

			<NFormItem :label="t('games.russianRoulette.surviveMessage')">
				<NInput
					v-model:value="formValue.surviveMessage" :maxlength="450" type="textarea" autosize
					:rows="1"
				/>
			</NFormItem>

			<NFormItem :label="t('games.russianRoulette.deathMessage')">
				<NInput
					v-model:value="formValue.deathMessage" :maxlength="450" type="textarea" autosize
					:rows="1"
				/>
			</NFormItem>
		</div>

		<NDivider />

		<NSpace vertical>
			<NButton block secondary type="warning" @click="resetSettings">
				{{ t('sharedButtons.setDefaultSettings') }}
			</NButton>

			<NButton block secondary type="success" @click="save">
				{{ t('sharedButtons.save') }}
			</NButton>
		</NSpace>
	</NModal>
</template>

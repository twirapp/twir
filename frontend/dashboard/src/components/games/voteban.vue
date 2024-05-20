<script setup lang="ts">
import { GavelIcon } from 'lucide-vue-next'
import {
	NButton,
	NDivider,
	NFormItem,
	NInput,
	NInputNumber,
	NModal,
	NSelect,
	NSwitch,
	useMessage,
} from 'naive-ui'
import { ref, toRaw, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import Card from './card.vue'

import type { VotebanGame } from '@/gql/graphql'

import { useGamesApi } from '@/api/games/games.js'
import CommandButton from '@/features/commands/components/command-button.vue'
import { VoteBanGameVotingMode } from '@/gql/graphql'

const isModalOpened = ref(false)

const gamesManager = useGamesApi()
const { data } = gamesManager.useGamesQuery()
const updater = gamesManager.useVotebanMutation()

const maxWords = 10

const formValue = ref<VotebanGame>({
	banMessage: '',
	banMessageModerators: '',
	chatVotesWordsNegative: [],
	chatVotesWordsPositive: [],
	enabled: false,
	initMessage: '',
	neededVotes: 3,
	surviveMessage: '',
	surviveMessageModerators: '',
	timeoutModerators: true,
	timeoutSeconds: 60,
	voteDuration: 60,
	votingMode: VoteBanGameVotingMode.Chat,
})

watch(data, (v) => {
	if (!v) return

	const raw = toRaw(v)

	formValue.value = structuredClone(raw.gamesVoteban)
}, { immediate: true })

const { t } = useI18n()

const notifications = useMessage()

async function save() {
	await updater.executeMutation({
		opts: formValue.value,
	})
	notifications.success(t('sharedTexts.saved'))
}
</script>

<template>
	<Card
		title="Voteban"
		:icon="GavelIcon"
		:icon-stroke="1"
		:description="t('games.voteban.description')"
		@open-settings="isModalOpened = true"
	/>

	<NModal
		v-model:show="isModalOpened"
		:mask-closable="false"
		:segmented="true"
		preset="card"
		title="Voteban"
		content-style="padding: 10px; width: 100%"
		style="width: 500px; max-width: calc(100vw - 40px)"
	>
		<div class="flex gap-6">
			<div class="flex flex-col gap-1 items-center">
				<span>{{ t('sharedTexts.enabled') }}</span>
				<NSwitch v-model:value="formValue.enabled"></NSwitch>
			</div>

			<CommandButton name="voteban" />
		</div>

		<NDivider title-placement="center">
			Messages
		</NDivider>

		<div class="flex flex-col gap-2">
			<NFormItem :label="t('games.voteban.initialMessage')">
				<NInput v-model:value="formValue.initMessage" type="textarea" autosize :maxlength="500" show-count />
			</NFormItem>

			<NFormItem :label="t('games.voteban.banMessage')">
				<NInput v-model:value="formValue.banMessage" type="textarea" autosize :maxlength="500" show-count />
			</NFormItem>

			<NFormItem :label="t('games.voteban.surviveMessage')">
				<NInput v-model:value="formValue.surviveMessage" type="textarea" autosize :maxlength="500" show-count />
			</NFormItem>
		</div>

		<NDivider title-placement="center">
			{{ t('sharedTexts.settings') }}
		</NDivider>

		<div class="flex flex-col gap-2">
			<NFormItem :label="t('games.voteban.voteMode')">
				<NSelect
					v-model:value="formValue.votingMode"
					:options="[
						{ label: 'Chat', value: VoteBanGameVotingMode.Chat },
						{ label: 'Twitch polls (soon)', value: VoteBanGameVotingMode.Polls, disabled: true },
					]"
				/>
			</NFormItem>

			<NFormItem :label="t('games.voteban.voteDuration')">
				<NInputNumber v-model:value="formValue.voteDuration" style="width: 100%" :min="1" :max="84600" />
			</NFormItem>

			<NFormItem :label="t('games.voteban.neededVotes')">
				<NInputNumber v-model:value="formValue.neededVotes" style="width: 100%" :min="1" :max="999999" />
			</NFormItem>

			<NFormItem :label="t('games.voteban.banDuration')">
				<NInputNumber v-model:value="formValue.timeoutSeconds" style="width: 100%" :min="1" :max="84600" />
			</NFormItem>
		</div>

		<NDivider title-placement="center">
			Moderators
		</NDivider>

		<div class="flex flex-col gap-2">
			<div class="flex justify-between w-full">
				<span>{{ t('games.voteban.timeoutModerators') }}</span>
				<NSwitch v-model:value="formValue.timeoutModerators" />
			</div>

			<NFormItem :label="t('games.voteban.banMessageModerators')">
				<NInput v-model:value="formValue.banMessageModerators" type="textarea" autosize :maxlength="500" show-count />
			</NFormItem>

			<NFormItem :label="t('games.voteban.surviveMessageModerators')">
				<NInput v-model:value="formValue.surviveMessageModerators" type="textarea" autosize :maxlength="500" show-count />
			</NFormItem>
		</div>

		<NDivider />

		<NButton block secondary type="success" @click="save">
			{{ t('sharedButtons.save') }}
		</NButton>
	</NModal>
</template>

<script setup lang="ts">
import { IconMessageCircleQuestion, IconTrash } from '@tabler/icons-vue'
import { NButton, NDivider, NInput, NModal, NSwitch, useMessage } from 'naive-ui'
import { ref, toRaw, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import Card from './card.vue'

import { useGamesApi } from '@/api/games/games.js'
import CommandButton from '@/features/commands/components/command-button.vue'

const isModalOpened = ref(false)

const maxAnswers = 25

const gamesManager = useGamesApi()
const { data } = gamesManager.useGamesQuery()
const updater = gamesManager.useEightBallMutation()

const formValue = ref({
	enabled: false,
	answers: ['Yes', 'No'],
})

watch(data, (v) => {
	if (!v) return

	const raw = toRaw(v)
	formValue.value.answers = raw.gamesEightBall.answers
	formValue.value.enabled = raw.gamesEightBall.enabled
}, { immediate: true })

const { t } = useI18n()

const notifications = useMessage()

async function save() {
	await updater.executeMutation({
		opts: {
			answers: formValue.value.answers,
			enabled: formValue.value.enabled,
		},
	})
	notifications.success(t('sharedTexts.saved'))
}
</script>

<template>
	<Card
		title="8ball"
		:icon="IconMessageCircleQuestion"
		:icon-stroke="1"
		:description="t('games.8ball.description')"
		@open-settings="isModalOpened = true"
	/>

	<NModal
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
				<NSwitch v-model:value="formValue.enabled"></NSwitch>
			</div>

			<CommandButton name="8ball" />
		</div>

		<NDivider />

		<h3>{{ t('games.8ball.answers') }} ({{ formValue.answers.length }}/{{ maxAnswers }})</h3>

		<div class="flex flex-col gap-2">
			<div
				v-for="(_, index) of formValue.answers"
				:key="index"
				class="flex gap-1"
			>
				<NInput
					v-model:value="formValue.answers[index]"
					placeholder="Yes"
				/>

				<NButton
					secondary
					type="error"
					@click="() => {
						formValue.answers = formValue.answers.filter((_, i) => i != index)
					}"
				>
					<IconTrash />
				</NButton>
			</div>

			<NButton
				secondary
				type="info"
				block
				:disabled="formValue.answers.length >= maxAnswers"
				@click="() => formValue.answers.push('')"
			>
				{{ t('sharedButtons.create') }}
			</NButton>
		</div>

		<NDivider />

		<NButton block secondary type="success" @click="save">
			{{ t('sharedButtons.save') }}
		</NButton>
	</NModal>
</template>

<script setup lang="ts">
import { SkullIcon } from 'lucide-vue-next'
import {
	NButton,
	NDivider,
	NFormItem,
	NInput,
	NInputNumber,
	NModal,
	NSwitch,
	useMessage,
} from 'naive-ui'
import { ref, toRaw, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import Card from './card.vue'

import type { SeppukuGame } from '@/gql/graphql'

import { useGamesApi } from '@/api/games/games.js'
import CommandButton from '@/features/commands/components/command-button.vue'

const isModalOpened = ref(false)

const gamesManager = useGamesApi()
const { data } = gamesManager.useGamesQuery()
const updater = gamesManager.useSeppukuMutation()

const formValue = ref<SeppukuGame>({
	enabled: false,
	message: '{sender} said: my honor tarnished, I reclaim it through death. May my spirit find peace. Farewell.',
	messageModerators: '{sender} drew his sword and ripped open his belly for the sad emperor.',
	timeoutModerators: false,
	timeoutSeconds: 60,
})

watch(data, (v) => {
	if (!v) return

	const raw = toRaw(v)

	formValue.value.enabled = raw.gamesSeppuku.enabled
	formValue.value.message = raw.gamesSeppuku.message
	formValue.value.messageModerators = raw.gamesSeppuku.messageModerators
	formValue.value.timeoutModerators = raw.gamesSeppuku.timeoutModerators
	formValue.value.timeoutSeconds = raw.gamesSeppuku.timeoutSeconds
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
		title="Seppuku"
		:icon="SkullIcon"
		:icon-stroke="1"
		:description="t('games.seppuku.description')"
		@open-settings="isModalOpened = true"
	/>

	<NModal
		v-model:show="isModalOpened"
		:mask-closable="false"
		:segmented="true"
		preset="card"
		title="Seppuku"
		content-style="padding: 10px; width: 100%"
		style="width: 500px; max-width: calc(100vw - 40px)"
	>
		<div class="flex gap-6">
			<div class="flex flex-col gap-1 items-center">
				<span>{{ t('sharedTexts.enabled') }}</span>
				<NSwitch v-model:value="formValue.enabled" />
			</div>

			<CommandButton name="seppuku" />
		</div>

		<NDivider />

		<NFormItem :label="t('games.seppuku.message')">
			<NInput
				v-model:value="formValue.message"
				clearable
			/>
		</NFormItem>

		<NFormItem :label="t('games.seppuku.timeoutSeconds')">
			<NInputNumber
				v-model:value="formValue.timeoutSeconds"
				min="1"
				max="84600"
			/>
		</NFormItem>

		<NFormItem :label="t('games.seppuku.timeoutModerators')">
			<NSwitch
				v-model:value="formValue.timeoutModerators"
			/>
		</NFormItem>

		<NFormItem :label="t('games.seppuku.messageModerators')">
			<NInput
				v-model:value="formValue.messageModerators"
				clearable
				:disabled="!formValue.timeoutModerators"
			/>
		</NFormItem>

		<NButton block secondary type="success" @click="save">
			{{ t('sharedButtons.save') }}
		</NButton>
	</NModal>
</template>

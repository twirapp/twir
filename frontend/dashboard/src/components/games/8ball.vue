<script setup lang="ts">
import { MessageCircle, Trash } from 'lucide-vue-next'
import { computed, ref, toRaw, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import { useProfile } from '@/api'
import { useGamesApi } from '@/api/games/games.js'
import Card from '@/components/games/card.vue'
import { Button } from '@/components/ui/button'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Separator } from '@/components/ui/separator'
import { Switch } from '@/components/ui/switch'
import { toast } from 'vue-sonner'
import CommandButton from '@/features/commands/ui/command-button.vue'

const isModalOpened = ref(false)
const { data: profile } = useProfile()

const maxAnswers = computed(() => {
	const selectedDashboard = profile.value?.availableDashboards.find(
		(d) => d.id === profile.value?.selectedDashboardId
	)
	return selectedDashboard?.plan.maxEightballAnswers ?? 25
})

const gamesManager = useGamesApi()
const { data } = gamesManager.useGamesQuery()
const updater = gamesManager.useEightBallMutation()

const formValue = ref({
	enabled: false,
	answers: ['Yes', 'No'],
})

watch(
	data,
	(v) => {
		if (!v) return

		const raw = toRaw(v)
		formValue.value.answers = raw.gamesEightBall.answers
		formValue.value.enabled = raw.gamesEightBall.enabled
	},
	{ immediate: true }
)

const { t } = useI18n()

async function save() {
	await updater.executeMutation({
		opts: {
			answers: formValue.value.answers,
			enabled: formValue.value.enabled,
		},
	})
	toast.success(t('sharedTexts.saved'))
}
</script>

<template>
	<Card
		title="8ball"
		:icon="MessageCircle"
		:icon-stroke="1"
		:description="t('games.8ball.description')"
		show-settings
		@open-settings="isModalOpened = true"
	/>

	<Dialog v-model:open="isModalOpened">
		<DialogContent class="sm:max-w-[500px]">
			<DialogHeader>
				<DialogTitle>8ball</DialogTitle>
			</DialogHeader>

			<div class="flex flex-col gap-3">
				<div class="flex flex-row gap-1 items-center">
					<span>{{ t('sharedTexts.enabled') }}</span>
					<Switch
						:model-value="formValue.enabled"
						@update:model-value="() => (formValue.enabled = !formValue.enabled)"
					/>
				</div>

				<CommandButton name="8ball" />
			</div>

			<Separator />

			<div class="space-y-4">
				<h3 class="font-medium">
					{{ t('games.8ball.answers') }} ({{ formValue.answers.length }}/{{ maxAnswers }})
				</h3>

				<div class="space-y-2">
					<div v-for="(_, index) of formValue.answers" :key="index" class="flex gap-2">
						<Input v-model="formValue.answers[index]" placeholder="Yes" class="flex-1" />

						<Button
							variant="destructive"
							size="icon"
							@click="
								() => {
									formValue.answers = formValue.answers.filter((_, i) => i != index)
								}
							"
						>
							<Trash class="h-4 w-4" />
						</Button>
					</div>

					<Button
						variant="secondary"
						class="w-full"
						:disabled="formValue.answers.length >= maxAnswers"
						@click="() => formValue.answers.push('')"
					>
						{{ t('sharedButtons.create') }}
					</Button>
				</div>
			</div>

			<Separator />

			<Button variant="default" class="w-full" @click="save">
				{{ t('sharedButtons.save') }}
			</Button>
		</DialogContent>
	</Dialog>
</template>

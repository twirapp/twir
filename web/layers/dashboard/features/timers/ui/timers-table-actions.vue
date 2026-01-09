<script setup lang="ts">
import { PencilIcon, TrashIcon } from 'lucide-vue-next'
import { ref } from 'vue'

import type { TimerResponse } from '#layers/dashboard/api/timers'

import { useTimersApi } from '#layers/dashboard/api/timers'




const props = defineProps<{ timer: TimerResponse }>()
const showDelete = ref(false)

const timersApi = useTimersApi()
const updateMutation = timersApi.useMutationUpdateTimer()
const removeMutation = timersApi.useMutationRemoveTimer()

function toggleEnabledGreetings() {
	updateMutation.executeMutation({
		id: props.timer.id,
		opts: { enabled: !props.timer.enabled },
	})
}

function deleteTimer() {
	removeMutation.executeMutation({ id: props.timer.id })
}
</script>

<template>
	<div class="flex items-center gap-2">
		<UiSwitch :model-value="timer.enabled" @update:model-value="toggleEnabledGreetings" />

		<RouterLink v-slot="{ navigate, href }" custom :to="`/dashboard/timers/${timer.id}`">
			<UiButton
				as="a"
				:href="href"
				variant="secondary"
				size="icon"
				@click="navigate"
			>
				<PencilIcon class="h-4 w-4" />
			</UiButton>
		</RouterLink>

		<UiButton variant="destructive" size="icon" @click="showDelete = true">
			<TrashIcon class="h-4 w-4" />
		</UiButton>
	</div>

	<UiActionConfirm v-model:open="showDelete" @confirm="deleteTimer" />
</template>

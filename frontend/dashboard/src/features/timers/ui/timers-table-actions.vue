<script setup lang="ts">
import { PencilIcon, TrashIcon } from 'lucide-vue-next'
import { ref } from 'vue'

import type { TimerResponse } from '@/api/timers'

import { useTimersApi } from '@/api/timers'
import ActionConfirm from '@/components/ui/action-confirm.vue'
import { Button } from '@/components/ui/button'
import { Switch } from '@/components/ui/switch'

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
		<Switch :checked="timer.enabled" @update:checked="toggleEnabledGreetings" />

		<RouterLink v-slot="{ navigate, href }" custom :to="`/dashboard/timers/${timer.id}`">
			<Button
				as="a"
				:href="href"
				variant="secondary"
				size="icon"
				@click="navigate"
			>
				<PencilIcon class="h-4 w-4" />
			</Button>
		</RouterLink>

		<Button variant="destructive" size="icon" @click="showDelete = true">
			<TrashIcon class="h-4 w-4" />
		</Button>
	</div>

	<ActionConfirm v-model:open="showDelete" @confirm="deleteTimer" />
</template>

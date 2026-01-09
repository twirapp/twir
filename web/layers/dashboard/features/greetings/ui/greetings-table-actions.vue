<script setup lang="ts">
import { PencilIcon, TrashIcon } from 'lucide-vue-next'
import { ref } from 'vue'

import GreetingsDialog from './greetings-dialog.vue'

import { type Greetings, useGreetingsApi } from '#layers/dashboard/api/greetings'




const props = defineProps<{ greetings: Greetings }>()
const showDelete = ref(false)

const greetingsApi = useGreetingsApi()
const updateGreetingsMutation = greetingsApi.useMutationUpdateGreetings()
const removeGreetingsMutation = greetingsApi.useMutationRemoveGreetings()

function toggleEnabledGreetings() {
	updateGreetingsMutation.executeMutation({
		id: props.greetings.id,
		opts: { enabled: !props.greetings.enabled },
	})
}

function deleteGreetings() {
	removeGreetingsMutation.executeMutation({ id: props.greetings.id })
}
</script>

<template>
	<div class="flex items-center gap-2">
		<UiSwitch :model-value="greetings.enabled" @update:model-value="toggleEnabledGreetings" />

		<GreetingsDialog :greeting="greetings">
			<template #dialog-trigger>
				<UiButton variant="secondary" size="icon">
					<PencilIcon class="h-4 w-4" />
				</UiButton>
			</template>
		</GreetingsDialog>

		<UiButton variant="destructive" size="icon" @click="showDelete = true">
			<TrashIcon class="h-4 w-4" />
		</UiButton>
	</div>

	<UiActionConfirm v-model:open="showDelete" @confirm="deleteGreetings" />
</template>

<script setup lang="ts">
import { useField } from 'vee-validate'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import type { EventOperation } from '@/api/events.ts'

import { useObsOverlayManager } from '@/api'
import {
	FormControl,
	FormField,
	FormItem,
	FormLabel,
	FormMessage,
} from '@/components/ui/form'
import {
	Select,
	SelectContent,
	SelectGroup,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from '@/components/ui/select'
import { EventOperationType } from '@/gql/graphql.ts'

const props = defineProps<{
	operationIndex: number
}>()

const currentOperationPath = computed(() => `operations.${props.operationIndex}`)
const { value: currentOperation } = useField<Omit<EventOperation, 'id'> | undefined>(currentOperationPath)

const { t } = useI18n()
const obsManager = useObsOverlayManager()
const obsSettings = obsManager.getSettings()

const obsScenes = computed(() => {
	return obsSettings.data.value?.scenes.map(s => ({
		value: s,
		label: s,
	})) ?? []
})

const obsSources = computed(() => {
	return obsSettings.data.value?.sources.map(s => ({
		value: s,
		label: s,
	})) ?? []
})

const obsAudioSources = computed(() => {
	return obsSettings.data.value?.audioSources.map(s => ({
		value: s,
		label: s,
	})) ?? []
})

const data = computed(() => {
	if (!currentOperation.value?.type) return

	if ([
		EventOperationType.ObsSetAudioVolume,
		EventOperationType.ObsIncreaseAudioVolume,
		EventOperationType.ObsDecreaseAudioVolume,
		EventOperationType.ObsEnableAudio,
		EventOperationType.ObsDisableAudio,
		EventOperationType.ObsToggleAudio,
	].includes(currentOperation.value?.type)) return obsAudioSources.value

	if (currentOperation.value.type === EventOperationType.ObsChangeScene) return obsScenes.value
	if (currentOperation.value.type === EventOperationType.ObsToggleSource) return obsSources.value

	return []
})
</script>

<template>
	<FormField
		v-slot="{ componentField }"
		:name="`operations.${operationIndex}.target`"
	>
		<FormItem>
			<FormLabel>Select source</FormLabel>
			<FormControl>
				<Select
					v-bind="componentField"
					:placeholder="t('events.targetAlertPlaceholder')"
				>
					<SelectTrigger>
						<SelectValue placeholder="Select" />
					</SelectTrigger>
					<SelectContent>
						<SelectGroup>
							<SelectItem v-for="item in data" :key="item.value" :value="item.value">
								{{ item.label }}
							</SelectItem>
						</SelectGroup>
					</SelectContent>
				</Select>
			</FormControl>
			<FormMessage />
		</FormItem>
	</FormField>
</template>

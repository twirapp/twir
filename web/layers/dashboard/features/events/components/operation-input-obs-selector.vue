<script setup lang="ts">
import { useField } from 'vee-validate'
import { computed } from 'vue'


import type { EventOperation } from '#layers/dashboard/api/events.ts'

import { useObsWebsocketApi } from '#layers/dashboard/api/overlays-obs'


import { EventOperationType } from '~/gql/graphql.ts'

const props = defineProps<{
	operationIndex: number
}>()

const currentOperationPath = computed(() => `operations.${props.operationIndex}`)
const { value: currentOperation } = useField<Omit<EventOperation, 'id'> | undefined>(
	currentOperationPath
)

const { t } = useI18n()
const obsApi = useObsWebsocketApi()
const { data: queryData } = obsApi.useQueryObsWebsocket()
const { data: subscriptionData } = obsApi.useSubscriptionObsData()

// Use subscription data for realtime updates, fallback to query data
const obsData = computed(() => {
	return subscriptionData.value?.obsWebsocketData ?? queryData.value?.obsWebsocketData
})

const obsScenes = computed(() => {
	return (
		obsData.value?.scenes?.map((s: string) => ({
			value: s,
			label: s,
		})) ?? []
	)
})

const obsSources = computed(() => {
	return (
		obsData.value?.sources?.map((s: string) => ({
			value: s,
			label: s,
		})) ?? []
	)
})

const obsAudioSources = computed(() => {
	return (
		obsData.value?.audioSources?.map((s: string) => ({
			value: s,
			label: s,
		})) ?? []
	)
})

const data = computed(() => {
	if (!currentOperation.value?.type) return

	if (
		[
			EventOperationType.ObsSetAudioVolume,
			EventOperationType.ObsIncreaseAudioVolume,
			EventOperationType.ObsDecreaseAudioVolume,
			EventOperationType.ObsEnableAudio,
			EventOperationType.ObsDisableAudio,
			EventOperationType.ObsToggleAudio,
		].includes(currentOperation.value?.type)
	)
		return obsAudioSources.value

	if (currentOperation.value.type === EventOperationType.ObsChangeScene) return obsScenes.value
	if (currentOperation.value.type === EventOperationType.ObsToggleSource) return obsSources.value

	return []
})
</script>

<template>
	<UiFormField v-slot="{ componentField }" :name="`operations.${operationIndex}.target`">
		<UiFormItem>
			<UiFormLabel>Select source</UiFormLabel>
			<UiFormControl>
				<UiSelect v-bind="componentField" :placeholder="t('events.targetAlertPlaceholder')">
					<UiSelectTrigger>
						<UiSelectValue placeholder="Select" />
					</UiSelectTrigger>
					<UiSelectContent>
						<UiSelectGroup>
							<UiSelectItem v-for="item in data" :key="item.value" :value="item.value">
								{{ item.label }}
							</UiSelectItem>
						</UiSelectGroup>
					</UiSelectContent>
				</UiSelect>
			</UiFormControl>
			<UiFormMessage />
		</UiFormItem>
	</UiFormField>
</template>

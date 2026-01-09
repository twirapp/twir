<script setup lang="ts">


import { useAlertsApi } from '#layers/dashboard/api/alerts.ts'



defineProps<{
	operationIndex: number
}>()

const { t } = useI18n()

const alertsApi = useAlertsApi()

const { data } = alertsApi.useAlertsQuery()
</script>

<template>
	<UiFormField
		v-slot="{ componentField }"
		:name="`operations.${operationIndex}.target`"
	>
		<UiFormItem>
			<UiFormLabel>{{ t('alerts.name') }}</UiFormLabel>
			<UiFormControl>
				<UiSelect
					v-bind="componentField"
					:placeholder="t('events.targetAlertPlaceholder')"
				>
					<UiSelectTrigger>
						<UiSelectValue placeholder="Select alert" />
					</UiSelectTrigger>
					<UiSelectContent>
						<UiSelectGroup>
							<UiSelectItem v-for="alert in data?.channelAlerts" :key="alert.id" :value="alert.id">
								{{ alert.name }}
							</UiSelectItem>
						</UiSelectGroup>
					</UiSelectContent>
				</UiSelect>
			</UiFormControl>
			<UiFormMessage />
		</UiFormItem>
	</UiFormField>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'

import { useAlertsApi } from '@/api/alerts.ts'
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

defineProps<{
	operationIndex: number
}>()

const { t } = useI18n()

const alertsApi = useAlertsApi()

const { data } = alertsApi.useAlertsQuery()
</script>

<template>
	<FormField
		v-slot="{ componentField }"
		:name="`operations.${operationIndex}.target`"
	>
		<FormItem>
			<FormLabel>{{ t('alerts.name') }}</FormLabel>
			<FormControl>
				<Select
					v-bind="componentField"
					:placeholder="t('events.targetAlertPlaceholder')"
				>
					<SelectTrigger>
						<SelectValue placeholder="Select alert" />
					</SelectTrigger>
					<SelectContent>
						<SelectGroup>
							<SelectItem v-for="alert in data?.channelAlerts" :key="alert.id" :value="alert.id">
								{{ alert.name }}
							</SelectItem>
						</SelectGroup>
					</SelectContent>
				</Select>
			</FormControl>
			<FormMessage />
		</FormItem>
	</FormField>
</template>

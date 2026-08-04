<script setup lang="ts">
import type { ImportReportData } from './types.js'

import { Button } from '@/components/ui/button/index.js'
import {
	Card,
	CardContent,
	CardDescription,
	CardFooter,
	CardHeader,
	CardTitle,
} from '@/components/ui/card/index.js'
import ImportResult from './import-result.vue'

withDefaults(
	defineProps<{
		connected: boolean
		canManageCommands: boolean
		canManageTimers: boolean
		commandsImporting?: boolean
		timersImporting?: boolean
		commandsReport?: ImportReportData | null
		timersReport?: ImportReportData | null
	}>(),
	{
		commandsImporting: false,
		timersImporting: false,
		commandsReport: null,
		timersReport: null,
	}
)

defineEmits<{
	importCommands: []
	importTimers: []
}>()

const { t } = useI18n()
</script>

<template>
	<div class="grid gap-4 md:grid-cols-2">
		<Card class="flex flex-col">
			<CardHeader>
				<CardTitle>{{ t('imports.commands.title') }}</CardTitle>
				<CardDescription>{{ t('imports.commands.description') }}</CardDescription>
			</CardHeader>
			<CardContent class="grow">
				<ImportResult v-if="commandsReport" :report="commandsReport" />
				<p v-else class="text-sm text-muted-foreground">
					{{ t('imports.result.waiting') }}
				</p>
			</CardContent>
			<CardFooter>
				<Button
					class="w-full"
					:disabled="!connected || !canManageCommands || commandsImporting"
					data-test="import-commands"
					@click="$emit('importCommands')"
				>
					<Icon v-if="commandsImporting" name="lucide:loader-circle" data-icon="inline-start" class="animate-spin" />
					<Icon v-else name="lucide:download" data-icon="inline-start" />
					{{ commandsImporting ? t('imports.actions.importing') : t('imports.actions.importCommands') }}
				</Button>
			</CardFooter>
		</Card>

		<Card class="flex flex-col">
			<CardHeader>
				<CardTitle>{{ t('imports.timers.title') }}</CardTitle>
				<CardDescription>{{ t('imports.timers.description') }}</CardDescription>
			</CardHeader>
			<CardContent class="grow">
				<ImportResult v-if="timersReport" :report="timersReport" />
				<p v-else class="text-sm text-muted-foreground">
					{{ t('imports.result.waiting') }}
				</p>
			</CardContent>
			<CardFooter>
				<Button
					class="w-full"
					:disabled="!connected || !canManageTimers || timersImporting"
					data-test="import-timers"
					@click="$emit('importTimers')"
				>
					<Icon v-if="timersImporting" name="lucide:loader-circle" data-icon="inline-start" class="animate-spin" />
					<Icon v-else name="lucide:download" data-icon="inline-start" />
					{{ timersImporting ? t('imports.actions.importing') : t('imports.actions.importTimers') }}
				</Button>
			</CardFooter>
		</Card>
	</div>
</template>

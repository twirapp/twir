<script setup lang="ts">
import type { ImportReportData } from './types.js'

import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert/index.js'
import { Badge } from '@/components/ui/badge/index.js'
import { ImportFailureReason } from '~/gql/graphql.js'

defineProps<{
	report: ImportReportData
}>()

const { t } = useI18n()

const reasonKeys: Record<ImportFailureReason, string> = {
	[ImportFailureReason.Duplicate]: 'duplicate',
	[ImportFailureReason.PlanLimit]: 'planLimit',
	[ImportFailureReason.UnsupportedRole]: 'unsupportedRole',
	[ImportFailureReason.UnsupportedResponseType]: 'unsupportedResponseType',
	[ImportFailureReason.IncompatibleIntervals]: 'incompatibleIntervals',
	[ImportFailureReason.InvalidRecord]: 'invalidRecord',
}
</script>

<template>
	<Alert>
		<Icon name="lucide:circle-check" />
		<AlertTitle>
			{{ report.failedCount ? t('imports.result.partial') : t('imports.result.success') }}
		</AlertTitle>
		<AlertDescription class="flex flex-col gap-3">
			<div class="flex flex-wrap gap-2">
				<Badge variant="secondary">
					{{ t('imports.result.imported', { count: report.importedCount }) }}
				</Badge>
				<Badge :variant="report.failedCount ? 'destructive' : 'outline'">
					{{ t('imports.result.failed', { count: report.failedCount }) }}
				</Badge>
			</div>

			<ul v-if="report.failures.length" class="flex max-h-52 flex-col gap-2 overflow-y-auto">
				<li
					v-for="failure in report.failures"
					:key="`${failure.name}:${failure.reason}`"
					class="flex items-start justify-between gap-3 rounded-md border p-2"
				>
					<span class="min-w-0 truncate font-medium">{{ failure.name }}</span>
					<span class="shrink-0 text-right text-muted-foreground">
						{{ t(`imports.failureReasons.${reasonKeys[failure.reason]}`) }}
					</span>
				</li>
			</ul>
		</AlertDescription>
	</Alert>
</template>

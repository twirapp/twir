<script setup lang="ts">
import { InfoIcon, PlusIcon } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'

import Table from '@/components/table.vue'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import CommandButton from '@/features/commands/ui/command-button.vue'
import {
	useExpiringVipsTable,
} from '@/features/expiring-vips/composables/use-expiring-vips-table.ts'
import PageLayout from '@/layout/page-layout.vue'

const expiringVipsTable = useExpiringVipsTable()
const { t } = useI18n()
</script>

<template>
	<PageLayout>
		<template #title>
			{{ t('expiringVips.title') }}
		</template>

		<template #title-footer>
			<div class="flex flex-col">
				<span>
					{{ t('expiringVips.description.line1') }}
				</span>
				<span>
					{{ t('expiringVips.description.line2', { module: 'VIPS' }) }}
				</span>
			</div>
		</template>

		<template #action>
			<Button disabled>
				<PlusIcon class="size-4 mr-2" />
				<div class="flex flex-col">
					<span>
						{{ t('sharedButtons.create') }}
					</span>
					<span class="text-xs">
						{{ t('expiringVips.comingSoon') }}
					</span>
				</div>
			</Button>
		</template>

		<template #content>
			<div class="flex flex-col gap-4">
				<Card>
					<CardHeader>
						<CardTitle>{{ t('sidebar.commands.label') }}</CardTitle>
					</CardHeader>
					<CardContent>
						<div class="flex flex-wrap gap-4">
							<CommandButton name="vips add" :title="t('expiringVips.commands.add')" />
							<CommandButton name="vips remove" :title="t('expiringVips.commands.remove')" />
							<CommandButton name="vips list" :title="t('expiringVips.commands.list')" />
							<CommandButton name="vips setexpire" :title="t('expiringVips.commands.setExpire')" />
						</div>
					</CardContent>
				</Card>

				<Alert>
					<InfoIcon class="h-4 w-4" />
					<AlertTitle>{{ t('expiringVips.alert.title') }}</AlertTitle>
					<AlertDescription>
						{{ t('expiringVips.alert.description') }}
					</AlertDescription>
				</Alert>

				<Table :table="expiringVipsTable.table" :is-loading="expiringVipsTable.isLoading.value" />
			</div>
		</template>
	</PageLayout>
</template>

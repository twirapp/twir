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
			Scheduled vips
		</template>

		<template #title-footer>
			<div class="flex flex-col">
				<span>
					With this feature you can creat vips for your channel, which will be removed after selected time.
				</span>
				<span>
					You can use commands from built-in <b class="font-bold">VIPS</b> module for that.
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
						Soon, but already available through commands
					</span>
				</div>
			</Button>
		</template>

		<template #content>
			<div class="flex flex-col gap-4">
				<Card>
					<CardHeader>
						<CardTitle>Commands</CardTitle>
					</CardHeader>
					<CardContent>
						<div class="flex flex-wrap gap-4">
							<CommandButton name="vips add" title="Add vip" />
							<CommandButton name="vips remove" title="Remove vip" />
							<CommandButton name="vips list" title="List vips" />
							<CommandButton name="vips setexpire" title="Extend expiring time for exited timed vip, or add expiration for some vip user" />
						</div>
					</CardContent>
				</Card>

				<Alert>
					<InfoIcon class="h-4 w-4" />
					<AlertTitle>Heads up!</AlertTitle>
					<AlertDescription>
						It's a list of expiring vips, not list of all channel vips.
					</AlertDescription>
				</Alert>

				<Table :table="expiringVipsTable.table" :is-loading="expiringVipsTable.isLoading.value" />
			</div>
		</template>
	</PageLayout>
</template>

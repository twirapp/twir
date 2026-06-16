<script setup lang="ts">
import { computed, ref } from 'vue'
import { useCommandsApi } from '~~/layers/dashboard/api/commands/commands.js'
import Table from '~~/layers/dashboard/components/table.vue'
import CommandsList from '~~/layers/dashboard/features/commands/ui/list.vue'
import { useExpiringVipsTable } from '~~/layers/dashboard/features/expiring-vips/composables/use-expiring-vips-table.js'
import ExpiringVipsCreateDialog from '~~/layers/dashboard/features/expiring-vips/ui/expiring-vips-create-dialog.vue'
import PageLayout from '~~/layers/dashboard/layout/page-layout.vue'

import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'

const expiringVipsTable = useExpiringVipsTable()
const { t } = useI18n()

const commandsApi = useCommandsApi()
const { data: commands } = commandsApi.useQueryCommands()

const expiringVipsCommands = computed(() => {
	return commands.value?.commands?.filter((c) => c.module === 'VIPS')
})

const showCreateDialog = ref(false)
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
			<Button @click="showCreateDialog = true">
				<Icon
					name="lucide:plus"
					class="mr-2 size-4"
				/>
				{{ t('sharedButtons.create') }}
			</Button>
		</template>

		<template #content>
			<div class="flex flex-col gap-4">
				<Card>
					<CardHeader>
						<CardTitle>{{ t('sidebar.commands.label') }}</CardTitle>
					</CardHeader>
					<CardContent>
						<CommandsList
							v-if="expiringVipsCommands"
							:commands="expiringVipsCommands"
							show-background
						/>
					</CardContent>
				</Card>

				<Alert>
					<Icon
						name="lucide:info"
						class="h-4 w-4"
					/>
					<AlertTitle>{{ t('expiringVips.alert.title') }}</AlertTitle>
					<AlertDescription>
						{{ t('expiringVips.alert.description') }}
					</AlertDescription>
				</Alert>

				<Table
					:table="expiringVipsTable.table"
					:is-loading="expiringVipsTable.isLoading.value"
				/>
			</div>
			<ExpiringVipsCreateDialog v-model:open="showCreateDialog" />
		</template>
	</PageLayout>
</template>

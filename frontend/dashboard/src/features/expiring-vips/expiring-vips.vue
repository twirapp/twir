<script setup lang="ts">
import { InfoIcon, PlusIcon } from 'lucide-vue-next'
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import { useCommandsApi } from '@/api/commands/commands.ts'
import Table from '@/components/table.vue'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import CommandsList from '@/features/commands/ui/list.vue'
import { useExpiringVipsTable } from '@/features/expiring-vips/composables/use-expiring-vips-table.ts'
import ExpiringVipsCreateDialog from '@/features/expiring-vips/ui/expiring-vips-create-dialog.vue'
import PageLayout from '@/layout/page-layout.vue'

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
				<PlusIcon class="size-4 mr-2" />
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
					<InfoIcon class="h-4 w-4" />
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

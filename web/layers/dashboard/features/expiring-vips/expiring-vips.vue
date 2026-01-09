<script setup lang="ts">
import { InfoIcon, PlusIcon } from 'lucide-vue-next'
import { computed } from 'vue'


import { useCommandsApi } from '#layers/dashboard/api/commands/commands.ts'
import Table from '#layers/dashboard/components/table.vue'
import CommandsList from '~/features/commands/ui/list.vue'
import {
	useExpiringVipsTable,
} from '~/features/expiring-vips/composables/use-expiring-vips-table.ts'
import PageLayout from '~/layout/page-layout.vue'

const expiringVipsTable = useExpiringVipsTable()
const { t } = useI18n()

const commandsApi = useCommandsApi()
const { data: commands } = commandsApi.useQueryCommands()

const expiringVipsCommands = computed(() => {
	return commands.value?.commands?.filter(c => c.module === 'VIPS')
})
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
			<UiButton disabled>
				<PlusIcon class="size-4 mr-2" />
				<div class="flex flex-col">
					<span>
						{{ t('sharedButtons.create') }}
					</span>
					<span class="text-xs">
						{{ t('expiringVips.comingSoon') }}
					</span>
				</div>
			</UiButton>
		</template>

		<template #content>
			<div class="flex flex-col gap-4">
				<UiCard>
					<UiCardHeader>
						<UiCardTitle>{{ t('sidebar.commands.label') }}</UiCardTitle>
					</UiCardHeader>
					<UiCardContent>
						<CommandsList
							v-if="expiringVipsCommands"
							:commands="expiringVipsCommands"
							show-background
						/>
					</UiCardContent>
				</UiCard>

				<UiAlert>
					<InfoIcon class="h-4 w-4" />
					<UiAlertTitle>{{ t('expiringVips.alert.title') }}</UiAlertTitle>
					<UiAlertDescription>
						{{ t('expiringVips.alert.description') }}
					</UiAlertDescription>
				</UiAlert>

				<Table :table="expiringVipsTable.table" :is-loading="expiringVipsTable.isLoading.value" />
			</div>
		</template>
	</PageLayout>
</template>

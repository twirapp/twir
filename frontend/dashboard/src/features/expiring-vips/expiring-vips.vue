<script setup lang="ts">
import { PlusIcon } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'

import Table from '@/components/table.vue'
import { Button } from '@/components/ui/button'
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
			<Table :table="expiringVipsTable.table" :is-loading="expiringVipsTable.isLoading.value" />
		</template>
	</PageLayout>
</template>

<script setup lang="ts">
import { PlusIcon } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'

import VariablesDialog from './variables-dialog.vue'

import { Button } from '@/components/ui/button'

const { t } = useI18n()
</script>

<template>
	<div class="flex gap-2">
		<VariablesDialog>
			<template #dialog-trigger>
				<Button>
					<PlusIcon class="size-4 mr-2" />
					{{ t('sharedButtons.create') }}
				</Button>
			</template>
		</VariablesDialog>
	</div>
</template>

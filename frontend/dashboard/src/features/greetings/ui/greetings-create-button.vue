<script setup lang="ts">
import { PlusIcon } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'

import GreetingsDialog from './greetings-dialog.vue'

import { Button } from '@/components/ui/button'

const { t } = useI18n()
</script>

<template>
	<div class="flex gap-2">
		<GreetingsDialog>
			<template #dialog-trigger>
				<Button>
					<PlusIcon class="size-4 mr-2" />
					{{ t('greetings.create') }}
				</Button>
			</template>
		</GreetingsDialog>
	</div>
</template>

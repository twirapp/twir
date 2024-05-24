<script setup lang="ts">
import { PlusIcon } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'

import KeywordsDialog from './keywords-dialog.vue'

import { Button } from '@/components/ui/button'

const { t } = useI18n()
</script>

<template>
	<div class="flex gap-2">
		<KeywordsDialog>
			<template #dialog-trigger>
				<Button>
					<PlusIcon class="size-4 mr-2" />
					{{ t('keywords.create') }}
				</Button>
			</template>
		</KeywordsDialog>
	</div>
</template>

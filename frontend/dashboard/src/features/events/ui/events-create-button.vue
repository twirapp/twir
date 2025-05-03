<script setup lang="ts">
import { PlusIcon } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

import { Button } from '@/components/ui/button'

const { t } = useI18n()
const router = useRouter()

function createEvent() {
	router.push('/dashboard/events/new')
}
</script>

<template>
	<div class="flex gap-2">
		<Button @click="createEvent">
			<PlusIcon class="size-4 mr-2" />
			{{ t('events.create') }}
		</Button>
	</div>
</template>

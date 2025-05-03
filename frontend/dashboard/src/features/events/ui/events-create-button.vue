<script setup lang="ts">
import { PlusIcon } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

import { useEventsApi } from '@/api/events.ts'
import { Button } from '@/components/ui/button'

const { t } = useI18n()
const router = useRouter()

function createEvent() {
	router.push('/dashboard/events/new')
}

const eventApi = useEventsApi()
const { data: events } = eventApi.useQueryEvents()
</script>

<template>
	<div class="flex gap-2">
		<Button type="button" :disabled="events?.events.length >= 50" @click="createEvent">
			<PlusIcon class="size-4 mr-2" />
			{{ t('sharedTexts.create') }}
		</Button>
	</div>
</template>

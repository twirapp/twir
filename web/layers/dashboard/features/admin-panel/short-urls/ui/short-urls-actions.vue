<script setup lang="ts">
import { TrashIcon } from 'lucide-vue-next'
import { ref } from 'vue'

import type { AdminShortUrl } from '@/api/admin/short-urls.ts'

import ActionConfirm from '@/components/ui/action-confirm.vue'
import { Button } from '@/components/ui/button'
import {
	useAdminShortUrlsApi,
} from '@/features/admin-panel/short-urls/composables/use-admin-short-urls-api.ts'

defineProps<{
	item: AdminShortUrl
}>()

const api = useAdminShortUrlsApi()
const showDelete = ref(false)
</script>

<template>
	<div class="flex gap-2">
		<Button
			variant="destructive"
			size="icon"
			@click="showDelete = true"
		>
			<TrashIcon class="size-4" />
		</Button>
	</div>

	<ActionConfirm
		v-model:open="showDelete"
		@confirm="api.removeShortUrl(item.id)"
	/>
</template>

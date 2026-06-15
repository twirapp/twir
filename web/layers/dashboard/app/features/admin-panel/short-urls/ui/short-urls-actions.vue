<script setup lang="ts">
import { ref } from 'vue'

import type { AdminShortUrl } from '~~/layers/dashboard/app/api/admin/short-urls.js'

import ActionConfirm from '@/components/ui/action-confirm'
import { Button } from '@/components/ui/button'
import {
	useAdminShortUrlsApi,
} from '~~/layers/dashboard/app/features/admin-panel/short-urls/composables/use-admin-short-urls-api.js'

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
			<Icon name="lucide:trash" class="size-4" />
		</Button>
	</div>

	<ActionConfirm
		v-model:open="showDelete"
		@confirm="api.removeShortUrl(item.id)"
	/>
</template>

<script setup lang="ts">
import { TrashIcon } from 'lucide-vue-next'
import { ref } from 'vue'

import type { AdminShortUrl } from '#layers/dashboard/api/admin/short-urls.ts'


import {
	useAdminShortUrlsApi,
} from '~/features/admin-panel/short-urls/composables/use-admin-short-urls-api.ts'

defineProps<{
	item: AdminShortUrl
}>()

const api = useAdminShortUrlsApi()
const showDelete = ref(false)
</script>

<template>
	<div class="flex gap-2">
		<UiButton
			variant="destructive"
			size="icon"
			@click="showDelete = true"
		>
			<TrashIcon class="size-4" />
		</UiButton>
	</div>

	<UiActionConfirm
		v-model:open="showDelete"
		@confirm="api.removeShortUrl(item.id)"
	/>
</template>

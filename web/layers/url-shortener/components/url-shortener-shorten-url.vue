<script setup lang="ts">
import type { LinkOutputDto } from '@twir/api/openapi'
import { toast } from 'vue-sonner'

const props = defineProps<{ url: LinkOutputDto }>()

const clipboardApi = useClipboard()

function copyUrl() {
	clipboardApi.copy(props.url.short_url)

	toast('Copied', {
		description: 'Shortened url copied to clipboard',
		duration: 2500,
	})
}
</script>

<template>
	<div class="block p-4 border-border border-2 rounded-md bg-zinc-800 w-full underline relative">
		<UiButton size="sm" variant="default" class="absolute -right-2 -top-4" @click="copyUrl">
			<Icon name="lucide:copy" class="size-4" />
		</UiButton>
		<a :href="url.short_url" target="_blank">
			{{ url.short_url }}
		</a>
	</div>
</template>

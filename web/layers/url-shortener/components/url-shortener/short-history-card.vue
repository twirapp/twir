<script setup lang="ts">
import type { LinkOutputDto } from '@twir/api/openapi'
import { toast } from 'vue-sonner'
import { useMetaExtractor } from '../../composables/use-meta-extractor'
import UrlShortenerQrCodeDialog from './qr-code-dialog.vue'

const props = defineProps<{
	url: LinkOutputDto
}>()

const clipboard = useClipboard()
const { extractMetaFromUrl, loading: metaLoading } = useMetaExtractor()
const hasLoaded = ref(false)
const metaData = ref<any>(null)
const showQrDialog = ref(false)

function removeProtocol(value: string) {
	return value.replace(/^https?:\/\//, '')
}

const displayShortUrl = computed(() => removeProtocol(props.url.short_url))
const displayUrl = computed(() => removeProtocol(props.url.url))
const createdAt = computed(() => {
	const date = new Date(props.url.created_at)
	if (Number.isNaN(date.getTime())) {
		return ''
	}

	return date.toLocaleString(import.meta.client ? navigator.language : 'en-US', {
		month: 'short',
		day: '2-digit',
		year: 'numeric',
		hour: '2-digit',
		minute: '2-digit',
	})
})

function copyShortUrl() {
	clipboard.copy(props.url.short_url)
	toast.success('Copied', {
		description: 'Short link copied to clipboard.',
		duration: 2000,
	})
}

async function fetchMetadata() {
	if (metaLoading.value || hasLoaded.value) return

	if (props.url.url) {
		hasLoaded.value = true

		try {
			const meta = await extractMetaFromUrl(props.url.url)
			metaData.value = meta
		} catch (err) {
			console.error('Failed to fetch metadata:', err)
			hasLoaded.value = false
		}
	}
}

watch(
	() => props.url.url,
	() => {
		fetchMetadata()
	},
	{ immediate: true }
)
</script>

<template>
	<div
		class="flex flex-col gap-2 rounded-2xl border border-[hsl(240,11%,18%)] bg-[hsl(240,11%,9%)] p-3 shadow-[0px_0px_30px_hsl(240,11%,6%)]"
	>
		<div class="flex items-start justify-between gap-3">
			<div class="flex items-start gap-3 min-w-0">
				<div
					class="flex-none size-9 rounded-full border border-[hsl(240,11%,25%)] bg-[hsl(240,11%,15%)] flex items-center justify-center"
				>
					<Icon
						v-if="!hasLoaded || !metaData || !metaData.favicon"
						name="lucide:link"
						class="h-4 w-4"
					/>
					<img v-else :src="metaData.favicon" class="h-4 w-4" />
				</div>
				<div class="min-w-0">
					<a
						:href="props.url.short_url"
						target="_blank"
						class="block font-semibold text-[hsl(240,11%,90%)] hover:text-white transition-colors truncate"
					>
						{{ displayShortUrl }}
					</a>
					<div class="flex items-center gap-1 text-xs text-[hsl(240,11%,60%)]">
						<Icon name="lucide:corner-down-right" class="h-4 w-4 shrink-0" />
						<span class="truncate">{{ displayUrl }}</span>
					</div>
				</div>
			</div>
			<div class="flex items-center gap-2">
				<button
					class="flex items-center justify-center rounded-lg border border-[hsl(240,11%,25%)] bg-[hsl(240,11%,15%)] p-1.5 text-[hsl(240,11%,80%)] hover:border-[hsl(240,11%,40%)] hover:bg-[hsl(240,11%,25%)] transition-colors"
					title="Copy short URL"
					@click="copyShortUrl"
				>
					<Icon name="lucide:copy" class="h-3.5 w-3.5" />
				</button>
				<button
					class="flex items-center justify-center rounded-lg border border-[hsl(240,11%,25%)] bg-[hsl(240,11%,15%)] p-1.5 text-[hsl(240,11%,80%)] hover:border-[hsl(240,11%,40%)] hover:bg-[hsl(240,11%,25%)] transition-colors"
					title="Show QR code"
					@click="showQrDialog = true"
				>
					<Icon name="lucide:qr-code" class="h-3.5 w-3.5" />
				</button>
				<a
					:href="props.url.short_url"
					target="_blank"
					class="flex items-center justify-center rounded-lg border border-[hsl(240,11%,25%)] bg-[hsl(240,11%,15%)] p-1.5 text-[hsl(240,11%,80%)] hover:border-[hsl(240,11%,40%)] hover:bg-[hsl(240,11%,25%)] transition-colors"
					title="Open short URL"
				>
					<Icon name="lucide:external-link" class="h-3.5 w-3.5" />
				</a>
			</div>
		</div>
		<p v-if="createdAt" class="text-xs text-[hsl(240,11%,55%)]">
			Created {{ createdAt }}
		</p>

		<ClientOnly>
			<UrlShortenerQrCodeDialog v-model:open="showQrDialog" :short-url="props.url.short_url" />
		</ClientOnly>
	</div>
</template>

<script setup lang="ts">
import { toast } from 'vue-sonner'
import { TwirLogo } from '@twir/brand'

import { useMetaExtractor } from '../../composables/use-meta-extractor'
import { useQRCode } from '../../composables/use-qr-code'

import type { LinkOutputDto } from '@twir/api/openapi'

import PixelBlast from '~/components/ui/bits/backgrounds/PixelBlast/pixel-blast.vue'
import { ColorPicker } from '~/components/ui/color-picker'

const props = defineProps<{ url: LinkOutputDto }>()

const clipboardApi = useClipboard()

const { extractMetaFromUrl, error, loading } = useMetaExtractor()

const hasLoaded = ref(false)
const metaData = ref<any>(null)

function removeProtocol(url: string) {
	return url.replace(/^https?:\/\//, '')
}

const displayShortUrl = computed(() => removeProtocol(props.url.short_url))
const displayUrl = computed(() => removeProtocol(props.url.url))

function copyUrl() {
	clipboardApi.copy(props.url.short_url)

	toast.success('Copied', {
		description: 'Shortened url copied to clipboard',
		duration: 2500,
	})
}

function formatViews(views: number) {
	const intl = new Intl.NumberFormat(import.meta.client ? navigator.language : 'en-US', {
		notation: 'compact',
		maximumFractionDigits: 1,
	})
	return intl.format(views)
}

const { generateQRCode, downloadQR } = useQRCode()
const showQRModal = ref(false)
const qrCodeUrl = ref<string>('')
const qrSettings = reactive({
	showLogo: true,
	color: '#ffffff',
	backgroundColor: '#00000000',
})

const logoRef = ref<HTMLElement | null>(null)

async function openQRCode() {
	if (!props.url?.short_url) return
	showQRModal.value = true
	await generateQR()
}

async function generateQR() {
	if (!props.url?.short_url) return

	qrCodeUrl.value = await generateQRCode({
		url: props.url.short_url,
		color: qrSettings.color,
		backgroundColor: qrSettings.backgroundColor,
		logoComponent: qrSettings.showLogo ? logoRef.value : null,
		logoSize: 70,
		size: 400,
	})
}

async function handleDownloadQR() {
	if (!qrCodeUrl.value) return
	await downloadQR(qrCodeUrl.value, `twir-qrcode.png`)
}

watch([() => qrSettings.showLogo, () => qrSettings.color], async () => {
	if (showQRModal.value) {
		await generateQR()
	}
})

async function fetchMetadata() {
	if (loading.value || hasLoaded.value) return

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
	<div class="flex w-full justify-center">
		<div
			class="flex justify-between items-center bg-[hsl(240,11%,9%)] border border-[hsl(240,11%,18%)] h-fit w-full max-w-xl rounded-2xl p-3 shadow-[0px_0px_30px_hsl(240,11%,6%)]"
		>
			<div class="flex justify-between items-center gap-3">
				<div
					class="flex size-fit p-3 rounded-full font-semibold border border-[hsl(240,11%,25%)] bg-[hsl(240,11%,20%)]"
				>
					<Icon v-if="!hasLoaded || !metaData" name="lucide:link" class="size-4" />
					<img v-else :src="metaData.favicon" class="size-4" />
				</div>
				<div class="flex flex-col gap-1">
					<div class="flex items-center">
						<a
							:href="props.url.short_url"
							target="_blank"
							class="font-bold hover:hover:text-[hsl(240,11%,85%)] transition-colors max-w-60 truncate"
						>
							{{ displayShortUrl }}
						</a>
					</div>
					<span class="flex gap-1">
						<Icon name="lucide:corner-down-right" class="size-4 text-[hsl(240,11%,50%)]" />
						<a
							:href="props.url.url"
							target="_blank"
							class="text-sm font-medium text-[hsl(240,11%,50%)] hover:text-[hsl(240,11%,65%)] transition-colors max-w-60 truncate"
						>
							{{ displayUrl }}
						</a>
					</span>
				</div>
			</div>
			<div class="flex items-center gap-2">
				<span
					class="cursor-pointer flex gap-1.5 items-center text-sm size-fit p-2.5 rounded-lg font-semibold border border-[hsl(240,11%,25%)] hover:border-[hsl(240,11%,40%)] bg-[hsl(240,11%,20%)] hover:bg-[hsl(240,11%,30%)] transition-colors"
				>
					<Icon name="lucide:mouse-pointer-click" class="size-4" />
					{{ formatViews(props.url.views) }} views
				</span>
				<UiButton
					variant="outline"
					class="size-fit p-3 rounded-lg font-semibold border border-[hsl(240,11%,25%)] hover:border-[hsl(240,11%,40%)] bg-[hsl(240,11%,20%)] hover:bg-[hsl(240,11%,30%)]"
					@click="copyUrl"
				>
					<Icon name="lucide:copy" class="size-4" />
				</UiButton>
				<UiButton
					variant="outline"
					class="size-fit p-3 rounded-lg font-semibold border border-[hsl(240,11%,25%)] hover:border-[hsl(240,11%,40%)] bg-[hsl(240,11%,20%)] hover:bg-[hsl(240,11%,30%)]"
					@click="openQRCode"
				>
					<Icon name="lucide:qr-code" class="size-4" />
				</UiButton>
			</div>
		</div>

		<div ref="logoRef" class="fixed -left-[9999px]">
			<TwirLogo />
		</div>
		<ClientOnly>
			<UiDialog v-model:open="showQRModal">
				<UiDialogContent class="max-w-lg">
					<UiDialogHeader>
						<UiDialogTitle>QR Code Settings</UiDialogTitle>
					</UiDialogHeader>

					<div class="flex flex-col gap-5">
						<div class="flex flex-col gap-2">
							<div class="flex items-center justify-between">
								<span class="text-sm font-semibold">Preview</span>
								<div class="flex gap-1 items-center">
									<UiButton class="size-7" variant="ghost" size="custom" @click="handleDownloadQR">
										<Icon name="lucide:download" class="size-4" />
									</UiButton>
								</div>
							</div>
							<PixelBlast
								:pixel-size="5"
								:pattern-scale="3.5"
								:pattern-density="1.75"
								:pixel-jiter="2"
								:speed="0.25"
								:edge-fade="0.32"
								:enable-ripples="false"
								class="flex justify-center p-6 rounded-lg border"
							>
								<div class="flex border bg-black/10 backdrop-blur rounded-md">
									<img
										v-if="qrCodeUrl"
										:src="qrCodeUrl"
										alt="QR Code"
										class="w-full max-w-[180px] rounded-lg p-3"
									/>
								</div>
							</PixelBlast>
						</div>

						<div class="flex flex-col gap-4">
							<div class="flex flex-row items-center justify-between">
								<label class="text-sm font-medium">Show Twir Logo</label>
								<UiSwitchToggle v-model="qrSettings.showLogo" />
							</div>

							<div class="flex flex-row items-center justify-between">
								<label class="text-sm font-medium">QR Code Color</label>
								<ColorPicker v-model="qrSettings.color" class="size-6" />
							</div>
						</div>
					</div>
				</UiDialogContent>
			</UiDialog>
		</ClientOnly>
	</div>
</template>

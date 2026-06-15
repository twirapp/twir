<script setup lang="ts">
import { useQRCode } from '../../composables/use-qr-code'
import PixelBlast from '~/components/ui/bits/backgrounds/PixelBlast/pixel-blast.vue'
import { TwirLogo } from '@twir/brand'
import { render } from 'vue'

const props = defineProps<{
	shortUrl: string
}>()

const { generateQRCode, downloadQR } = useQRCode()
const qrCodeUrl = ref<string>('')
const qrSettings = reactive({
	showLogo: true,
	color: '#ffffff',
	backgroundColor: '#00000000',
})

const logoContainer = ref<HTMLDivElement | null>(null)

onMounted(async () => {
	const container = document.createElement('div')
	const logoVnode = h(TwirLogo)
	render(logoVnode, container)
	logoContainer.value = container

	await generateQR()
})

async function generateQR() {
	if (!props.shortUrl) return

	qrCodeUrl.value = await generateQRCode({
		url: props.shortUrl,
		color: qrSettings.color,
		backgroundColor: qrSettings.backgroundColor,
		logoComponent: qrSettings.showLogo && logoContainer.value ? logoContainer.value : null,
		logoSize: 70,
		size: 400,
	})
}

async function handleDownloadQR() {
	if (!qrCodeUrl.value) return
	await downloadQR(qrCodeUrl.value, `twir-qrcode.png`)
}

watch([() => qrSettings.showLogo, () => qrSettings.color], async () => {
	await generateQR()
})
</script>

<template>
	<UiDialog>
		<UiDialogTrigger as-child>
			<slot name="trigger" />
		</UiDialogTrigger>
		<UiDialogContent class="max-w-lg">
			<UiDialogHeader>
				<UiDialogTitle>QR Code Settings</UiDialogTitle>
			</UiDialogHeader>

			<div class="flex flex-col gap-5">
				<div class="flex flex-col gap-2">
					<div class="flex items-center justify-between">
						<span class="text-sm font-semibold">Preview</span>
						<div class="flex gap-1 items-center">
							<UiButton
								class="size-7"
								variant="ghost"
								size="default"
								@click="handleDownloadQR"
							>
								<Icon
									name="lucide:download"
									class="size-4"
								/>
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
						<UiColorPicker
							v-model="qrSettings.color"
							class="size-6"
						/>
					</div>
				</div>
			</div>
		</UiDialogContent>
	</UiDialog>
</template>

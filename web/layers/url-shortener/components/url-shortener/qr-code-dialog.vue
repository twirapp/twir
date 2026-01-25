<script setup lang="ts">
import Dialog from '@/components/ui/dialog/Dialog.vue'
import DialogContent from '@/components/ui/dialog/DialogContent.vue'
import DialogDescription from '@/components/ui/dialog/DialogDescription.vue'
import DialogHeader from '@/components/ui/dialog/DialogHeader.vue'
import DialogTitle from '@/components/ui/dialog/DialogTitle.vue'
import Button from '@/components/ui/button/Button.vue'

import { useQRCode } from '../../composables/use-qr-code'

const props = defineProps<{
	open: boolean
	shortUrl: string
}>()

const emit = defineEmits<{
	(e: 'update:open', value: boolean): void
}>()

const { generateQRCode, downloadQR } = useQRCode()
const qrCode = ref<string | null>(null)
const isLoading = ref(false)
const errorMessage = ref<string | null>(null)

const downloadFileName = computed(() => {
	if (!props.shortUrl) return 'twir-qrcode.png'

	const cleaned = props.shortUrl
		.replace(/^https?:\/\//, '')
		.replace(/[^a-zA-Z0-9-_]+/g, '-')
		.replace(/-+/g, '-')
		.replace(/^-|-$/g, '')

	if (!cleaned) return 'twir-qrcode.png'

	return `twir-${cleaned}.png`
})

watch(
	() => [props.open, props.shortUrl],
	async ([isOpen]) => {
		if (!isOpen || !props.shortUrl) return
		if (!import.meta.client) return

		isLoading.value = true
		errorMessage.value = null

		try {
			qrCode.value = await generateQRCode({ url: props.shortUrl })
		} catch (error) {
			console.error('Failed to generate QR code:', error)
			qrCode.value = null
			errorMessage.value = 'Failed to generate QR code'
		} finally {
			isLoading.value = false
		}
	}
)

function handleDownload() {
	if (!qrCode.value) return
	downloadQR(qrCode.value, downloadFileName.value)
}

function closeDialog() {
	emit('update:open', false)
}
</script>

<template>
	<Dialog :open="open" @update:open="closeDialog">
		<DialogContent class="max-w-sm">
			<DialogHeader>
				<DialogTitle>QR Code</DialogTitle>
				<DialogDescription>
					Scan to open
					<span class="break-all">{{ shortUrl }}</span>
				</DialogDescription>
			</DialogHeader>

			<div class="mt-4 flex items-center justify-center">
				<div
					class="flex h-48 w-48 items-center justify-center rounded-2xl border border-[hsl(240,11%,18%)] bg-[hsl(240,11%,12%)] p-3 shadow-[0px_0px_20px_hsl(240,11%,6%)]"
				>
					<Icon v-if="isLoading" name="lucide:loader-2" class="h-8 w-8 animate-spin" />
					<Icon v-else-if="errorMessage" name="lucide:alert-triangle" class="h-8 w-8 text-red-500" />
					<img v-else-if="qrCode" :src="qrCode" alt="QR code" class="h-full w-full" />
					<Icon v-else name="lucide:qr-code" class="h-8 w-8 text-[hsl(240,11%,50%)]" />
				</div>
			</div>

			<p v-if="errorMessage" class="mt-3 text-sm text-red-400">
				{{ errorMessage }}
			</p>

			<div class="mt-4 flex items-center justify-between gap-2">
				<Button variant="outline" size="sm" :disabled="!qrCode || isLoading" @click="handleDownload">
					<Icon name="lucide:download" class="mr-2 h-4 w-4" />
					Download
				</Button>
				<Button variant="outline" size="sm" @click="closeDialog">
					Close
				</Button>
			</div>
		</DialogContent>
	</Dialog>
</template>

<script setup lang="ts">
import type { UploadedFileWithDeleteLinkDto } from '@twir/api/openapi'

const { t } = useI18n()
const uploader = useUploader()

// Mirrors backend constraints (apps/api-gql uploader service): 25 MiB, 6 image types.
const ACCEPTED_MIME_TYPES = [
	'image/png',
	'image/jpeg',
	'image/gif',
	'image/webp',
	'image/avif',
	'image/bmp',
]
const acceptAttribute = ACCEPTED_MIME_TYPES.join(',')
const maxFileSizeBytes = 26_214_400

const isDragging = ref(false)
const currentUpload = ref<UploadedFileWithDeleteLinkDto | null>(null)
const currentError = ref<string | null>(null)

const maxFileSizeLabel = computed(() => formatBytes(maxFileSizeBytes, 0))

async function handleFile(file: File | null | undefined) {
	if (!file || uploader.isUploading) return

	currentError.value = null
	currentUpload.value = null

	if (!ACCEPTED_MIME_TYPES.includes(file.type)) {
		currentError.value = t('uploader.errors.unsupportedType')
		return
	}

	if (file.size > maxFileSizeBytes) {
		currentError.value = t('uploader.errors.tooLarge', { maxSize: maxFileSizeLabel.value })
		return
	}

	const { data, error } = await uploader.uploadFile(file)
	if (error) {
		currentError.value = error
		return
	}

	currentUpload.value = data
}

function onFileInputChange(event: Event) {
	const input = event.target as HTMLInputElement
	handleFile(input.files?.[0])
	// Allow picking the same file again.
	input.value = ''
}

function isFileDrag(event: DragEvent) {
	return Array.from(event.dataTransfer?.types ?? []).includes('Files')
}

// Page-wide listeners. String-only useEventListener overload binds to window
// on the client and no-ops during SSR, with automatic cleanup on unmount.
useEventListener('dragover', (event: DragEvent) => {
	if (!isFileDrag(event)) return
	event.preventDefault()
	isDragging.value = true
})

useEventListener('dragleave', (event: DragEvent) => {
	if (!event.relatedTarget) {
		isDragging.value = false
	}
})

useEventListener('drop', (event: DragEvent) => {
	if (!isFileDrag(event)) return
	event.preventDefault()
	isDragging.value = false
	handleFile(event.dataTransfer?.files?.[0])
})

useEventListener('paste', (event: ClipboardEvent) => {
	const clipboardData = event.clipboardData
	if (!clipboardData) return

	let file = clipboardData.files?.[0]
	if (!file) {
		file =
			Array.from(clipboardData.items)
				.find((item) => item.kind === 'file')
				?.getAsFile() ?? undefined
	}

	if (!file) return
	event.preventDefault()
	handleFile(file)
})
</script>

<template>
	<div class="flex flex-col w-full max-w-xl gap-3">
		<div
			class="flex flex-col items-center justify-center w-full min-h-[240px] rounded-2xl border border-dashed p-8 text-center transition-colors"
			:class="
				isDragging
					? 'border-[hsl(240,11%,55%)] bg-[hsl(240,11%,15%)]'
					: 'border-[hsl(240,11%,18%)] bg-[hsl(240,11%,9%)] shadow-[0px_0px_30px_hsl(240,11%,6%)]'
			"
		>
			<label
				for="uploader-file-input"
				class="flex flex-col items-center justify-center gap-3 w-full h-full cursor-pointer text-[hsl(240,11%,90%)]"
			>
				<div
					class="flex rounded-xl border border-[hsl(240,11%,18%)] bg-[hsl(240,11%,10%)] p-4 text-[hsl(240,11%,90%)]"
				>
					<Icon
						v-if="uploader.isUploading"
						name="lucide:loader-2"
						class="w-6 h-6 animate-spin"
					/>
					<Icon v-else name="lucide:image-plus" class="w-6 h-6" />
				</div>
				<span class="font-semibold">
					{{
						uploader.isUploading
							? t('uploader.dropzone.uploading')
							: t('uploader.dropzone.hint')
					}}
				</span>
				<span class="text-xs text-[hsl(240,11%,55%)]">
					{{ t('uploader.dropzone.formats', { maxSize: maxFileSizeLabel }) }}
				</span>
			</label>
			<input
				id="uploader-file-input"
				type="file"
				:accept="acceptAttribute"
				class="sr-only"
				:disabled="uploader.isUploading"
				@change="onFileInputChange"
			/>
		</div>

		<p v-if="currentError" class="px-2 text-sm text-red-500">
			{{ currentError }}
		</p>

		<UploaderResultCard v-if="currentUpload" :upload="currentUpload" />

		<Teleport to="body">
			<div
				v-if="isDragging"
				class="fixed inset-0 z-50 flex items-center justify-center bg-[hsl(240,11%,6%)]/80 pointer-events-none"
			>
				<span class="text-xl font-semibold text-[hsl(240,11%,90%)]">{{ t('uploader.dropzone.dropHere') }}</span>
			</div>
		</Teleport>
	</div>
</template>

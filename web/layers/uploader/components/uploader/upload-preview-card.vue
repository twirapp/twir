<script setup lang="ts">
import { toast } from 'vue-sonner'

import type { UploadedFileOutputDto } from '@twir/api/openapi'

import ActionConfirm from '@/components/ui/action-confirm'

const props = defineProps<{
	upload: UploadedFileOutputDto
}>()

const emit = defineEmits<{
	(e: 'deleted'): void
}>()

const { t } = useI18n()
const uploader = useUploader()
const clipboard = useClipboard()

const showDeleteConfirm = ref(false)

const displayLink = computed(() => props.upload.link.replace(/^https?:\/\//, ''))

function formatDate(value: string) {
	const date = new Date(value)
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
}

const createdAt = computed(() => formatDate(props.upload.created_at))
const expiresAt = computed(() => formatDate(props.upload.expires_at))

function copyLink() {
	clipboard.copy(props.upload.link)
	toast.success(t('uploader.copied'), {
		description: t('uploader.copiedDescription'),
		duration: 2000,
	})
}

async function handleDelete() {
	const { error } = await uploader.deleteUpload(props.upload.id)

	if (error) {
		toast.error(t('uploader.errors.title'), {
			description: error,
		})
		return
	}

	toast.success(t('uploader.delete.success'))
	emit('deleted')
}
</script>

<template>
	<div
		class="flex flex-col overflow-hidden rounded-2xl border border-[hsl(240,11%,18%)] bg-[hsl(240,11%,9%)] shadow-[0px_0px_30px_hsl(240,11%,6%)]"
	>
		<a :href="upload.link" target="_blank" class="block bg-[hsl(240,11%,12%)]">
			<img
				:src="upload.link"
				:alt="upload.name ?? displayLink"
				class="h-40 w-full object-contain"
				loading="lazy"
			/>
		</a>
		<div class="flex flex-col gap-2 p-3">
			<a
				:href="upload.link"
				target="_blank"
				class="font-semibold text-[hsl(240,11%,90%)] hover:text-white transition-colors truncate"
			>
				{{ displayLink }}
			</a>
			<div class="flex items-center gap-1 text-xs text-[hsl(240,11%,60%)]">
				<span class="truncate">{{ upload.name ?? $t('uploader.untitled') }}</span>
				<span aria-hidden="true">·</span>
				<span class="flex-none">{{ formatBytes(upload.size) }}</span>
			</div>
			<div class="flex flex-col text-xs text-[hsl(240,11%,55%)]">
				<span v-if="createdAt">{{ $t('uploader.created', { date: createdAt }) }}</span>
				<span v-if="expiresAt">{{ $t('uploader.expires', { date: expiresAt }) }}</span>
			</div>
			<div class="flex items-center gap-2 pt-1">
				<button
					class="flex items-center justify-center rounded-lg border border-[hsl(240,11%,25%)] bg-[hsl(240,11%,15%)] p-1.5 text-[hsl(240,11%,80%)] hover:border-[hsl(240,11%,40%)] hover:bg-[hsl(240,11%,25%)] transition-colors"
					:title="$t('uploader.actions.copyLink')"
					@click="copyLink"
				>
					<Icon name="lucide:copy" class="h-3.5 w-3.5" />
				</button>
				<a
					:href="upload.link"
					target="_blank"
					class="flex items-center justify-center rounded-lg border border-[hsl(240,11%,25%)] bg-[hsl(240,11%,15%)] p-1.5 text-[hsl(240,11%,80%)] hover:border-[hsl(240,11%,40%)] hover:bg-[hsl(240,11%,25%)] transition-colors"
					:title="$t('uploader.actions.open')"
				>
					<Icon name="lucide:external-link" class="h-3.5 w-3.5" />
				</a>
				<button
					class="flex items-center justify-center rounded-lg border border-red-900/50 bg-red-950/30 p-1.5 text-red-400 hover:border-red-700 hover:bg-red-950/50 hover:text-red-300 transition-colors"
					:title="$t('uploader.actions.delete')"
					@click="showDeleteConfirm = true"
				>
					<Icon name="lucide:trash-2" class="h-3.5 w-3.5" />
				</button>
			</div>
		</div>

		<ClientOnly>
			<ActionConfirm
				v-model:open="showDeleteConfirm"
				:confirm-text="t('uploader.delete.confirmText')"
				@confirm="handleDelete"
			/>
		</ClientOnly>
	</div>
</template>

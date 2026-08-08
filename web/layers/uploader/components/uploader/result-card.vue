<script setup lang="ts">
import type { UploadedFileWithDeleteLinkDto } from '@twir/api/openapi'

const props = defineProps<{
	upload: UploadedFileWithDeleteLinkDto
}>()

const { t, locale } = useI18n()

const showDeleteLink = ref(false)

const displayLink = computed(() => props.upload.link.replace(/^https?:\/\//, ''))

const expiresAtAbsolute = computed(() => {
	const date = new Date(props.upload.expires_at)
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

const expiresIn = computed(() => formatTimeUntil(props.upload.expires_at, locale.value))
const expiresLabel = computed(() => {
	if (expiresIn.value) {
		return t('uploader.expiresIn', { time: expiresIn.value })
	}
	if (expiresAtAbsolute.value) {
		return t('uploader.expires', { date: expiresAtAbsolute.value })
	}
	return ''
})
</script>

<template>
	<div
		class="flex flex-col gap-3 w-full rounded-2xl border border-[hsl(240,11%,18%)] bg-[hsl(240,11%,9%)] p-3 shadow-[0px_0px_30px_hsl(240,11%,6%)]"
	>
		<div class="flex items-start gap-3">
			<a :href="upload.link" target="_blank" class="flex-none">
				<img
					:src="upload.link"
					:alt="upload.name ?? displayLink"
					class="size-16 rounded-xl border border-[hsl(240,11%,18%)] bg-[hsl(240,11%,15%)] object-cover"
				/>
			</a>
			<div class="min-w-0 flex-1">
				<p class="text-sm font-semibold text-[hsl(240,11%,90%)]">{{ t('uploader.result.title') }}</p>
				<a
					:href="upload.link"
					target="_blank"
					class="block text-sm text-[hsl(240,11%,65%)] hover:text-white transition-colors truncate"
				>
					{{ displayLink }}
				</a>
				<UploaderCopyInput :text="upload.link" class="mt-2" />
			</div>
		</div>

		<p
			v-if="expiresLabel"
			class="text-xs text-[hsl(240,11%,55%)]"
			:title="t('uploader.expires', { date: expiresAtAbsolute })"
		>
			{{ expiresLabel }}
		</p>

		<div class="border-t border-[hsl(240,11%,18%)] pt-2">
			<button
				type="button"
				class="flex items-center gap-1.5 text-xs text-[hsl(240,11%,55%)] hover:text-white transition-colors"
				@click="showDeleteLink = !showDeleteLink"
			>
				<Icon
					name="lucide:chevron-down"
					class="h-3.5 w-3.5 transition-transform"
					:class="{ 'rotate-180': showDeleteLink }"
				/>
				{{ t('uploader.result.showDeleteLink') }}
			</button>
			<div v-if="showDeleteLink" class="flex flex-col gap-1.5 mt-2">
				<p class="text-xs text-[hsl(240,11%,55%)]">
					{{ t('uploader.result.deleteLinkHint') }}
				</p>
				<UploaderCopyInput :text="upload.delete_link" />
			</div>
		</div>
	</div>
</template>

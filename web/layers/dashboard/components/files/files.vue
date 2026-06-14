<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { toast } from 'vue-sonner'
import { useFilesApi } from '~~/layers/dashboard/api/files'
import { convertBytesToSize } from '~~/layers/dashboard/helpers/convertBytesToSize.js'

import { Alert, AlertDescription } from '@/components/ui/alert'
import { Button, FileButton } from '@/components/ui/button'
import { Progress } from '@/components/ui/progress'

const props = withDefaults(
	defineProps<{
		tab?: string
		mode?: 'list' | 'picker'
	}>(),
	{
		tab: 'Audios',
		mode: 'list',
	}
)

defineEmits<{
	select: [id: string]
	delete: [id: string]
}>()

const { t } = useI18n()

const api = useFilesApi()
const deleter = api.useDelete()
const uploader = api.useUpload()

const { data: files } = api.useQuery()

const uploadedFilesSize = computed(() => {
	if (!files.value?.files) return 0

	return files.value?.files.reduce((acc, curr) => acc + Number(curr.size), 0)
})

interface Tab {
	name: string
	disabled?: boolean
	accept: string
}
const tabs: Array<Tab> = [
	{
		name: 'Audios',
		accept: 'audio/*',
	},
	{
		name: 'Images (soon)',
		disabled: true,
		accept: 'image/*',
	},
]
const activeTab = ref<Tab>(tabs.at(0)!)

onMounted(() => {
	const neededTab = tabs.find((t) => t.name === props.tab)
	if (!neededTab) return
	activeTab.value = neededTab
})

const audios = computed(
	() => files.value?.files.filter((f) => f.mimetype.startsWith('audio')) ?? []
)

async function upload(f: File) {
	if (!f.type.startsWith(activeTab.value.accept.split('*').at(0)!)) return

	try {
		await uploader.executeMutation({
			file: f,
		})
	} catch (error: any) {
		toast.error(error?.message)
	}
}

const maxFileSize = 104857600 // 100mb

const uploadedFilesSizeSlider = computed(() => {
	if (uploadedFilesSize.value === 0) {
		return 0
	}

	return Math.max(Math.abs(uploadedFilesSize.value / 1000000), 5)
})
</script>

<template>
	<div class="flex flex-col md:flex-row">
		<div
			class="bg-card border-border md:min-w-max-w-52 flex h-full w-full flex-col justify-between gap-y-8 border-r-2 border-b-2 p-2 md:max-w-52 md:border-b-0"
		>
			<div class="flex w-full flex-col gap-1">
				<Button
					class="flex w-full items-center justify-center gap-2"
					size="sm"
					variant="secondary"
				>
					<Icon
						name="lucide:music"
						class="size-4"
					/>
					Audios
				</Button>
				<Button
					class="flex w-full items-center justify-center gap-2"
					size="sm"
					variant="secondary"
					disabled
				>
					<Icon
						name="lucide:image"
						class="size-4"
					/>
					Images (soon)
				</Button>
			</div>
			<div class="flex flex-col items-center justify-center gap-2">
				<FileButton
					class="size-28"
					:accept="activeTab.accept"
					:disabled="uploader.fetching.value || uploadedFilesSize >= (1 << 20) * 100"
					:loading="uploader.fetching.value"
					@file-selected="
						(file) => {
							if (!file) return
							upload(file)
						}
					"
				/>
				<Progress v-model="uploadedFilesSizeSlider" />
				<span
					>{{ convertBytesToSize(uploadedFilesSize) }} / {{ convertBytesToSize(maxFileSize) }}</span
				>
			</div>
			<!-- <div>
				<NUpload
					multiple
					directory-dnd
					:max="1"
					:accept="activeTab.accept"
					:file-list="[]"
					:disabled="uploader.fetching.value || uploadedFilesSize >= ((1 << 20) * 100)"
					@change="(data) => {
						if (!data.fileList.length) return
						upload(data.file.file!)
					}"
				>
					<NUploadDragger>
						<div v-if="!uploader.fetching.value">
							<div class="mb-3">
								<NIcon size="30" :depth="3">
									<IconArchive />
								</NIcon>
							</div>
							<NText class="text-xs">
								{{ t('filePicker.innerText', { type: activeTab.name.toLowerCase() }) }}
							</NText>
						</div>
						<NSpin v-else />
					</NUploadDragger>
				</NUpload>

				<NText>
					{{
						t('filePicker.usedSpace', {
							used: convertBytesToSize(uploadedFilesSize),
							max: 100,
						})
					}}
				</NText>
			</div> -->
		</div>

		<div class="min-h-0 w-full overflow-auto p-4">
			<div v-if="activeTab.name === 'Audios'">
				<Alert v-if="!audios.length">
					<AlertDescription>
						{{ t('filePicker.emptyText', { type: 'audios' }) }}
					</AlertDescription>
				</Alert>

				<div
					v-else
					class="flex flex-col gap-4 overflow-y-auto"
				>
					<div
						v-for="audio of audios"
						:key="audio.id"
						class="bg-card border-border flex flex-col rounded-md border-2 p-2"
					>
						<h1 class="break-all">
							{{ audio.name }}
						</h1>
						<audio
							controls
							:src="api.computeFileUrl(audio.channelId, audio.id)"
							class="w-full"
						/>

						<div class="flex flex-wrap items-end justify-between">
							<span>{{ convertBytesToSize(audio.size) }}</span>

							<div class="flex gap-1">
								<Button
									variant="destructive"
									size="sm"
									@click="
										async () => {
											await deleter.executeMutation({ id: audio.id })
											$emit('delete', audio.id)
										}
									"
								>
									{{ t('sharedButtons.delete') }}
								</Button>
								<Button
									size="sm"
									variant="secondary"
									@click="$emit('select', audio.id)"
								>
									{{ t('sharedButtons.select') }}
								</Button>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

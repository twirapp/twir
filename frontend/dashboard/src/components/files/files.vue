<script setup lang="ts">
import { RpcError } from '@protobuf-ts/runtime-rpc'
import { ImageIcon, MusicIcon } from 'lucide-vue-next'
import {
	NAlert,
	useMessage,
} from 'naive-ui'
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import { useFilesApi } from '@/api/index.js'
import { Button, FileButton } from '@/components/ui/button'
import { Progress } from '@/components/ui/progress'
import { convertBytesToSize } from '@/helpers/convertBytesToSize.js'

const props = withDefaults(defineProps<{
	tab?: string
	mode?: 'list' | 'picker'
}>(), {
	tab: 'Audios',
	mode: 'list',
})

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
	const neededTab = tabs.find(t => t.name === props.tab)
	if (!neededTab) return
	activeTab.value = neededTab
})

const audios = computed(() => files.value?.files.filter(f => f.mimetype.startsWith('audio')) ?? [])

const message = useMessage()

async function upload(f: File) {
	if (!f.type.startsWith(activeTab.value.accept.split('*').at(0)!)) return

	try {
		await uploader.executeMutation({
			file: f,
		})
	} catch (error) {
		if (error instanceof RpcError) {
			message.error(error.message)
		}
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
		<div class="flex flex-col p-2 justify-between bg-card h-full border-r-2 border-border border-b-2 md:border-b-0 w-full md:min-w-max-w-52 md:max-w-52 gap-y-8">
			<div class="flex flex-col gap-1 w-full">
				<Button class="w-full flex items-center gap-2 justify-center" size="sm" variant="secondary">
					<MusicIcon class="size-4" />
					Audios
				</Button>
				<Button class="w-full flex items-center gap-2 justify-center" size="sm" variant="secondary" disabled>
					<ImageIcon class="size-4" />
					Images (soon)
				</Button>
			</div>
			<div class="flex flex-col gap-2 justify-center items-center">
				<FileButton
					class="size-28"
					:accept="activeTab.accept"
					:disabled="uploader.fetching.value || uploadedFilesSize >= ((1 << 20) * 100)"
					:loading="uploader.fetching.value"
					@file-selected="(file) => {
						if (!file) return
						upload(file)
					}"
				/>
				<Progress v-model="uploadedFilesSizeSlider" />
				<span>{{ convertBytesToSize(uploadedFilesSize) }} / {{ convertBytesToSize(maxFileSize) }}</span>
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

		<div class="p-4 w-full min-h-0 overflow-auto">
			<div v-if="activeTab.name === 'Audios'">
				<NAlert v-if="!audios.length" type="info">
					{{ t('filePicker.emptyText', { type: 'audios' }) }}
				</NAlert>

				<div v-else class="flex flex-col gap-4 overflow-y-auto">
					<div
						v-for="audio of audios"
						:key="audio.id"
						class="flex flex-col rounded-md bg-card p-2 border-2 border-border"
					>
						<h1 class="break-all">
							{{ audio.name }}
						</h1>
						<audio controls :src="api.computeFileUrl(audio.channelId, audio.id)" class="w-full" />

						<div class="flex justify-between flex-wrap items-end">
							<span>{{ convertBytesToSize(audio.size) }}</span>

							<div class="flex gap-1">
								<Button
									variant="destructive"
									size="sm"
									@click="async () => {
										await deleter.executeMutation({ id: audio.id })
										$emit('delete', audio.id)
									}"
								>
									{{ t('sharedButtons.delete') }}
								</Button>
								<Button size="sm" variant="secondary" @click="$emit('select', audio.id)">
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

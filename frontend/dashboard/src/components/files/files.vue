<script setup lang="ts">
import { RpcError } from '@protobuf-ts/runtime-rpc'
import { IconArchive } from '@tabler/icons-vue'
import {
	NAlert,
	NButton,
	NCard,
	NGrid,
	NGridItem,
	NIcon,
	NSpin,
	NText,
	NUpload,
	NUploadDragger,
	useMessage,
} from 'naive-ui'
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import type { FileMeta } from '@twir/api/messages/files/files'

import { useFileUpload, useFiles, userFileDelete } from '@/api/index.js'
import { convertBytesToSize } from '@/helpers/convertBytesToSize.js'

const props = withDefaults(defineProps<{
	tab: string
	mode: 'list' | 'picker'
}>(), {
	tab: 'Audios',
	mode: 'list',
})

defineEmits<{
	select: [id: string]
	delete: [id: string]
}>()

const { t } = useI18n()

const uploader = useFileUpload()
const deleter = userFileDelete()

const { data: files } = useFiles()

const uploadedFilesSize = computed(() => {
	if (!files.value?.files) return 0

	return files.value?.files.reduce((acc, curr) => acc + Number(curr.size), 0)
})

function computeFileUrl(f: FileMeta) {
	const query = new URLSearchParams({
		channel_id: f.channelId,
		file_id: f.id,
	})
	return `${window.location.origin}/api-old/files/?${query}`
}

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
		await uploader.mutateAsync(f)
	} catch (error) {
		if (error instanceof RpcError) {
			message.error(error.message)
		}
	}
}
</script>

<template>
	<div class="flex gap-5">
		<div class="flex flex-col pr-1">
			<div v-if="mode === 'list'" class="flex flex-col gap-1">
				<NButton
					v-for="tabItem of tabs"
					:key="tabItem.name"
					dashed
					size="large"
					:disabled="tabItem.disabled"
					:type="tabItem.name === activeTab.name ? 'success' : 'default'"
					block
					@click="activeTab = tabItem"
				>
					{{ tabItem.name }}
				</NButton>
			</div>

			<div>
				<NUpload
					multiple
					directory-dnd
					:max="1"
					:accept="activeTab.accept"
					:file-list="[]"
					:disabled="uploader.isLoading.value || uploadedFilesSize >= ((1 << 20) * 100)"
					@change="(data) => {
						if (!data.fileList.length) return
						upload(data.file.file!)
					}"
				>
					<NUploadDragger>
						<div v-if="!uploader.isLoading.value">
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
			</div>
		</div>

		<div v-if="activeTab.name === 'Audios'">
			<NAlert v-if="!audios.length" type="info">
				{{ t('filePicker.emptyText', { type: 'audios' }) }}
			</NAlert>

			<NGrid v-else cols="1 s:1 m:2 l:3" responsive="screen" :x-gap="8" :y-gap="8">
				<NGridItem
					v-for="f of audios"
					:key="f.id"
					:span="1"
				>
					<NCard
						:title="`${f.name} (${convertBytesToSize(Number(f.size))})`"
						size="small"
						segmented
					>
						<template #header-extra>
							<NButton
								secondary
								type="error"
								size="small"
								@click="async () => {
									await deleter.mutateAsync(f.id)
									$emit('delete', f.id)
								}"
							>
								{{ t('sharedButtons.delete') }}
							</NButton>
						</template>

						<audio controls :src="computeFileUrl(f)" class="w-full" />

						<template v-if="mode === 'picker'" #footer>
							<NButton block @click="$emit('select', f.id)">
								{{ t('sharedButtons.select') }}
							</NButton>
						</template>
					</NCard>
				</NGridItem>
			</NGrid>
		</div>
	</div>
</template>

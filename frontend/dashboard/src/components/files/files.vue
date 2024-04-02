<script setup lang="ts">
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { IconArchive } from '@tabler/icons-vue';
import { FileMeta } from '@twir/api/messages/files/files';
import {
	NAlert,
	NButton,
	NCard,
	NDivider,
	NGrid,
	NGridItem,
	NIcon,
	NSpin,
	NText,
	NUpload,
	NUploadDragger,
	useMessage,
} from 'naive-ui';
import { computed, onMounted, ref } from 'vue';
import { useI18n } from 'vue-i18n';

import { useFiles, useFileUpload, userFileDelete } from '@/api/index.js';
import { convertBytesToSize } from '@/helpers/convertBytesToSize.js';


const { t } = useI18n();

const uploader = useFileUpload();
const deleter = userFileDelete();

const { data: files } = useFiles();

const uploadedFilesSize = computed(() => {
	if (!files.value?.files) return 0;

	return files.value?.files.reduce((acc, curr) => acc + Number(curr.size), 0);
});

const computeFileUrl = (f: FileMeta) => {
	const query = new URLSearchParams({
		channel_id: f.channelId,
		file_id: f.id,
	});
	return `${window.location.origin}/api/files/?${query}`;
};

const props = withDefaults(defineProps<{
	tab: string,
	mode: 'list' | 'picker',
}>(), {
	tab: 'Audios',
	mode: 'list',
});

type Tab = { name: string, disabled?: boolean, accept: string }
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
];
const activeTab = ref<Tab>(tabs.at(0)!);

onMounted(() => {
	const neededTab = tabs.find(t => t.name === props.tab);
	if (!neededTab) return;
	activeTab.value = neededTab;
});

const audios = computed(() => files.value?.files.filter(f => f.mimetype.startsWith('audio')) ?? []);

defineEmits<{
	select: [id: string],
	delete: [id: string]
}>();

const message = useMessage();

async function upload(f: File) {
	if (!f.type.startsWith(activeTab.value.accept.split('*').at(0)!)) return;

	try {
		await uploader.mutateAsync(f);
	} catch (error) {
		if (error instanceof RpcError) {
			message.error(error.message);
		}
	}
}
</script>

<template>
	<div class="flex gap-5">
		<div class="flex flex-col pr-1 border-r-[color:var(--n-border-color)] border-r border-solid">
			<div v-if="mode === 'list'" class="flex flex-col gap-1">
				<n-button
					v-for="tab of tabs"
					:key="tab.name"
					dashed
					size="large"
					:disabled="tab.disabled"
					:type="tab.name === activeTab.name ? 'success' : 'default'"
					block
					@click="activeTab = tab"
				>
					{{ tab.name }}
				</n-button>
				<n-divider />
			</div>

			<div>
				<n-upload
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
					<n-upload-dragger>
						<div v-if="!uploader.isLoading.value">
							<div class="mb-3">
								<n-icon size="30" :depth="3">
									<IconArchive />
								</n-icon>
							</div>
							<n-text class="text-xs">
								{{ t('filePicker.innerText', { type: activeTab.name.toLowerCase() }) }}
							</n-text>
						</div>
						<n-spin v-else />
					</n-upload-dragger>
				</n-upload>

				<n-text>
					{{
						t('filePicker.usedSpace', {
							used: convertBytesToSize(uploadedFilesSize),
							max: 100
						})
					}}
				</n-text>
			</div>
		</div>

		<div v-if="activeTab.name === 'Audios'">
			<n-alert v-if="!audios.length" type="info">
				{{ t('filePicker.emptyText', { type: 'audios' }) }}
			</n-alert>

			<n-grid v-else cols="1 s:1 m:2 l:3" responsive="screen" :x-gap="8" :y-gap="8">
				<n-grid-item
					v-for="f of audios"
					:key="f.id"
					:span="1"
				>
					<n-card
						:title="`${f.name} (${convertBytesToSize(Number(f.size))})`"
						size="small"
						segmented
					>
						<template #header-extra>
							<n-button
								secondary
								type="error"
								size="small"
								@click="async () => {
									await deleter.mutateAsync(f.id)
									$emit('delete', f.id)
								}"
							>
								{{ t('sharedButtons.delete') }}
							</n-button>
						</template>

						<audio controls :src="computeFileUrl(f)" class="w-full" />

						<template v-if="mode === 'picker'" #footer>
							<n-button block @click="$emit('select', f.id)">
								{{ t('sharedButtons.select') }}
							</n-button>
						</template>
					</n-card>
				</n-grid-item>
			</n-grid>
		</div>
	</div>
</template>

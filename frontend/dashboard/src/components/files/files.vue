<script setup lang="ts">
import { IconArchive } from '@tabler/icons-vue';
import { FileMeta } from '@twir/grpc/generated/api/api/files';
import { NIcon, NText, NUpload, NUploadDragger, NCard, NGrid, NGridItem, NButton, NDivider, NSpin, NAlert } from 'naive-ui';
import { computed, onMounted, ref } from 'vue';
import { useI18n } from 'vue-i18n';

import { useFileUpload, useFiles, userFileDelete } from '@/api/index.js';
import { convertBytesToSize } from '@/helpers/convertBytesToSize.js';


const { t } = useI18n();

const uploader = useFileUpload();
const deleter = userFileDelete();

async function upload(f: File) {
	await uploader.mutateAsync(f);
}

const { data: files } = useFiles();

const uploadedFilesSize = computed(() => {
	if (!files.value?.files) return 0;

	return files.value?.files.reduce((acc, curr) => acc + Number(curr.size), 0);
});

const computeFileUrl = (f: FileMeta) => {
	return `${window.location.origin}/cdn/twir/channels/${f.channelId}/${f.id}`;
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
</script>

<template>
	<div style="display: flex; gap: 20px;">
		<div class="sidebar">
			<div v-if="mode === 'list'" style="display: flex; flex-direction: column; gap: 4px;">
				<n-button
					v-for="t of tabs"
					:key="t.name"
					dashed
					size="large"
					:disabled="t.disabled"
					:type="t.name === activeTab.name ? 'success' : 'default'"
					block
					@click="activeTab = t"
				>
					{{ t.name }}
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
							<div style="margin-bottom: 12px">
								<n-icon size="30" :depth="3">
									<IconArchive />
								</n-icon>
							</div>
							<n-text style="font-size: 13px">
								Select or drop {{ activeTab.name.toLowerCase() }}
							</n-text>
						</div>
						<n-spin v-else />
					</n-upload-dragger>
				</n-upload>

				<n-text>Used {{ convertBytesToSize(uploadedFilesSize) }} / 100 MB</n-text>
			</div>
		</div>

		<div v-if="activeTab.name === 'Audios'">
			<n-alert v-if="!audios.length" type="info">
				No audios uploaded
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

						<audio controls :src="computeFileUrl(f)" style="width: 100%;" />

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

<style scoped>
.sidebar {
	display: flex;
	flex-direction: column;
	border-right: 1px solid #373636;
	padding-right: 5px;
}
</style>
@/helpers/convertBytesToSize.js

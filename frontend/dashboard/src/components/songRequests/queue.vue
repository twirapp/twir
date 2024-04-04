<script setup lang="ts">
import { dragAndDrop } from '@formkit/drag-and-drop/vue';
import {
	IconTrash,
	IconBan,
	IconGripVertical,
} from '@tabler/icons-vue';
import {
	NCard,
	NButton,
	NTime,
	NPopconfirm,
} from 'naive-ui';
import { storeToRefs } from 'pinia';
import { computed, ref } from 'vue';
import { useI18n } from 'vue-i18n';

import { useYoutubeSocket, Video } from '@/components/songRequests/hook.js';
import {
	Table,
	TableBody,
	TableCell,
	TableHead,
	TableHeader,
	TableRow,
} from '@/components/ui/table';
import { convertMillisToTime } from '@/helpers/convertMillisToTime.js';

const socket = useYoutubeSocket();
const { videos } = storeToRefs(socket);

const { t } = useI18n();

const totalSongsLength = computed(() => {
	return convertMillisToTime(videos.value.reduce((acc, cur) => acc + cur.duration * 1000, 0));
});

const parentRef = ref<HTMLElement>();
dragAndDrop({
	// eslint-disable-next-line @typescript-eslint/ban-ts-comment
	// @ts-expect-error
	parent: parentRef,
	// eslint-disable-next-line @typescript-eslint/ban-ts-comment
	// @ts-expect-error
	values: videos,
	dragHandle: '.drag-handle',
	draggable(child) {
		return !child.classList.contains('no-drag');
	},
	handleEnd(data) {
		const item = data.targetData.node.data.value as Video;
		socket.moveVideo(item.id, data.targetData.node.data.index);
	},
});
</script>

<template>
	<n-card
		:title="t('songRequests.table.title')"
		content-style="padding: 0;"
		header-style="padding: 10px;"
		segmented
	>
		<template #header-extra>
			<n-button tertiary size="small" @click="$emit('deleteAllVideos')">
				<IconTrash />
			</n-button>
		</template>
		<Table class="w-full">
			<TableHeader>
				<TableRow>
					<TableHead class="w-[1%]"></TableHead>
					<TableHead class="w-[5%]">
						#
					</TableHead>
					<TableHead>Name</TableHead>
					<TableHead>Author</TableHead>
					<TableHead></TableHead>
					<TableHead>
						Duration
					</TableHead>
					<TableHead>
						Actions
					</TableHead>
				</TableRow>
			</TableHeader>
			<TableBody ref="parentRef">
				<TableRow v-for="(video, index) of videos" :key="video.id">
					<TableCell>
						<IconGripVertical class="w-4 drag-handle cursor-move" />
					</TableCell>
					<TableCell>
						{{ index + 1 }}
					</TableCell>
					<TableCell>
						<div class="flex items-center gap-2">
							<span>{{ video.title }}</span>

							<n-popconfirm
								:positive-text="t('deleteConfirmation.confirm')"
								:negative-text="t('deleteConfirmation.cancel')"
								@positive-click="socket.banSong(video.videoId)"
							>
								<template #trigger>
									<n-button size="tiny" text>
										<IconBan />
									</n-button>
								</template>
								{{ t('songRequests.ban.songConfirm') }}
							</n-popconfirm>
						</div>
					</TableCell>
					<TableCell>
						<div class="flex items-center gap-2">
							<span>{{ video.orderedByDisplayName || video.orderedByName }}</span>
							<n-popconfirm
								:positive-text="t('deleteConfirmation.confirm')"
								:negative-text="t('deleteConfirmation.cancel')"
								@positive-click="socket.banUser(video.orderedById)"
							>
								<template #trigger>
									<n-button size="tiny" text>
										<IconBan />
									</n-button>
								</template>
								{{ t('songRequests.ban.userConfirm') }}
							</n-popconfirm>
						</div>
					</TableCell>
					<TableCell>
						<n-time type="relative" :time="0" :to="Date.now() - new Date(video.createdAt).getTime()" />
					</TableCell>
					<TableCell>
						{{ convertMillisToTime(video.duration * 1000) }}
					</TableCell>
					<TableCell>
						<n-button size="tiny" type="error" text @click="socket.deleteVideo(video.id)">
							<IconTrash />
						</n-button>
					</TableCell>
				</TableRow>
				<TableRow class="no-drag">
					<TableCell></TableCell>
					<TableCell>
						{{ videos.length }}
					</TableCell>
					<TableCell></TableCell>
					<TableCell></TableCell>
					<TableCell></TableCell>
					<TableCell>
						{{ totalSongsLength }}
					</TableCell>
					<TableCell></TableCell>
				</TableRow>
			</TableBody>
		</Table>
	</n-card>
</template>

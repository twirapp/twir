<script setup lang="ts">
import { dragAndDrop } from '@formkit/drag-and-drop/vue'
import { IconBan, IconGripVertical, IconTrash } from '@tabler/icons-vue'
import { NCard, NPopconfirm, NTime } from 'naive-ui'
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import ActionConfirm from '../ui/action-confirm.vue'

import type { Video } from '@/components/songRequests/hook.js'

import { useYoutubeSocket } from '@/components/songRequests/hook.js'
import { Button } from '@/components/ui/button'
import {
	Table,
	TableBody,
	TableCell,
	TableHead,
	TableHeader,
	TableRow,
} from '@/components/ui/table'
import { convertMillisToTime } from '@/helpers/convertMillisToTime.js'

const { videos, moveVideo, banSong, banUser, deleteVideo, deleteAllVideos } = useYoutubeSocket()

const { t } = useI18n()

const totalSongsLength = computed(() => {
	return convertMillisToTime(videos.value.reduce((acc, cur) => acc + cur.duration * 1000, 0))
})

const parentRef = ref<HTMLElement>()
dragAndDrop({
	// @ts-expect-error
	parent: parentRef,
	// @ts-expect-error
	values: videos,
	dragHandle: '.drag-handle',
	draggable(child) {
		return !child.classList.contains('no-drag')
	},
	handleEnd(data) {
		const item = data.targetData.node.data.value as Video
		moveVideo(item.id, data.targetData.node.data.index)
	},
})

const showConfirmClear = ref(false)
</script>

<template>
	<NCard
		:title="t('songRequests.table.title')"
		content-style="padding: 0;"
		header-style="padding: 10px;"
		segmented
	>
		<template #header-extra>
			<Button
				size="icon"
				class="size-8"
				variant="secondary"
				:disabled="!videos.length"
				@click="showConfirmClear = true"
			>
				<IconTrash class="size-4" />
			</Button>
		</template>
		<Table class="w-full">
			<TableHeader>
				<TableRow>
					<TableHead class="w-[1%]"></TableHead>
					<TableHead class="w-[5%]"> # </TableHead>
					<TableHead>Name</TableHead>
					<TableHead>Author</TableHead>
					<TableHead></TableHead>
					<TableHead> Duration </TableHead>
					<TableHead> Actions </TableHead>
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

							<NPopconfirm
								:positive-text="t('deleteConfirmation.confirm')"
								:negative-text="t('deleteConfirmation.cancel')"
								@positive-click="banSong(video.videoId)"
							>
								<template #trigger>
									<Button class="min-w-5" size="icon" variant="ghost">
										<IconBan class="size-5" />
									</Button>
								</template>
								{{ t('songRequests.ban.songConfirm') }}
							</NPopconfirm>
						</div>
					</TableCell>
					<TableCell>
						<div class="flex items-center gap-2">
							<span>{{ video.orderedByDisplayName || video.orderedByName }}</span>
							<NPopconfirm
								:positive-text="t('deleteConfirmation.confirm')"
								:negative-text="t('deleteConfirmation.cancel')"
								@positive-click="banUser(video.orderedById)"
							>
								<template #trigger>
									<Button class="min-w-5" size="icon" variant="ghost">
										<IconBan class="size-5" />
									</Button>
								</template>
								{{ t('songRequests.ban.userConfirm') }}
							</NPopconfirm>
						</div>
					</TableCell>
					<TableCell>
						<NTime
							type="relative"
							:time="0"
							:to="Date.now() - new Date(video.createdAt).getTime()"
						/>
					</TableCell>
					<TableCell>
						{{ convertMillisToTime(video.duration * 1000) }}
					</TableCell>
					<TableCell>
						<Button
							class="min-w-5"
							size="icon"
							variant="destructive"
							@click="deleteVideo(video.id)"
						>
							<IconTrash class="size-5" />
						</Button>
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
	</NCard>

	<ActionConfirm
		v-model:open="showConfirmClear"
		:confirm-text="t('songRequests.settings.confirmClearQueue')"
		@confirm="deleteAllVideos"
	/>
</template>

<script setup lang="ts">
import { dragAndDrop } from '@formkit/drag-and-drop/vue'
import { Ban, GripVertical, Trash2 } from 'lucide-vue-next'
import { formatDistanceToNow } from 'date-fns'
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import ActionConfirm from '../ui/action-confirm.vue'

import type { Video } from '@/components/songRequests/hook.js'

import { useYoutubeSocket } from '@/components/songRequests/hook.js'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardTitle } from '@/components/ui/card'
import {
	Table,
	TableBody,
	TableCell,
	TableHead,
	TableHeader,
	TableRow,
} from '@/components/ui/table'
import {
	AlertDialog,
	AlertDialogAction,
	AlertDialogCancel,
	AlertDialogContent,
	AlertDialogFooter,
	AlertDialogHeader,
	AlertDialogTitle,
	AlertDialogTrigger,
} from '@/components/ui/alert-dialog'
import { convertMillisToTime } from '@/helpers/convertMillisToTime.js'

const { videos, moveVideo, banSong, banUser, deleteVideo, deleteAllVideos } = useYoutubeSocket()

const { t } = useI18n()

const totalSongsLength = computed(() => {
	return convertMillisToTime(videos.value.reduce((acc, cur) => acc + cur.duration * 1000, 0))
})

const parentRef = ref<HTMLElement>()
dragAndDrop({
	parent: parentRef,
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

function formatRelativeTime(dateStr: string) {
	return formatDistanceToNow(new Date(dateStr), { addSuffix: true })
}
</script>

<template>
	<Card class="p-0">
		<CardContent class="p-0">
			<div class="flex flex-row justify-between items-center px-2 py-2 border-b">
				<CardTitle class="text-base">{{ t('songRequests.table.title') }}</CardTitle>
				<Button
					size="icon"
					class="size-8"
					variant="secondary"
					:disabled="!videos.length"
					@click="showConfirmClear = true"
				>
					<Trash2 class="size-4" />
				</Button>
			</div>
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
							<GripVertical class="w-4 drag-handle cursor-move" />
						</TableCell>
						<TableCell>
							{{ index + 1 }}
						</TableCell>
						<TableCell>
							<div class="flex items-center gap-2">
								<span>{{ video.title }}</span>

								<AlertDialog>
									<AlertDialogTrigger as-child>
										<Button class="min-w-5" size="icon" variant="ghost">
											<Ban class="size-5" />
										</Button>
									</AlertDialogTrigger>
									<AlertDialogContent>
										<AlertDialogHeader>
											<AlertDialogTitle>{{ t('songRequests.ban.songConfirm') }}</AlertDialogTitle>
										</AlertDialogHeader>
										<AlertDialogFooter>
											<AlertDialogCancel>{{ t('deleteConfirmation.cancel') }}</AlertDialogCancel>
											<AlertDialogAction @click="banSong(video.videoId)">
												{{ t('deleteConfirmation.confirm') }}
											</AlertDialogAction>
										</AlertDialogFooter>
									</AlertDialogContent>
								</AlertDialog>
							</div>
						</TableCell>
						<TableCell>
							<div class="flex items-center gap-2">
								<span>{{ video.orderedByDisplayName || video.orderedByName }}</span>
								<AlertDialog>
									<AlertDialogTrigger as-child>
										<Button class="min-w-5" size="icon" variant="ghost">
											<Ban class="size-5" />
										</Button>
									</AlertDialogTrigger>
									<AlertDialogContent>
										<AlertDialogHeader>
											<AlertDialogTitle>{{ t('songRequests.ban.userConfirm') }}</AlertDialogTitle>
										</AlertDialogHeader>
										<AlertDialogFooter>
											<AlertDialogCancel>{{ t('deleteConfirmation.cancel') }}</AlertDialogCancel>
											<AlertDialogAction @click="banUser(video.orderedById)">
												{{ t('deleteConfirmation.confirm') }}
											</AlertDialogAction>
										</AlertDialogFooter>
									</AlertDialogContent>
								</AlertDialog>
							</div>
						</TableCell>
						<TableCell>
							{{ formatRelativeTime(video.createdAt) }}
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
								<Trash2 class="size-5" />
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
		</CardContent>
	</Card>

	<ActionConfirm
		v-model:open="showConfirmClear"
		:confirm-text="t('songRequests.settings.confirmClearQueue')"
		@confirm="deleteAllVideos"
	/>
</template>

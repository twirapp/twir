<script setup lang="ts">
import { dragAndDrop } from '@formkit/drag-and-drop/vue'
import { Ban, GripVertical, Trash2 } from 'lucide-vue-next'
import { formatDistanceToNow } from 'date-fns'

import ActionConfirm from '../ui/action-confirm.vue'

import type { Video } from '#layers/dashboard/components/songRequests/hook.js'
import { useYoutubeSocket } from '#layers/dashboard/components/songRequests/hook.js'
import { convertMillisToTime } from '#layers/dashboard/helpers/convertMillisToTime.js'

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
	<UiCard class="p-0">
		<UiCardContent class="p-0">
			<div class="flex flex-row justify-between items-center px-2 py-2 border-b">
				<UiCardTitle class="text-base">{{ t('songRequests.table.title') }}</UiCardTitle>
				<UiButton
					size="icon"
					class="size-8"
					variant="secondary"
					:disabled="!videos.length"
					@click="showConfirmClear = true"
				>
					<Trash2 class="size-4" />
				</UiButton>
			</div>
			<UiTable class="w-full">
				<UiTableHeader>
					<UiTableRow>
						<UiTableHead class="w-[1%]"></UiTableHead>
						<UiTableHead class="w-[5%]"> # </UiTableHead>
						<UiTableHead>Name</UiTableHead>
						<UiTableHead>Author</UiTableHead>
						<UiTableHead></UiTableHead>
						<UiTableHead> Duration </UiTableHead>
						<UiTableHead> Actions </UiTableHead>
					</UiTableRow>
				</UiTableHeader>
				<UiTableBody ref="parentRef">
					<UiTableRow v-for="(video, index) of videos" :key="video.id">
						<UiTableCell>
							<GripVertical class="w-4 drag-handle cursor-move" />
						</UiTableCell>
						<UiTableCell>
							{{ index + 1 }}
						</UiTableCell>
						<UiTableCell>
							<div class="flex items-center gap-2">
								<span>{{ video.title }}</span>

								<UiAlertDialog>
									<UiAlertDialogTrigger as-child>
										<UiButton class="min-w-5" size="icon" variant="ghost">
											<Ban class="size-5" />
										</UiButton>
									</UiAlertDialogTrigger>
									<UiAlertDialogContent>
										<UiAlertDialogHeader>
											<AlertDialogTitle>{{ t('songRequests.ban.songConfirm') }}</AlertDialogTitle>
										</UiAlertDialogHeader>
										<UiAlertDialogFooter>
											<UiAlertDialogCancel>{{ t('deleteConfirmation.cancel') }}</UiAlertDialogCancel>
											<UiAlertDialogAction @click="banSong(video.videoId)">
												{{ t('deleteConfirmation.confirm') }}
											</UiAlertDialogAction>
										</UiAlertDialogFooter>
									</UiAlertDialogContent>
								</UiAlertDialog>
							</div>
						</UiTableCell>
						<UiTableCell>
							<div class="flex items-center gap-2">
								<span>{{ video.orderedByDisplayName || video.orderedByName }}</span>
								<UiAlertDialog>
									<UiAlertDialogTrigger as-child>
										<UiButton class="min-w-5" size="icon" variant="ghost">
											<Ban class="size-5" />
										</UiButton>
									</UiAlertDialogTrigger>
									<UiAlertDialogContent>
										<UiAlertDialogHeader>
											<AlertDialogTitle>{{ t('songRequests.ban.userConfirm') }}</AlertDialogTitle>
										</UiAlertDialogHeader>
										<UiAlertDialogFooter>
											<UiAlertDialogCancel>{{ t('deleteConfirmation.cancel') }}</UiAlertDialogCancel>
											<UiAlertDialogAction @click="banUser(video.orderedById)">
												{{ t('deleteConfirmation.confirm') }}
											</UiAlertDialogAction>
										</UiAlertDialogFooter>
									</UiAlertDialogContent>
								</UiAlertDialog>
							</div>
						</UiTableCell>
						<UiTableCell>
							{{ formatRelativeTime(video.createdAt) }}
						</UiTableCell>
						<UiTableCell>
							{{ convertMillisToTime(video.duration * 1000) }}
						</UiTableCell>
						<UiTableCell>
							<Button
								class="min-w-5"
								size="icon"
								variant="destructive"
								@click="deleteVideo(video.id)"
							>
								<Trash2 class="size-5" />
							</Button>
						</UiTableCell>
					</UiTableRow>
					<UiTableRow class="no-drag">
						<UiTableCell></UiTableCell>
					<UiTableCell>
							{{ videos.length }}
						</UiTableCell>
						<UiTableCell></UiTableCell>
						<UiTableCell></UiTableCell>
						<UiTableCell></UiTableCell>
						<UiTableCell>
							{{ totalSongsLength }}
						</UiTableCell>
						<UiTableCell></UiTableCell>
					</UiTableRow>
				</UiTableBody>
			</UiTable>
		</UiCardContent>
	</UiCard>

	<ActionConfirm
		v-model:open="showConfirmClear"
		:confirm-text="t('songRequests.settings.confirmClearQueue')"
		@confirm="deleteAllVideos"
	/>
</template>

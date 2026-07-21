<script setup lang="ts">
import { dragAndDrop } from '@formkit/drag-and-drop/vue'
import { formatDistanceToNow } from 'date-fns'
import { useProfile } from '~~/layers/dashboard/api/auth'
import { useSongRequestsApi } from '~~/layers/dashboard/api/song-requests'
import { useSongRequestGql } from '~~/layers/dashboard/composables/useSongRequestGql.js'
import { convertMillisToTime } from '~~/layers/dashboard/helpers/convertMillisToTime.js'

import type { SongRequestsSettingsOpts } from '~/gql/graphql.js'

import ActionConfirm from '@/components/ui/action-confirm'
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

const { data: profile } = useProfile()
const channelId = computed(() => profile.value?.selectedDashboardId ?? '')

const { queue, deleteFromQueue, clearQueue, reorder } = useSongRequestGql(channelId)

const youtubeModuleManager = useSongRequestsApi()
const { data: youtubeSettings } = youtubeModuleManager.useSongRequestQuery()
const youtubeModuleUpdater = youtubeModuleManager.useSongRequestMutation()

const { t } = useI18n()

function toSettingsOpts(settings: Record<string, unknown>): SongRequestsSettingsOpts {
	const { channelApiKey, __typename, ...rest } = settings
	return rest as unknown as SongRequestsSettingsOpts
}

const totalSongsLength = computed(() => {
	return convertMillisToTime(queue.value.reduce((acc, cur) => acc + cur.durationSeconds * 1000, 0))
})

const parentRef = ref<HTMLElement>()
dragAndDrop({
	parent: parentRef,
	values: queue,
	dragHandle: '.drag-handle',
	draggable(child) {
		return !child.classList.contains('no-drag')
	},
	handleEnd() {
		reorder(queue.value.map((item) => item.id))
	},
})

async function banUser(userName: string) {
	if (!youtubeSettings.value?.songRequests) return
	const settings = toSettingsOpts(youtubeSettings.value.songRequests as Record<string, unknown>)

	await youtubeModuleUpdater.executeMutation({
		opts: {
			...settings,
			denyList: {
				...settings.denyList,
				users: [...settings.denyList!.users, userName],
			},
		},
	})

	const userVideos = queue.value.filter((video) => video.orderedByName === userName)
	for (const video of userVideos) {
		deleteFromQueue(video.id)
	}
}

async function banSong(queueItemId: string) {
	if (!youtubeSettings.value?.songRequests) return
	const settings = toSettingsOpts(youtubeSettings.value.songRequests as Record<string, unknown>)

	const video = queue.value.find((v) => v.id === queueItemId)
	const videoId = video?.songLink?.match(/(?:v=|youtu\.be\/)([^&?/]+)/)?.[1] ?? ''

	await youtubeModuleUpdater.executeMutation({
		opts: {
			...settings,
			denyList: {
				...settings.denyList,
				songs: [...settings.denyList!.songs, videoId],
			},
		},
	})

	deleteFromQueue(queueItemId)
}

const showConfirmClear = ref(false)

function formatRelativeTime(dateStr: string) {
	return formatDistanceToNow(new Date(dateStr), { addSuffix: true })
}
</script>

<template>
	<Card class="p-0">
		<CardContent class="p-0">
			<div class="flex flex-row items-center justify-between border-b px-2 py-2">
				<CardTitle class="text-base">{{ t('songRequests.table.title') }}</CardTitle>
				<Button
					size="icon"
					class="size-8"
					variant="secondary"
					:disabled="!queue.length"
					@click="showConfirmClear = true"
				>
					<Icon
						name="lucide:trash2"
						class="size-4"
					/>
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
					<TableRow
						v-for="(video, index) of queue"
						:key="video.id"
					>
						<TableCell>
							<Icon
								name="lucide:grip-vertical"
								class="drag-handle w-4 cursor-move"
							/>
						</TableCell>
						<TableCell>
							{{ index + 1 }}
						</TableCell>
						<TableCell>
							<div class="flex items-center gap-2">
								<span>{{ video.title }}</span>

								<AlertDialog>
									<AlertDialogTrigger as-child>
										<Button
											class="min-w-5"
											size="icon"
											variant="ghost"
										>
											<Icon
												name="lucide:ban"
												class="size-5"
											/>
										</Button>
									</AlertDialogTrigger>
									<AlertDialogContent>
										<AlertDialogHeader>
											<AlertDialogTitle>{{ t('songRequests.ban.songConfirm') }}</AlertDialogTitle>
										</AlertDialogHeader>
										<AlertDialogFooter>
											<AlertDialogCancel>{{ t('deleteConfirmation.cancel') }}</AlertDialogCancel>
											<AlertDialogAction @click="banSong(video.id)">
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
										<Button
											class="min-w-5"
											size="icon"
											variant="ghost"
										>
											<Icon
												name="lucide:ban"
												class="size-5"
											/>
										</Button>
									</AlertDialogTrigger>
									<AlertDialogContent>
										<AlertDialogHeader>
											<AlertDialogTitle>{{ t('songRequests.ban.userConfirm') }}</AlertDialogTitle>
										</AlertDialogHeader>
										<AlertDialogFooter>
											<AlertDialogCancel>{{ t('deleteConfirmation.cancel') }}</AlertDialogCancel>
											<AlertDialogAction @click="banUser(video.orderedByName)">
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
							{{ convertMillisToTime(video.durationSeconds * 1000) }}
						</TableCell>
						<TableCell>
							<Button
								class="min-w-5"
								size="icon"
								variant="destructive"
								@click="deleteFromQueue(video.id)"
							>
								<Icon
									name="lucide:trash"
									class="size-5"
								/>
							</Button>
						</TableCell>
					</TableRow>
					<TableRow class="no-drag">
						<TableCell></TableCell>
						<TableCell>
							{{ queue.length }}
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
		@confirm="clearQueue"
	/>
</template>

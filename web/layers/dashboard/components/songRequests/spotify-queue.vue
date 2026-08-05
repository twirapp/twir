<script setup lang="ts">
import { formatDistanceToNow } from 'date-fns'
import { toast } from 'vue-sonner'
import { useProfile } from '~~/layers/dashboard/api/auth'
import { useIntegrationsPageData } from '~~/layers/dashboard/api/integrations/integrations-page.js'
import { useSongRequestsApi } from '~~/layers/dashboard/api/song-requests'
import { convertMillisToTime } from '~~/layers/dashboard/helpers/convertMillisToTime.js'

import { SpotifySongRequestStatus } from '~/gql/graphql.js'

import ActionConfirm from '@/components/ui/action-confirm'
import { Badge } from '@/components/ui/badge'
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

const { t } = useI18n()

const songRequestsApi = useSongRequestsApi()
const settingsQuery = songRequestsApi.useSongRequestQuery()

const integrationsPage = useIntegrationsPageData()
const spotifyAuthLink = computed(() => integrationsPage.spotifyAuthLink.value)

const capabilities = computed(() => settingsQuery.data.value?.songRequests?.spotifyCapabilities)

const queuePaused = computed(() => !capabilities.value?.canUseSpotify)
const queueQuery = songRequestsApi.useSpotifyQueueQuery(queuePaused)

const { data: profile } = useProfile()
const channelId = computed(() => profile.value?.selectedDashboardId ?? '')
const queueSubscription = songRequestsApi.useSpotifyQueueSubscription(channelId, queuePaused)

const queue = computed(() => {
	if (queueSubscription.data.value !== undefined) {
		return queueSubscription.data.value.spotifySongRequestsQueueUpdated.requests ?? []
	}
	return queueQuery.data.value?.spotifySongRequestsQueue.requests ?? []
})
const currentDevice = computed(() => {
	if (queueSubscription.data.value !== undefined) {
		return queueSubscription.data.value.spotifySongRequestsQueueUpdated.currentDevice ?? null
	}
	return queueQuery.data.value?.spotifySongRequestsQueue.currentDevice
})

function connectSpotify() {
	if (!spotifyAuthLink.value) return
	window.open(spotifyAuthLink.value, 'Twir connect integration', 'width=800,height=600')
}

watch(
	() => integrationsPage.spotifyData.value,
	(data, previous) => {
		if (!data || data === previous) return
		settingsQuery.executeQuery({ requestPolicy: 'network-only' })
		queueQuery.executeQuery({ requestPolicy: 'network-only' })
	}
)

const skipMutation = songRequestsApi.useSpotifySkipMutation()
const cancelMutation = songRequestsApi.useSpotifyCancelMutation()
const refreshDeviceMutation = songRequestsApi.useSpotifyRefreshDeviceMutation()

const totalSongsLength = computed(() => {
	return convertMillisToTime(queue.value.reduce((acc, cur) => acc + cur.durationMs, 0))
})

const skipTargetId = ref<string | null>(null)
const cancelTargetId = ref<string | null>(null)

const showSkipConfirm = computed({
	get: () => skipTargetId.value !== null,
	set: (v: boolean) => {
		if (!v) skipTargetId.value = null
	},
})

const showCancelConfirm = computed({
	get: () => cancelTargetId.value !== null,
	set: (v: boolean) => {
		if (!v) cancelTargetId.value = null
	},
})

interface RequestIdMutation {
	executeMutation: (vars: { requestId: string }) => Promise<{ error?: { message: string } }>
}

async function executeWithToast(mutation: RequestIdMutation, requestId: string) {
	const result = await mutation.executeMutation({ requestId })
	if (result.error) {
		toast.error(result.error.message, { duration: 5000 })
		return
	}
	toast.success(t('sharedTexts.saved'), { duration: 2500 })
}

async function skipRequest() {
	if (!skipTargetId.value) return
	await executeWithToast(skipMutation, skipTargetId.value)
}

async function cancelRequest() {
	if (!cancelTargetId.value) return
	await executeWithToast(cancelMutation, cancelTargetId.value)
}

async function refreshDevice() {
	const result = await refreshDeviceMutation.executeMutation({})
	if (result.error) {
		toast.error(result.error.message, { duration: 5000 })
		return
	}
	toast.success(t('songRequests.spotify.deviceRefreshed'), { duration: 2500 })
}

function statusVariant(status: SpotifySongRequestStatus): 'default' | 'secondary' | 'destructive' | 'outline' {
	switch (status) {
		case SpotifySongRequestStatus.Playing:
			return 'default'
		case SpotifySongRequestStatus.Queued:
			return 'secondary'
		case SpotifySongRequestStatus.SkippedByTwir:
		case SpotifySongRequestStatus.CancelledPendingSkip:
		case SpotifySongRequestStatus.RemovedOrReconciled:
			return 'destructive'
		default:
			return 'outline'
	}
}

function statusLabel(status: SpotifySongRequestStatus): string {
	return t(`songRequests.spotify.status.${status}`)
}

function formatRelativeTime(dateStr: string) {
	return formatDistanceToNow(new Date(dateStr), { addSuffix: true })
}
</script>

<template>
	<Card class="mb-4">
		<CardContent class="flex flex-wrap items-center gap-2 p-4">
			<Icon
				name="simple-icons:spotify"
				class="size-5 text-[#1DB954]"
			/>
			<template v-if="!capabilities?.connected">
				<span class="text-sm">{{ t('songRequests.spotify.notConnected') }}</span>
				<Button
					size="sm"
					variant="outline"
					:disabled="!spotifyAuthLink"
					@click="connectSpotify"
				>
					{{ t('songRequests.spotify.connect') }}
				</Button>
			</template>
			<template v-else-if="!capabilities.hasPlaybackScope">
				<span class="text-sm">{{ t('songRequests.spotify.missingScope') }}</span>
				<Button
					size="sm"
					variant="outline"
					:disabled="!spotifyAuthLink"
					@click="connectSpotify"
				>
					{{ t('songRequests.spotify.reconnect') }}
				</Button>
			</template>
			<template v-else>
				<span class="text-sm">
					{{
						currentDevice
							? t('songRequests.spotify.device', {
									name: currentDevice.name,
									type: currentDevice.type,
								})
							: t('songRequests.spotify.noDevice')
					}}
				</span>
				<Button
					size="sm"
					variant="outline"
					@click="refreshDevice"
				>
					<Icon
						name="lucide:refresh-cw"
						class="size-4"
					/>
					{{ t('songRequests.spotify.refreshDevice') }}
				</Button>
			</template>
		</CardContent>
	</Card>

	<Card class="p-0">
		<CardContent class="p-0">
			<div class="flex flex-row items-center justify-between border-b px-2 py-2">
				<CardTitle class="text-base">{{ t('songRequests.table.title') }}</CardTitle>
			</div>
			<Table class="w-full">
				<TableHeader>
					<TableRow>
						<TableHead class="w-[5%]">#</TableHead>
						<TableHead>{{ t('songRequests.table.columns.title') }}</TableHead>
						<TableHead>{{ t('songRequests.spotify.columns.artist') }}</TableHead>
						<TableHead>{{ t('songRequests.table.columns.added') }}</TableHead>
						<TableHead>{{ t('songRequests.spotify.columns.status') }}</TableHead>
						<TableHead>{{ t('songRequests.table.columns.duration') }}</TableHead>
						<TableHead>{{ t('songRequests.spotify.columns.actions') }}</TableHead>
					</TableRow>
				</TableHeader>
				<TableBody>
					<TableRow
						v-for="(request, index) of queue"
						:key="request.id"
					>
						<TableCell>{{ index + 1 }}</TableCell>
						<TableCell>{{ request.title }}</TableCell>
						<TableCell>{{ request.artist }}</TableCell>
						<TableCell>
							<div class="flex flex-col">
								<span>{{ request.requesterDisplayName || request.requesterName }}</span>
								<span class="text-muted-foreground text-xs">
									{{ formatRelativeTime(request.createdAt) }}
								</span>
							</div>
						</TableCell>
						<TableCell>
							<Badge :variant="statusVariant(request.status)">
								{{ statusLabel(request.status) }}
							</Badge>
						</TableCell>
						<TableCell>{{ convertMillisToTime(request.durationMs) }}</TableCell>
						<TableCell>
							<div class="flex gap-1">
								<Button
									class="min-w-5"
									size="icon"
									variant="secondary"
									:title="t('songRequests.spotify.skip')"
									@click="skipTargetId = request.id"
								>
									<Icon
										name="lucide:skip-forward"
										class="size-5"
									/>
								</Button>
								<Button
									class="min-w-5"
									size="icon"
									variant="destructive"
									:title="t('songRequests.spotify.cancel')"
									@click="cancelTargetId = request.id"
								>
									<Icon
										name="lucide:trash"
										class="size-5"
									/>
								</Button>
							</div>
						</TableCell>
					</TableRow>
					<TableRow v-if="queue.length">
						<TableCell>{{ queue.length }}</TableCell>
						<TableCell></TableCell>
						<TableCell></TableCell>
						<TableCell></TableCell>
						<TableCell></TableCell>
						<TableCell>{{ totalSongsLength }}</TableCell>
						<TableCell></TableCell>
					</TableRow>
					<TableRow v-else>
						<TableCell
							colspan="7"
							class="text-muted-foreground text-center"
						>
							{{ t('songRequests.spotify.emptyQueue') }}
						</TableCell>
					</TableRow>
				</TableBody>
			</Table>
		</CardContent>
	</Card>

	<ActionConfirm
		v-model:open="showSkipConfirm"
		:confirm-text="t('songRequests.spotify.confirmSkip')"
		@confirm="skipRequest"
	/>
	<ActionConfirm
		v-model:open="showCancelConfirm"
		:confirm-text="t('songRequests.spotify.confirmCancel')"
		@confirm="cancelRequest"
	/>
</template>

<script setup lang="ts">
import { computed } from 'vue'

import {
	useProfile,
	useUserAccessFlagChecker,
} from '~~/layers/dashboard/api/auth'
import {
	useChannelOverlayDelete,
	useChannelOverlaysQuery,
} from '~~/layers/dashboard/api/overlays/custom'
import Brb from '~~/layers/dashboard/features/overlays/brb/card.vue'
import Chat from '~~/layers/dashboard/components/overlays/chat.vue'
import Dudes from '~~/layers/dashboard/components/overlays/dudes.vue'
import Kappagen from '~~/layers/dashboard/components/overlays/kappagen.vue'
import NowPlaying from '~~/layers/dashboard/components/overlays/now-playing.vue'
import OBS from '~~/layers/dashboard/features/overlays/obs/card.vue'
import TTS from '~~/layers/dashboard/features/overlays/tts/card.vue'
import { convertOverlayLayerTypeToText } from '~~/layers/dashboard/components/registry/overlays/helpers.ts'
import FaceitStats from '~~/layers/dashboard/features/overlays/faceit-stats/ui/card.vue'
import ValorantStats from '~~/layers/dashboard/features/overlays/valorant-stats/ui/card.vue'
import { ChannelRolePermissionEnum } from '~/gql/graphql.ts'
import { copyToClipBoard } from '~~/layers/dashboard/helpers/index.ts'
import { Button } from '@/components/ui/button'
import { CardContent, CardDescription, CardFooter, CardHeader, CardTitle, Card as ShadCard } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import {
	AlertDialog,
	AlertDialogAction,
	AlertDialogCancel,
	AlertDialogContent,
	AlertDialogDescription,
	AlertDialogFooter,
	AlertDialogHeader,
	AlertDialogTitle,
	AlertDialogTrigger,
} from '@/components/ui/alert-dialog'
import { toast } from 'vue-sonner'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })

const { t } = useI18n()
const userCanManageOverlays = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageOverlays)
const { data: profile } = useProfile()
const selectedDashboardTwitchUser = computed(() => {
	return profile.value?.availableDashboards.find((d) => d.id === profile.value?.selectedDashboardId)
})

async function copyUrl(id: string) {
	await copyToClipBoard(
		`${window.location.origin}/overlays/${selectedDashboardTwitchUser.value?.apiKey}/registry/overlays/${id}`
	)
	toast.success(t('overlays.copied'))
}

const { data: customOverlaysData, executeQuery: refetchOverlays } = useChannelOverlaysQuery()
const customOverlays = computed(() => customOverlaysData.value?.channelOverlays ?? [])
const deleteOverlay = useChannelOverlayDelete()

async function handleDelete(id: string) {
	await deleteOverlay.executeMutation({ id })
	refetchOverlays({ requestPolicy: 'network-only' })
	toast.success(t('sharedTexts.saved'))
}

const router = useRouter()

const maxCustomOverlays = computed(() => {
	const selectedDashboard = profile.value?.availableDashboards.find(
		(d) => d.id === profile.value?.selectedDashboardId
	)
	return selectedDashboard?.plan?.maxCustomOverlays ?? 10
})

const isCreateDisabled = computed(() => {
	return customOverlays.value.length >= maxCustomOverlays.value || !userCanManageOverlays.value
})

function editCustomOverlay(id?: string) {
	if (!id && isCreateDisabled.value) {
		return
	}

	return router.push({
		name: 'registry-overlays-id',
		params: {
			id: id ?? 'new',
		},
	})
}
</script>

<template>
	<div class="flex items-center justify-center max-w-[60vw] mx-auto my-0">
		<div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
			<div>
				<FaceitStats />
			</div>
			<div>
				<ValorantStats />
			</div>
			<div>
				<NowPlaying />
			</div>
			<div>
				<TTS />
			</div>
			<div>
				<OBS />
			</div>
			<div>
				<Chat />
			</div>
			<div>
				<Kappagen />
			</div>
			<div>
				<Dudes />
			</div>
			<div>
				<Brb />
			</div>

			<ShadCard
				v-for="overlay of customOverlays"
				:key="overlay.id"
				class="flex flex-col h-full"
			>
				<CardHeader>
					<CardTitle class="flex items-center justify-between">
						<span>{{ overlay.name }}</span>
						<Badge variant="outline" class="ml-2">
							{{ overlay.layers.length }} {{ overlay.layers.length === 1 ? 'layer' : 'layers' }}
						</Badge>
					</CardTitle>
					<CardDescription v-if="overlay.layers.length">
						{{ t('overlaysRegistry.layers') }}
					</CardDescription>
				</CardHeader>

				<CardContent class="flex-1">
					<div v-if="overlay.layers.length" class="flex gap-2 flex-wrap">
						<Badge
							v-for="layer of overlay.layers"
							:key="layer.id"
							variant="secondary"
						>
							{{ convertOverlayLayerTypeToText(layer.type) }}
						</Badge>
					</div>
					<Alert v-else variant="default" class="border-yellow-500/50">
						<AlertTitle>{{ t('overlaysRegistry.noLayersCreated.title') }}</AlertTitle>
						<AlertDescription>
							{{ t('overlaysRegistry.noLayersCreated.description') }}
						</AlertDescription>
					</Alert>
				</CardContent>

				<CardFooter class="flex gap-2 flex-wrap">
					<Button
						variant="outline"
						size="sm"
						:disabled="!userCanManageOverlays"
						@click="editCustomOverlay(overlay.id)"
					>
						<Icon name="lucide:pencil" class="h-4 w-4 mr-2" />
						<span>{{ t('sharedButtons.settings') }}</span>
					</Button>

					<Button
						variant="outline"
						size="sm"
						:disabled="!userCanManageOverlays"
						@click="copyUrl(overlay.id)"
					>
						<Icon name="lucide:copy" class="h-4 w-4 mr-2" />
						<span>{{ t('overlays.copyOverlayLink') }}</span>
					</Button>

					<AlertDialog>
						<AlertDialogTrigger as-child>
							<Button
								variant="outline"
								size="sm"
								:disabled="!userCanManageOverlays"
								class="text-destructive hover:text-destructive"
							>
								<Icon name="lucide:trash2" class="h-4 w-4 mr-2" />
								<span>{{ t('sharedButtons.delete') }}</span>
							</Button>
						</AlertDialogTrigger>
						<AlertDialogContent>
							<AlertDialogHeader>
								<AlertDialogTitle>{{ t('deleteConfirmation.title') }}</AlertDialogTitle>
								<AlertDialogDescription>
									{{ t('deleteConfirmation.text') }}
								</AlertDialogDescription>
							</AlertDialogHeader>
							<AlertDialogFooter>
								<AlertDialogCancel>{{ t('deleteConfirmation.cancel') }}</AlertDialogCancel>
								<AlertDialogAction @click="() => handleDelete(overlay.id)">
									{{ t('deleteConfirmation.confirm') }}
								</AlertDialogAction>
							</AlertDialogFooter>
						</AlertDialogContent>
					</AlertDialog>
				</CardFooter>
			</ShadCard>

			<ShadCard
				class="flex flex-col h-full cursor-pointer hover:bg-accent/50 transition-colors"
				:class="{ 'cursor-not-allowed opacity-50': isCreateDisabled }"
				@click="() => !isCreateDisabled && editCustomOverlay()"
			>
				<CardContent class="flex-1 flex items-center justify-center p-6">
					<div class="flex flex-col items-center justify-center text-muted-foreground">
						<Icon name="lucide:plus" class="size-16 mb-4" />
						<p class="text-sm font-medium">
							{{ customOverlays.length >= maxCustomOverlays ? t('overlaysRegistry.limitExceeded') : t('overlaysRegistry.createNew') }}
							({{ customOverlays.length }}/{{ maxCustomOverlays }})
						</p>
					</div>
				</CardContent>
			</ShadCard>
		</div>
	</div>
</template>

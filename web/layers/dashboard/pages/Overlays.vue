<script setup lang="ts">
import { Copy, Pencil, Plus, Trash2 } from 'lucide-vue-next'
import { computed } from 'vue'

import { useRouter } from 'vue-router'

import { useUserAccessFlagChecker } from '#layers/dashboard/api/auth'
import {
	useChannelOverlayDelete,
	useChannelOverlaysQuery,
} from '#layers/dashboard/api/overlays/custom'
import Brb from '~/features/overlays/brb/card.vue'
import Chat from '#layers/dashboard/components/overlays/chat.vue'
import Dudes from '#layers/dashboard/components/overlays/dudes.vue'
import Kappagen from '#layers/dashboard/components/overlays/kappagen.vue'
import NowPlaying from '#layers/dashboard/components/overlays/now-playing.vue'
import OBS from '~/features/overlays/obs/card.vue'
import TTS from '~/features/overlays/tts/card.vue'
import { convertOverlayLayerTypeToText } from '#layers/dashboard/components/registry/overlays/helpers.js'
import FaceitStats from '~/features/overlays/faceit-stats/ui/card.vue'
import ValorantStats from '~/features/overlays/valorant-stats/ui/card.vue'
import { ChannelRolePermissionEnum } from '~/gql/graphql'
import { copyToClipBoard } from '#layers/dashboard/helpers/index'





import { toast } from 'vue-sonner'

const { t } = useI18n()
const userCanManageOverlays = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageOverlays)
const { user: profile } = storeToRefs(useDashboardAuth())
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
		name: 'RegistryOverlayEdit',
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

			<!-- Custom Overlays -->
			<UiCard
				v-for="overlay of customOverlays"
				:key="overlay.id"
				class="flex flex-col h-full"
			>
				<UiCardHeader>
					<UiCardTitle class="flex items-center justify-between">
						<span>{{ overlay.name }}</span>
						<UiBadge variant="outline" class="ml-2">
							{{ overlay.layers.length }} {{ overlay.layers.length === 1 ? 'layer' : 'layers' }}
						</UiBadge>
					</UiCardTitle>
					<UiCardDescription v-if="overlay.layers.length">
						{{ t('overlaysRegistry.layers') }}
					</UiCardDescription>
				</UiCardHeader>

				<UiCardContent class="flex-1">
					<div v-if="overlay.layers.length" class="flex gap-2 flex-wrap">
						<UiBadge
							v-for="layer of overlay.layers"
							:key="layer.id"
							variant="secondary"
						>
							{{ convertOverlayLayerTypeToText(layer.type) }}
						</UiBadge>
					</div>
					<UiAlert v-else variant="default" class="border-yellow-500/50">
						<UiAlertTitle>{{ t('overlaysRegistry.noLayersCreated.title') }}</UiAlertTitle>
						<UiAlertDescription>
							{{ t('overlaysRegistry.noLayersCreated.description') }}
						</UiAlertDescription>
					</UiAlert>
				</UiCardContent>

				<UiCardFooter class="flex gap-2 flex-wrap">
					<UiButton
						variant="outline"
						size="sm"
						:disabled="!userCanManageOverlays"
						@click="editCustomOverlay(overlay.id)"
					>
						<Pencil class="h-4 w-4 mr-2" />
						<span>{{ t('sharedButtons.settings') }}</span>
					</UiButton>

					<UiButton
						variant="outline"
						size="sm"
						:disabled="!userCanManageOverlays"
						@click="copyUrl(overlay.id)"
					>
						<Copy class="h-4 w-4 mr-2" />
						<span>{{ t('overlays.copyOverlayLink') }}</span>
					</UiButton>

					<UiAlertDialog>
						<UiAlertDialogTrigger as-child>
							<UiButton
								variant="outline"
								size="sm"
								:disabled="!userCanManageOverlays"
								class="text-destructive hover:text-destructive"
							>
								<Trash2 class="h-4 w-4 mr-2" />
								<span>{{ t('sharedButtons.delete') }}</span>
							</UiButton>
						</UiAlertDialogTrigger>
						<UiAlertDialogContent>
							<UiAlertDialogHeader>
								<UiAlertDialogTitle>{{ t('deleteConfirmation.title') }}</UiAlertDialogTitle>
								<UiAlertDialogDescription>
									{{ t('deleteConfirmation.text') }}
								</UiAlertDialogDescription>
							</UiAlertDialogHeader>
							<UiAlertDialogFooter>
								<UiAlertDialogCancel>{{ t('deleteConfirmation.cancel') }}</UiAlertDialogCancel>
								<UiAlertDialogAction @click="() => handleDelete(overlay.id)">
									{{ t('deleteConfirmation.confirm') }}
								</UiAlertDialogAction>
							</UiAlertDialogFooter>
						</UiAlertDialogContent>
					</UiAlertDialog>
				</UiCardFooter>
			</UiCard>

			<!-- Add New Overlay Card -->
			<UiCard
				class="flex flex-col h-full cursor-pointer hover:bg-accent/50 transition-colors"
				:class="{ 'cursor-not-allowed opacity-50': isCreateDisabled }"
				@click="() => !isCreateDisabled && editCustomOverlay()"
			>
				<UiCardContent class="flex-1 flex items-center justify-center p-6">
					<div class="flex flex-col items-center justify-center text-muted-foreground">
						<Plus class="size-16 mb-4" />
						<p class="text-sm font-medium">
							{{ customOverlays.length >= maxCustomOverlays ? t('overlaysRegistry.limitExceeded') : t('overlaysRegistry.createNew') }}
							({{ customOverlays.length }}/{{ maxCustomOverlays }})
						</p>
					</div>
				</UiCardContent>
			</UiCard>
		</div>
	</div>
</template>

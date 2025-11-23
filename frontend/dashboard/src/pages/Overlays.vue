<script setup lang="ts">
import { IconSettings, IconTrash } from '@tabler/icons-vue'
import { SquarePen } from 'lucide-vue-next'
import { NAlert, NButton, NCard, NPopconfirm, NTag, useNotification } from 'naive-ui'
import { computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

import { useOverlaysRegistry, useProfile, useUserAccessFlagChecker } from '@/api/index.js'
import Card from '@/components/card/card.vue'
import Brb from '@/features/overlays/brb/card.vue'
import Chat from '@/components/overlays/chat.vue'
import Dudes from '@/components/overlays/dudes.vue'
import Kappagen from '@/components/overlays/kappagen.vue'
import NowPlaying from '@/components/overlays/now-playing.vue'
import OBS from '@/components/overlays/obs.vue'
import TTS from '@/features/overlays/tts/card.vue'
import { convertOverlayLayerTypeToText } from '@/components/registry/overlays/helpers.js'
import FaceitStats from '@/features/overlays/faceit-stats/ui/card.vue'
import ValorantStats from '@/features/overlays/valorant-stats/ui/card.vue'
import { ChannelRolePermissionEnum } from '@/gql/graphql'
import { copyToClipBoard } from '@/helpers/index.js'

const { t } = useI18n()
const userCanManageOverlays = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageOverlays)
const { data: profile } = useProfile()
const selectedDashboardTwitchUser = computed(() => {
	return profile.value?.availableDashboards.find((d) => d.id === profile.value?.selectedDashboardId)
})

const message = useNotification()
async function copyUrl(id: string) {
	await copyToClipBoard(
		`${window.location.origin}/overlays/${selectedDashboardTwitchUser.value?.apiKey}/registry/overlays/${id}`
	)
	message.success({
		title: t('overlays.copied'),
		duration: 2500,
	})
}
const overlaysManager = useOverlaysRegistry()
const deleter = overlaysManager.deleteOne
const { data: customOverlays, refetch } = overlaysManager.getAll({})
onMounted(refetch)

const router = useRouter()

function editCustomOverlay(id?: string) {
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

			<div v-for="overlay of customOverlays?.overlays" :key="overlay.id" :span="1">
				<Card :title="overlay.name" style="height: 100%">
					<template #content>
						<div v-if="overlay.layers.length" class="flex gap-1 flex-wrap">
							<NTag v-for="layer of overlay.layers" :key="layer.id" type="success">
								{{ convertOverlayLayerTypeToText(layer.type) }}
							</NTag>
						</div>
						<NAlert v-else type="warning" :title="t('overlaysRegistry.noLayersCreated.title')">
							{{ t('overlaysRegistry.noLayersCreated.description') }}
						</NAlert>
					</template>

					<template #footer>
						<div class="flex gap-2 flex-wrap">
							<NButton
								secondary
								size="large"
								:disabled="!userCanManageOverlays"
								@click="editCustomOverlay(overlay.id)"
							>
								<span>{{ t('sharedButtons.settings') }}</span>
								<IconSettings />
							</NButton>

							<NButton
								secondary
								type="info"
								size="large"
								:disabled="!userCanManageOverlays"
								@click="copyUrl(overlay.id)"
							>
								<span>{{ t('overlays.copyOverlayLink') }}</span>
								<IconSettings />
							</NButton>

							<NPopconfirm
								:positive-text="t('deleteConfirmation.confirm')"
								:negative-text="t('deleteConfirmation.cancel')"
								@positive-click="() => deleter.mutate({ id: overlay.id })"
							>
								{{ t('deleteConfirmation.text') }}
								<template #trigger>
									<NButton secondary size="large" type="error" :disabled="!userCanManageOverlays">
										<span>{{ t('sharedButtons.delete') }}</span>
										<IconTrash />
									</NButton>
								</template>
							</NPopconfirm>
						</div>
					</template>
				</Card>
			</div>

			<div>
				<NCard
					class="h-full"
					:style="{ cursor: userCanManageOverlays ? 'pointer' : 'not-allowed' }"
					embedded
					@click="() => editCustomOverlay()"
				>
					<div class="flex items-center justify-center h-full">
						<SquarePen class="size-16" />
					</div>
				</NCard>
			</div>
		</div>
	</div>
</template>

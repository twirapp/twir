<script setup lang="ts">
import { IconSettings, IconTrash } from '@tabler/icons-vue'
import { SquarePen } from 'lucide-vue-next'
import {
	NAlert,
	NButton,
	NCard,
	NGrid,
	NGridItem,
	NPopconfirm,
	NTag,
	useNotification,
} from 'naive-ui'
import { onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

import { useOverlaysRegistry, useProfile, useUserAccessFlagChecker } from '@/api/index.js'
import Card from '@/components/card/card.vue'
import { responsiveCols } from '@/components/consants.js'
import Brb from '@/components/overlays/brb.vue'
import Chat from '@/components/overlays/chat.vue'
import Dudes from '@/components/overlays/dudes.vue'
import Kappagen from '@/components/overlays/kappagen.vue'
import NowPlaying from '@/components/overlays/now-playing.vue'
import OBS from '@/components/overlays/obs.vue'
import TTS from '@/components/overlays/tts.vue'
import { convertOverlayLayerTypeToText } from '@/components/registry/overlays/helpers.js'
import { ChannelRolePermissionEnum } from '@/gql/graphql'
import { copyToClipBoard } from '@/helpers/index.js'

const { t } = useI18n()
const userCanManageOverlays = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageOverlays)
const { data: profile } = useProfile()

const message = useNotification()
async function copyUrl(id: string) {
	await copyToClipBoard(`${window.location.origin}/overlays/${profile.value?.apiKey}/registry/overlays/${id}`)
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
		<NGrid :cols="responsiveCols" :x-gap="16" :y-gap="16" responsive="screen">
			<NGridItem :span="1">
				<NowPlaying />
			</NGridItem>
			<NGridItem :span="1">
				<TTS />
			</NGridItem>
			<NGridItem :span="1">
				<OBS />
			</NGridItem>
			<NGridItem :span="1">
				<Chat />
			</NGridItem>
			<NGridItem :span="1">
				<Kappagen />
			</NGridItem>
			<NGridItem :span="1">
				<Dudes />
			</NGridItem>
			<NGridItem :span="1">
				<Brb />
			</NGridItem>

			<NGridItem v-for="overlay of customOverlays?.overlays" :key="overlay.id" :span="1">
				<Card
					:title="overlay.name"
					style="height: 100%;"
				>
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
								:disabled="!userCanManageOverlays || profile?.selectedDashboardId !== profile?.id"
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
			</NGridItem>

			<NGridItem :span="1">
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
			</NGridItem>
		</NGrid>
	</div>
</template>

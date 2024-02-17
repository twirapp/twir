<script setup lang="ts">
import { IconPlus, IconSettings, IconTrash } from '@tabler/icons-vue';
import {
	NAlert,
	NButton,
	NCard,
	NGrid,
	NGridItem,
	NPopconfirm,
	NTag,
	useNotification,
} from 'naive-ui';
import { onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';

import { useOverlaysRegistry, useProfile, useUserAccessFlagChecker } from '@/api/index.js';
import Card from '@/components/card/card.vue';
import { responsiveCols } from '@/components/consants.js';
import Brb from '@/components/overlays/brb.vue';
import Chat from '@/components/overlays/chat.vue';
import Dudes from '@/components/overlays/dudes.vue';
import Kappagen from '@/components/overlays/kappagen.vue';
import NowPlaying from '@/components/overlays/now-playing.vue';
import OBS from '@/components/overlays/obs.vue';
import TTS from '@/components/overlays/tts.vue';
import { convertOverlayLayerTypeToText } from '@/components/registry/overlays/helpers.js';
import { copyToClipBoard } from '@/helpers/index.js';


const { t } = useI18n();
const userCanManageOverlays = useUserAccessFlagChecker('MANAGE_OVERLAYS');
const userProfile = useProfile();

const message = useNotification();
const copyUrl = async (id: string) => {
	await copyToClipBoard(`${window.location.origin}/overlays/${userProfile.data.value?.apiKey}/registry/overlays/${id}`);
	message.success({
		title: t('overlays.copied'),
		duration: 2500,
	});
};
const overlaysManager = useOverlaysRegistry();
const deleter = overlaysManager.deleteOne;
const { data: customOverlays, refetch } = overlaysManager.getAll({});
onMounted(refetch);


const router = useRouter();

const editCustomOverlay = (id?: string) => router.push({
	name: 'RegistryOverlayEdit',
	params: {
		id: id ?? 'new',
	},
});
</script>

<template>
	<div
		style="
			display: flex;
			align-items: center;
			justify-content: center;
			max-width: 60vw;
			margin: 0 auto;
		"
	>
		<n-grid :cols="responsiveCols" :x-gap="16" :y-gap="16" responsive="screen">
			<n-grid-item :span="1">
				<now-playing />
			</n-grid-item>
			<n-grid-item :span="1">
				<TTS />
			</n-grid-item>
			<n-grid-item :span="1">
				<OBS />
			</n-grid-item>
			<n-grid-item :span="1">
				<Chat />
			</n-grid-item>
			<n-grid-item :span="1">
				<Kappagen />
			</n-grid-item>
			<n-grid-item :span="1">
				<Dudes />
			</n-grid-item>
			<n-grid-item :span="1">
				<Brb />
			</n-grid-item>

			<n-grid-item v-for="overlay of customOverlays?.overlays" :key="overlay.id" :span="1">
				<card
					:title="overlay.name"
					style="height: 100%;"
				>
					<template #content>
						<div v-if="overlay.layers.length" style="display: flex; gap: 4px; flex-wrap: wrap;">
							<n-tag v-for="layer of overlay.layers" :key="layer.id" type="success">
								{{ convertOverlayLayerTypeToText(layer.type) }}
							</n-tag>
						</div>
						<n-alert v-else type="warning" :title="t('overlaysRegistry.noLayersCreated.title')">
							{{ t('overlaysRegistry.noLayersCreated.description') }}
						</n-alert>
					</template>

					<template #footer>
						<div style="display: flex; gap: 8px; flex-wrap: wrap">
							<n-button
								secondary
								size="large"
								:disabled="!userCanManageOverlays"
								@click="editCustomOverlay(overlay.id)"
							>
								<span>{{ t('sharedButtons.settings') }}</span>
								<IconSettings />
							</n-button>

							<n-button
								secondary
								type="info"
								size="large"
								:disabled="!userCanManageOverlays || userProfile.data.value?.selectedDashboardId != userProfile.data.value?.id"
								@click="copyUrl(overlay.id)"
							>
								<span>{{ t('overlays.copyOverlayLink') }}</span>
								<IconSettings />
							</n-button>

							<n-popconfirm
								:positive-text="t('deleteConfirmation.confirm')"
								:negative-text="t('deleteConfirmation.cancel')"
								@positive-click="() => deleter.mutate({ id: overlay.id })"
							>
								{{ t('deleteConfirmation.text') }}
								<template #trigger>
									<n-button secondary size="large" type="error" :disabled="!userCanManageOverlays">
										<span>{{ t('sharedButtons.delete') }}</span>
										<IconTrash />
									</n-button>
								</template>
							</n-popconfirm>
						</div>
					</template>
				</card>
			</n-grid-item>

			<n-grid-item :span="1">
				<n-card
					content-style="
						display: flex;
						align-items: center;
						justify-content: center;
					"
					:style="{
						cursor: userCanManageOverlays ? 'pointer' : 'not-allowed',
						height: '100%',
					}"
					embedded
					@click="() => editCustomOverlay()"
				>
					<IconPlus style="height: 80px; width: 80px;" />
				</n-card>
			</n-grid-item>
		</n-grid>
	</div>
</template>

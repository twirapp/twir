<script setup lang="ts">
import { IconSettings, IconTrash, IconPlus } from '@tabler/icons-vue';
import { NTag, NGrid, NGridItem, NButton, NPopconfirm, NCard, NAlert } from 'naive-ui';
import { onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';

import { useOverlaysRegistry, useUserAccessFlagChecker } from '@/api';
import { useProfile } from '@/api/index.js';
import Card from '@/components/card/card.vue';
import { convertOverlayLayerTypeToText } from '@/components/registry/overlays/helpers.js';
import { copyToClipBoard } from '@/helpers';

const overlaysManager = useOverlaysRegistry();
const deleter = overlaysManager.deleteOne;
const { data: overlays, refetch } = overlaysManager.getAll({});
onMounted(refetch);

const { t } = useI18n();

const router = useRouter();
const goToEditPage = (id?: string) => router.push({
	name: 'RegistryOverlayEdit',
	params: {
		id: id ?? 'new',
	},
});

const userCanManageOverlays = useUserAccessFlagChecker('MANAGE_OVERLAYS');

const userProfile = useProfile();
const copyUrl = async (id: string) => {
	await copyToClipBoard(`${window.location.origin}/overlays/${userProfile.data.value?.apiKey}/registry/overlays/${id}`);
};
</script>

<template>
	<n-alert type="info">{{ t('overlaysRegistry.description') }}</n-alert>

	<n-grid style="margin-top: 15px;" responsive="screen" cols="1 s:2 m:3 l:3" :x-gap="10" :y-gap="10" item-responsive>
		<n-grid-item :span="1">
			<n-card
				class="c"
				content-style="
						display: flex;
						align-items: center;
						justify-content: center;
					"
				:style="{
					cursor: userCanManageOverlays ? 'pointer' : 'not-allowed',
				}"
				embedded
				@click="() => goToEditPage()"
			>
				<IconPlus style="height: 80px; width: 80px;" />
			</n-card>
		</n-grid-item>

		<n-grid-item v-for="overlay of overlays?.overlays" :key="overlay.id" :span="1">
			<card
				:title="overlay.name"
				class="c"
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
							@click="goToEditPage(overlay.id)"
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
	</n-grid>
</template>

<style scoped>
.c {
	height: 200px;
}
</style>

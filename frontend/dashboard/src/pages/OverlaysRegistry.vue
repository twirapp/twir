<script setup lang="ts">
import { IconSettings, IconTrash } from '@tabler/icons-vue';
import { NTag, NGrid, NGridItem, NButton, NPopconfirm } from 'naive-ui';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';

import { useOverlaysRegistry } from '@/api';
import Card from '@/components/card/card.vue';
import { convertOverlayLayerTypeToText } from '@/components/registry/overlays/helpers.js';

const overlaysManager = useOverlaysRegistry();
const deleter = overlaysManager.deleteOne;
const { data: overlays } = overlaysManager.getAll({});

const { t } = useI18n();

const router = useRouter();
const goToEditPage = (id?: string) => router.push({
	name: 'RegistryOverlayEdit',
	params: {
		id: id ?? 'new',
	},
});
</script>

<template>
	<n-grid responsive="screen" cols="1 s:2 m:2 l:3" :x-gap="10" :y-gap="10">
		<n-grid-item>
			<card
				v-for="overlay of overlays?.overlays"
				:key="overlay.id"
				:title="overlay.name"
			>
				<template #content>
					<div style="display: flex; gap: 4px; flex-wrap: wrap;">
						<n-tag v-for="layer of overlay.layers" :key="layer.id" type="success">
							{{ convertOverlayLayerTypeToText(layer.type) }}
						</n-tag>
					</div>
				</template>

				<template #footer>
					<div style="display: flex; gap: 8px">
						<n-button secondary size="large" @click="goToEditPage(overlay.id)">
							<span>{{ t('sharedButtons.settings') }}</span>
							<IconSettings />
						</n-button>

						<n-popconfirm
							:positive-text="t('deleteConfirmation.confirm')"
							:negative-text="t('deleteConfirmation.cancel')"
							@positive-click="() => deleter.mutate({ id: overlay.id })"
						>
							{{ t('deleteConfirmation.text') }}
							<template #trigger>
								<n-button secondary size="large" type="error">
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

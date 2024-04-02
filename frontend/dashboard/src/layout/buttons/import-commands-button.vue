<script setup lang="ts">
import { IconDownload } from '@tabler/icons-vue';
import { NButton, NA } from 'naive-ui';
import { storeToRefs } from 'pinia';
import { useI18n } from 'vue-i18n';

import { useSidebarCollapseStore } from '../use-sidebar-collapse';

withDefaults(defineProps<{ isDrawer: boolean }>(), {
	isDrawer: false,
});

const { t } = useI18n();
const collapsedStore = useSidebarCollapseStore();
const { isCollapsed } = storeToRefs(collapsedStore);
</script>

<template>
	<div class="flex px-2 mb-2">
		<router-link :to="{ name: 'Import'}" #="{ navigate, href }" custom>
			<n-a :href="href" type="info" secondary block class="w-full" @click="navigate">
				<n-button type="info" secondary block>
					<template #icon>
						<IconDownload />
					</template>
					<template v-if="isDrawer || !isCollapsed">
						{{ t('sidebar.import') }}
					</template>
				</n-button>
			</n-a>
		</router-link>
	</div>
</template>

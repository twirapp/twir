<script setup lang="ts">
import { IconDownload } from '@tabler/icons-vue';
import { NButton, NA } from 'naive-ui';
import { storeToRefs } from 'pinia';
import { useI18n } from 'vue-i18n';

import { useSidebarCollapseStore } from '../use-sidebar-collapse';

const { t } = useI18n();
const collapsedStore = useSidebarCollapseStore();
const { isCollapsed } = storeToRefs(collapsedStore);
</script>

<template>
	<div style="display: flex; padding-left: 8px; padding-right: 8px; margin-bottom: 8px">
		<router-link :to="{ name: 'Import'}" #="{ navigate, href }" custom>
			<n-a :href="href" type="info" secondary block class="w-full" @click="navigate">
				<n-button type="info" secondary block>
					<template #icon>
						<IconDownload />
					</template>
					<template v-if="!isCollapsed">
						{{ t('sidebar.import') }}
					</template>
				</n-button>
			</n-a>
		</router-link>
	</div>
</template>

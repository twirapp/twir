<script setup lang="ts">
import { IconDownload } from '@tabler/icons-vue'
import { NA, NButton } from 'naive-ui'
import { useI18n } from 'vue-i18n'

import { useSidebarCollapseStore } from '../use-sidebar-collapse'

withDefaults(defineProps<{ isDrawer?: boolean }>(), {
	isDrawer: false,
})

const { t } = useI18n()
const { isCollapsed } = useSidebarCollapseStore()
</script>

<template>
	<div class="flex px-2 mb-2">
		<router-link :to="{ name: 'Import' }" #="{ navigate, href }" custom>
			<NA :href="href" type="info" secondary block class="w-full" @click="navigate">
				<NButton type="info" secondary block>
					<template #icon>
						<IconDownload />
					</template>
					<template v-if="isDrawer || !isCollapsed">
						{{ t('sidebar.import') }}
					</template>
				</NButton>
			</NA>
		</router-link>
	</div>
</template>

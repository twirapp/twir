<script setup lang="ts">
import { IconDownload } from '@tabler/icons-vue';
import { NScrollbar, NButton, NA } from 'naive-ui';

import BotStatus from './components/BotStatus.vue';
import Navigation from './navigation.vue';

import FeedbackButton from '@/layout/feedback/feedback-button.vue';
</script>

<template>
	<div class="sidebar">
		<BotStatus />

		<n-scrollbar trigger="none">
			<navigation />
		</n-scrollbar>

		<div style="display: flex; padding-left: 8px; padding-right: 8px; margin-bottom: 8px">
			<router-link :to="{ name: 'Import'}" #="{ navigate, href }" custom>
				<n-a :href="href" type="info" secondary block class="w-full" @click="navigate">
					<n-button type="info" secondary block>
						<template #icon>
							<IconDownload />
						</template>
						Import
					</n-button>
				</n-a>
			</router-link>
		</div>
		<feedback-button />
	</div>
</template>

<style scoped>
.sidebar {
	display: flex;
	flex-direction: column;
	justify-content: space-between;
	height: calc(100vh - var(--layout-header-height));
	width: 100%;
}
</style>

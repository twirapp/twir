<script setup lang="ts">
import {
	NGrid,
	NGridItem,
	NModal,
} from 'naive-ui';
import { computed, ref } from 'vue';

import { useYoutubeModuleSettings } from '@/api/modules/ytsr.js';
import { useYoutubeSocket } from '@/components/songRequests/hook.js';
import Player from '@/components/songRequests/player.vue';
import VideosQueue from '@/components/songRequests/queue.vue';
import SettingsModal from '@/components/songRequests/settings.vue';

const isSettingsModalOpened = ref(false);
const openSettingsModal = () => isSettingsModalOpened.value = true;

useYoutubeSocket();

const youtubeModuleManager = useYoutubeModuleSettings();
const youtubeModuleData = youtubeModuleManager.getAll();

const noCookie = computed(() => {
	return youtubeModuleData.data.value?.data?.playerNoCookieMode ?? false;
});
</script>

<template>
	<n-grid cols="1 s:1 m:1 l:3" responsive="screen" :y-gap="15" :x-gap="15">
		<n-grid-item :span="1">
			<player
				v-if="!youtubeModuleData.isLoading.value"
				:no-cookie="noCookie"
				:open-settings-modal="openSettingsModal"
			/>
		</n-grid-item>

		<n-grid-item :span="2">
			<videos-queue />
		</n-grid-item>
	</n-grid>

	<n-modal
		v-model:show="isSettingsModalOpened"
		:span="10"
		:mask-closable="false"
		:segmented="true"
		preset="card"
		title="Settings"
		:style="{ width: '70%',top: '50px' }"
	>
		<settings-modal />
	</n-modal>
</template>

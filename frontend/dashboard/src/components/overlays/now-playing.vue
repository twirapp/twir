<script setup lang="ts">
import { IconBroadcast } from '@tabler/icons-vue';
import { NModal } from 'naive-ui';
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';

import Card from '@/components/overlays/card.vue';
import NowPlaying from '@/pages/overlays/now-playing/now-playing.vue';

const { t } = useI18n();
const isModalOpened = ref(false);
</script>

<template>
	<card
		:icon="IconBroadcast"
		title="Now Playing"
		:description="t('overlays.nowPlaying.description')"
		:show-copy="false"
		@open-settings="isModalOpened = true"
	>
	</card>

	<n-modal
		v-model:show="isModalOpened"
		:mask-closable="false"
		:segmented="true"
		preset="card"
		title="Now playing"
		content-style="padding: 10px; width: 100%"
		style="width: 500px; max-width: calc(100vw - 40px)"
	>
		<NowPlaying />
	</n-modal>
</template>

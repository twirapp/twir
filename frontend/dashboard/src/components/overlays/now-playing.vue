<script setup lang="ts">
import { IconBroadcast } from '@tabler/icons-vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';

import Card from '@/components/overlays/card.vue';

const { t } = useI18n();

const router = useRouter();
</script>

<template>
	<card
		:icon="IconBroadcast"
		title="Now Playing"
		:description="t('overlays.nowPlaying.description')"
		:show-copy="false"
		@open-settings="router.push({ name: 'NowPlaying' })"
	>
	</card>
</template>

<script setup lang="ts">
import { IconBell } from '@tabler/icons-vue';
import { computed } from 'vue';

import { useProfile } from '@/api/index.js';
import Card from '@/components/overlays/card.vue';

const userProfile = useProfile();
const overlayLink = computed(() => {
	return `${window.location.origin}/overlays/${userProfile.data?.value?.apiKey}/alerts`;
});
</script>

<template>
	<card
		:icon="IconBell"
		title="Alerts"
		:description="`This overlay used to play configured alerts`"
		:overlay-link="overlayLink"
		:show-settings="false"
	>
	</card>
</template>

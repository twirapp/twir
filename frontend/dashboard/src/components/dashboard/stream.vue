<script lang="ts" setup>
import { computed } from 'vue';

import Card from './card.vue';

import { useProfile, useTwitchGetUsers } from '@/api/index.js';
import { storeToRefs } from 'pinia';

const { data: profile } = storeToRefs(useProfile());

const selectedTwitchId = computed(() => profile.value?.selectedDashboardId ?? '');
const selectedDashboardTwitchUser = useTwitchGetUsers({ ids: selectedTwitchId });

const streamUrl = computed(() => {
	if (!selectedDashboardTwitchUser.data.value?.users.length) return;

	const user = selectedDashboardTwitchUser.data.value.users.at(0)!;

	const url = `https://player.twitch.tv/?channel=${user.login}&parent=${window.location.host}&autoplay=false`;

	return url;
});
</script>

<template>
	<card>
		<iframe
			v-if="streamUrl"
			:src="streamUrl"
			width="100%"
			height="100%"
			frameborder="0"
			scrolling="no"
			allowfullscreen="true"
		>
		</iframe>
	</card>
</template>

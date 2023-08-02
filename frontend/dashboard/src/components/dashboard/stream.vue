<script lang="ts" setup>
import { NCard } from 'naive-ui';
import { computed } from 'vue';

import { usePositions } from './positions.js';

import { useProfile, useTwitchGetUsers } from '@/api/index.js';

const { data: profile } = useProfile();

const selectedTwitchId = computed(() => profile.value?.selectedDashboardId ?? '');
const selectedDashboardTwitchUser = useTwitchGetUsers({ ids: selectedTwitchId });

const streamUrl = computed(() => {
	if (!selectedDashboardTwitchUser.data.value?.users.length) return;

	const user = selectedDashboardTwitchUser.data.value.users.at(0)!;

	const url = `https://player.twitch.tv/?channel=${user.login}&parent=${window.location.host}`;

	return url;
});

const positions = usePositions();
</script>

<template>
	<n-card embedded title="Stream" header-style="height: 30px; padding: 5" content-style="padding: 0px;">
		<iframe
			v-if="streamUrl"
			:src="streamUrl"
			width="100%"
			:height="positions.stream.height - 50"
			frameborder="0"
			scrolling="no"
			allowfullscreen="true"
		>
		</iframe>
	</n-card>
</template>

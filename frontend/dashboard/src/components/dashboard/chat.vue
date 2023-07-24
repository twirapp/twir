<script lang="ts" setup>
import { IconSun, IconMoon } from '@tabler/icons-vue';
import { useLocalStorage } from '@vueuse/core';
import { NCard, NButton } from 'naive-ui';
import { computed } from 'vue';

import { usePositions } from './positions.js';

import { useProfile, useTwitchGetUsers } from '@/api/index.js';
import { type Theme } from '@/hooks/index.js';

const { data: profile } = useProfile();


const chatTheme = useLocalStorage<Theme>('twirTwitchChatTheme', 'dark');
const toggleTheme = () => {
	chatTheme.value = chatTheme.value === 'light' ? 'dark' : 'light';
};

const selectedTwitchId = computed(() => profile.value?.selectedDashboardId ?? '');
const selectedDashboardTwitchUser = useTwitchGetUsers({ ids: selectedTwitchId });

const chatUrl = computed(() => {
	if (!selectedDashboardTwitchUser.data.value?.users.length) return;

	const user = selectedDashboardTwitchUser.data.value.users.at(0)!;

	let url = `https://www.twitch.tv/embed/${user.login}/chat?parent=${window.location.host}`;

	if (chatTheme.value === 'dark') {
		url += '&darkpopout';
	}

	return url;
});

const positions = usePositions();
</script>

<template>
	<n-card embedded title="Chat" header-style="height: 30px; padding: 5" content-style="padding: 0px;">
		<template #header-extra>
			<n-button tertiary style="padding: 5px" @click="toggleTheme">
				<IconSun v-if="chatTheme === 'dark'" color="orange" />
				<IconMoon v-else />
			</n-button>
		</template>
		<iframe
			v-if="chatUrl"
			:src="chatUrl"
			width="100%"
			:height="positions.chat.height - 50"
			frameborder="0"
			scrolling="no"
			allowfullscreen="true"
		>
		</iframe>
	</n-card>
</template>

<script lang="ts" setup>
import { IconSun, IconMoon } from '@tabler/icons-vue';
import { useLocalStorage } from '@vueuse/core';
import { NCard, NButton } from 'naive-ui';
import { computed } from 'vue';

import { useProfile } from '@/api/index.js';
import { type Theme } from '@/hooks/index.js';

const { data: profile } = useProfile();


const chatTheme = useLocalStorage<Theme>('twirTwitchChatTheme', 'dark');
const toggleTheme = () => {
	chatTheme.value = chatTheme.value === 'light' ? 'dark' : 'light';
};

const chatUrl = computed(() => {
	if (!profile.value) return;

	let url = `https://www.twitch.tv/embed/${profile.value.login}/chat?parent=${window.location.host}`;

	if (chatTheme.value === 'dark') {
		url += '&darkpopout';
	}

	return url;
});
</script>

<template>
	<n-card embedded title="Chat" style="min-height: 600px" content-style="padding: 0px">
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
			height="100%"
			style="border: 0;"
		>
		</iframe>
	</n-card>
</template>

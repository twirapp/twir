<script lang="ts" setup>
import { IconSun, IconMoon } from '@tabler/icons-vue';
import { NButton } from 'naive-ui';
import { computed } from 'vue';

import Card from './card.vue';


import { useProfile, useTwitchGetUsers } from '@/api/index.js';
import { useTheme } from '@/composables/use-theme.js';

const { data: profile } = useProfile();

const { theme: chatTheme, toggleTheme } = useTheme('twirTwitchChatTheme');

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
</script>

<template>
	<card :content-style="{ 'margin-bottom': '10px', padding: '0px' }">
		<template #header-extra>
			<n-button size="small" text @click="toggleTheme">
				<IconSun v-if="chatTheme === 'dark'" color="orange" />
				<IconMoon v-else />
			</n-button>
		</template>

		<iframe
			v-if="chatUrl"
			:src="chatUrl"
			frameborder="0"
			class="w-full h-full"
		>
		</iframe>
	</card>
</template>

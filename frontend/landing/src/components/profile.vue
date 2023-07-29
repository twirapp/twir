<script lang="ts" setup>
import { Profile } from '@twir/grpc/generated/api/api/auth';
import { computed, onMounted, ref } from 'vue';

import { browserProtectedClient, browserUnProtectedClient } from '../api/twirp-browser.js';
import TwitchIcon from '../assets/login-twitch.svg';

const profile = ref<Profile | null>(null);
const authLink = ref<string | null>(null);

const isLoading = computed(() => authLink.value === null && profile.value === null);

onMounted(() => {
	browserUnProtectedClient.authGetLink({
			state: window.btoa(window.location.origin),
		})
			.then((r) => authLink.value = r.response.link);

		browserProtectedClient.authUserProfile({})
			.then((r) => profile.value = r.response);
});
</script>

<template>
	<a
		class="flex flex-row px-[18px] py-[10px] items-center gap-[8px] bg-[#5D58F5] text-white rounded-lg"
		:href="profile ? '/dashboard' : authLink"
	>
		<template v-if="isLoading">
			<svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
				<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
				<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
			</svg>
			Loading...
		</template>

		<template v-else>
			<template v-if="!profile">
				<span class="font-['Inter'] font-semibold text-white text-base">Login</span>
				<img :src="TwitchIcon" width="18" height="19" />
			</template>
			<template v-else>
				<img :src="profile.avatar" class="rounded-full" width="30" height="30" />
				{profile.displayName}
			</template>
		</template>
	</a>
</template>

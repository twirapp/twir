<script lang="ts" setup>
import { useStore } from '@nanostores/vue';
import { computed, onMounted, ref } from 'vue';

import Loader from '../icons/Loader.svg?component';
import TwitchIcon from '../icons/TwitchLogo.svg?component';
import { profileStore, authLinkStore } from '../stores/auth.js';

const profile = useStore(profileStore);
const authLink = useStore(authLinkStore);

const isLoading = computed(
    () => authLinkStore.value === null && profileStore.value === null && isError.value === false,
);
const isError = ref(false);
const error = ref<string | null>(null);

onMounted(async () => {
  if (typeof window === 'undefined') return;
  const { browserProtectedClient, browserUnProtectedClient } = await import('../api/twirp-browser.js');

  browserUnProtectedClient
		.authGetLink({
			state: window.btoa(window.location.origin),
		})
		.then((r) => authLinkStore.set(r.response.link))
		.catch((err) => {
			isError.value = true;
			console.error(err);
			error.value = String(err);
		});

  browserProtectedClient
		.authUserProfile({})
		.then((r) => profileStore.set(r.response))
		.catch((err) => {
			isError.value = true;
			console.error(err);
			error.value = String(err);
		});
});
</script>

<template>
	<Loader v-if="isLoading" class="w-6 h-6 animate-spin stroke-white/75 stroke-[1.5] m-2" />
	<a
		v-else-if="!profile && !isLoading"
		class="flex flex-row px-4 py-2 items-center gap-2 bg-[#5D58F5] text-white rounded-lg font-medium cursor-pointer hover:bg-[#6964FF]"
		:href="authLink"
	>
		Login
		<TwitchIcon class="w-5 h-5 fill-white" />
	</a>
	<a v-else href="/dashboard" class="text-white/75 inline-flex gap-x-3 items-center">
		{{ profile.displayName }}
		<img :src="profile.avatar" class="rounded-full w-9 h-9 object-cover" />
	</a>
</template>

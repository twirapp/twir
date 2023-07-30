<script lang="ts" setup>
import { Profile } from '@twir/grpc/generated/api/api/auth';
import { computed, ref } from 'vue';

import { browserProtectedClient, browserUnProtectedClient } from '../api/twirp-browser.js';
import ErrorIcon from '../icons/Error.svg?component';
import Loader from '../icons/Loader.svg?component';
import TwitchIcon from '../icons/TwitchLogo.svg?component';

const profile = ref<Profile | null>(null);
const authLink = ref<string | null>(null);

const isLoading = computed(
	() => authLink.value === null && profile.value === null && isError.value === false,
);
const isError = ref(false);
const error = ref<string | null>(null);

browserUnProtectedClient
	.authGetLink({
		state: window.btoa(window.location.origin),
	})
	.then((r) => (authLink.value = r.response.link))
	.catch((err) => {
		isError.value = true;
		console.error(err);
		error.value = String(err);
	});

browserProtectedClient
	.authUserProfile({})
	.then((r) => (profile.value = r.response))
	.catch((err) => {
		isError.value = true;
		console.error(err);
		error.value = String(err);
	});
</script>

<template>
	<Loader v-if="isLoading" class="w-6 h-6 animate-spin stroke-white/75 stroke-[1.5] m-2" />
	<span v-else-if="isError" class="text-[#FD6675] inline-flex items-center gap-x-2 max-w-sm">
		<ErrorIcon class="fill-[#FD6675] h-5 w-5 flex-shrink-0 items-center" />
		<span class="line-clamp-1">{{ error }}</span>
	</span>
	<a v-else-if="profile" href="/dashboard" class="text-white/75 inline-flex gap-x-3 items-center">
		{{ profile.displayName }}
		<img :src="profile.avatar" class="rounded-full w-9 h-9 object-cover" />
	</a>
	<a
		v-else
		class="flex flex-row px-4 py-2 items-center gap-2 bg-[#5D58F5] text-white rounded-lg font-medium cursor-pointer hover:bg-[#6964FF]"
		:href="authLink">
		Login
		<TwitchIcon class="w-5 h-5 fill-white" />
	</a>
</template>

<script setup lang="ts">
import { computed, watch } from 'vue';
import { useRouter } from 'vue-router';

import { useProfile } from '@/api/index.js';

const router = useRouter();

const channelName = computed<string>(() => {
	if (typeof router.currentRoute.value.params.channelName != 'string') {
		return '';
	}
	return router.currentRoute.value.params.channelName;
});

const { data: profile, isLoading } = useProfile(channelName);

watch(profile, (v) => {
	if (!v) return;

	window.document.title = `Twir - ${v.displayName}`;
});
</script>

<template>
	<div
		class="mx-auto p-6 flex items-center bg-neutral-700 text-slate-200 rounded-xl shadow-lg space-x-4"
	>
		<svg
			v-if="isLoading || !profile"
			class="animate-spin w-[50px] text-white" xmlns="http://www.w3.org/2000/svg"
			fill="none"
			viewBox="0 0 24 24"
		>
			<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
			<path
				class="opacity-75" fill="currentColor"
				d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
			></path>
		</svg>

		<div v-else class="flex gap-6">
			<img :src="profile.profileImageUrl" class="rounded-full w-[68px]" />
			<div class="flex flex-col">
				<p class="text-4xl">
					{{ profile?.displayName }}
				</p>
				<a
					:href="'https://twitch.tv/'+profile.login" target="_blank"
					class="text-purple-400"
				>twitch.tv/{{ profile.login }}</a>
			</div>
		</div>
	</div>
</template>

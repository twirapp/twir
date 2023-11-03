<script lang="ts" setup>
import { useStore } from '@nanostores/vue';
import type { Profile } from '@twir/grpc/generated/api/api/auth';

import { protectedClient, unProtectedClient } from '../api/twirp.js';
import TwitchIcon from '../icons/TwitchLogo.svg?component';
import { authLinkStore } from '../stores/auth.js';

const authLink = useStore(authLinkStore);

const props = defineProps<{
	session?: string
	location: string
}>();

let profile: Profile | undefined;

if (props.session) {
	try {
		const request = await protectedClient.authUserProfile({}, {
			meta: { Cookie: `session=${props.session}` },
		});
		profile = request.response;
	// eslint-disable-next-line no-empty
	} catch {}
}

const request = await unProtectedClient.authGetLink({ state: btoa(props.location) });
authLinkStore.set(request.response.link);
</script>

<template>
	<a
		v-if="!profile"
		class="flex flex-row px-4 py-2 items-center gap-2 bg-[#5D58F5] text-white rounded-lg font-medium focus-visible:outline-none focus-visible:ring-4 focus-visible:ring-[#5D58F5]/50 cursor-pointer hover:bg-[#6964FF] transition-shadow"
		:href="authLink"
	>
		Login
		<TwitchIcon class="w-5 h-5 fill-white" />
	</a>
	<a v-else href="/dashboard" class="text-white/75 inline-flex gap-x-3 items-center">
		{{ profile.displayName }}
		<img
			:src="profile.avatar"
			:alt="`${profile.displayName} avatar image`"
			class="rounded-full w-9 h-9 object-cover"
		/>
	</a>
</template>

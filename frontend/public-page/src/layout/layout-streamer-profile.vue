<script setup lang="ts">
import { storeToRefs } from 'pinia';
import { watch } from 'vue';

import { useStreamerProfile, useStreamerPublicSettings } from '@/api/use-streamer-profile';
import { Avatar, AvatarImage, AvatarFallback } from '@/components/ui/avatar';
import { Skeleton } from '@/components/ui/skeleton';

const { data: profile, fetching } = storeToRefs(useStreamerProfile());
const { data: publicSettings } = storeToRefs(useStreamerPublicSettings());

watch(profile, (v) => {
	if (!v) return;

	window.document.title = `Twir - ${v.twitchGetUserByName?.displayName}`;
});
</script>

<template>
	<div class="flex gap-4 rounded-md border p-10">
		<div v-if="fetching || !profile?.twitchGetUserByName" class="flex gap-y-4 flex-wrap space-x-4">
			<Skeleton class="h-32 w-32 rounded-full" />
			<div class="space-y-2">
				<Skeleton class="w-[250px] h-[2.5rem]" />
				<Skeleton class="w-[200px] h-[1.25rem]" />
			</div>
		</div>

		<div v-else class="flex gap-2 justify-between flex-wrap w-full">
			<div class="flex flex-wrap gap-4">
				<Avatar class="w-32 h-32">
					<AvatarImage :src="profile.twitchGetUserByName.profileImageUrl" alt="streamer-profile-image" />
					<AvatarFallback>{{ profile.twitchGetUserByName.login.slice(0, 2) }}</AvatarFallback>
				</Avatar>
				<div class="flex flex-col gap-4">
					<span class="text-4xl break-all">
						{{ profile.twitchGetUserByName.displayName }}
					</span>
					<span class="text-sm break-all">
						{{ publicSettings?.userPublicSettings?.description || profile.twitchGetUserByName.description }}
					</span>
				</div>
			</div>
			<div class="flex md:flex-col sm:flex-row flex-wrap gap-2">
				<a :href="`https://twitch.tv/${profile.twitchGetUserByName.login}`" class="underline" target="_blank">
					Twitch
				</a>
				<a
					v-for="(link, idx) of publicSettings?.userPublicSettings.socialLinks"
					:key="idx"
					:href="link.href"
					class="underline" target="_blank"
				>
					{{ link.title }}
				</a>
			</div>
		</div>
	</div>
</template>

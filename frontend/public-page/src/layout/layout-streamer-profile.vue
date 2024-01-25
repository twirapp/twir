<script setup lang="ts">
import { storeToRefs } from 'pinia';
import { watch } from 'vue';

import { Avatar, AvatarImage, AvatarFallback } from '@/components/ui/avatar';
import { Skeleton } from '@/components/ui/skeleton';
import { useStreamerProfile } from '@/composables/use-streamer-profile';

const emits = defineEmits<{
	updateChannelId: [channelId: string]
}>();

const profileStore = useStreamerProfile();
const { profile, isLoading } = storeToRefs(profileStore);

watch(profile, (v) => {
	if (!v) return;

	window.document.title = `Twir - ${v.displayName}`;
	emits('updateChannelId', v.id);
});
</script>

<template>
	<div class="flex gap-4 rounded-md border p-10">
		<div v-if="isLoading || !profile" class="flex items-center space-x-4">
			<Skeleton class="h-12 w-12 rounded-full" />
			<div class="space-y-2">
				<Skeleton class="h-4 w-[250px]" />
				<Skeleton class="h-4 w-[200px]" />
			</div>
		</div>

		<div v-else class="flex gap-2 justify-between w-full">
			<div class="flex gap-4">
				<Avatar :class="$style.avatar">
					<AvatarImage :src="profile.profileImageUrl" alt="streamer-profile-image" />
					<AvatarFallback>{{ profile.login.slice(0, 2) }}</AvatarFallback>
				</Avatar>
				<div class="flex flex-col gap-4">
					<a class="text-4xl" target="_blank" :href="`https://twitch.tv/${profile.login}`">
						{{ profile.displayName }}
					</a>
					<span class="text-sm">
						Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed vitae
					</span>
				</div>
			</div>
			<div>
				<a :href="`https://twitch.tv/${profile.login}`" class="underline" target="_blank">
					Twitch
				</a>
			</div>
		</div>
	</div>
</template>

<style module>
.avatar {
	width: 8rem;
	height: 8rem;
}
</style>

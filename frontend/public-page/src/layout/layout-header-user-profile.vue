<script setup lang="ts">
import TwitchIcon from '@/assets/icons/socials/twitch.svg?use';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import { Skeleton } from '@/components/ui/skeleton';
import { useUserProfile } from '@/composables/use-user-profile';

const { data, isLoading, isError } = useUserProfile();
</script>

<template>
	<div v-if="isLoading" class="flex items-center space-x-4">
		<Skeleton class="h-12 w-12 rounded-full" />
		<div class="space-y-2">
			<Skeleton class="h-4 w-[250px]" />
			<Skeleton class="h-4 w-[200px]" />
		</div>
	</div>
	<!--	neogate data for testing login button-->
	<div v-else-if="!isError && !data" class="flex items-center gap-2">
		<Avatar>
			<AvatarImage :src="data?.avatar" alt="streamer-profile-image" />
			<AvatarFallback>{{ data?.login.slice(0, 2) }}</AvatarFallback>
		</Avatar>
	</div>
	<div v-else>
		<Button variant="secondary">
			<TwitchIcon class="w-5 h-5 fill-white" />
			Login
		</Button>
	</div>
</template>


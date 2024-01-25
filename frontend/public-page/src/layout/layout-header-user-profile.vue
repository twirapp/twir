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
		<div class="space-y-2">
			<Skeleton class="h-4 w-[50px]" />
			<Skeleton class="h-4 w-[50px]" />
		</div>
		<Skeleton class="h-12 w-12 rounded-full" />
	</div>
	<!--	use !data for test login button -->
	<div v-else-if="!isError && data" class="flex items-center gap-2">
		<div class="flex flex-col">
			<small class="text-xs font-medium leading-none text-muted-foreground">
				Logged as
			</small>
			<span>{{ data.displayName }}</span>
		</div>
		<Avatar>
			<AvatarImage :src="data?.avatar" alt="streamer-profile-image" />
			<AvatarFallback>{{ data?.login.slice(0, 2) }}</AvatarFallback>
		</Avatar>
	</div>
	<div v-else>
		<Button variant="secondary">
			<div class="flex items-center gap-2">
				<span>
					Login
				</span>
				<TwitchIcon class="w-4 h-4 fill-white" />
			</div>
		</Button>
	</div>
</template>


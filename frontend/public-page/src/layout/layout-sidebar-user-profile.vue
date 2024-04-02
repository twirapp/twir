<script setup lang="ts">
import { IconLogout } from '@tabler/icons-vue';

import TwitchIcon from '@/assets/icons/socials/twitch.svg?use';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import { Skeleton } from '@/components/ui/skeleton';
import { useUserProfile, useLogout, useLoginLink } from '@/composables/use-user-profile';

const { data, isLoading, isError } = useUserProfile();
const logout = useLogout();

const { data: loginLink, isError: isLoginLinkError } = useLoginLink();
</script>

<template>
	<div class="w-full">
		<Transition appear mode="out-in">
			<div v-if="isLoading" class="flex items-center gap-2">
				<Skeleton class="h-9 w-9 rounded-full" />
				<div class="space-y-2 w-fit">
					<Skeleton class="h-4 w-[50px]" />
					<Skeleton class="h-4 w-[100px]" />
				</div>
			</div>
			<!--	use !data for test login button -->
			<div v-else-if="!isError && data" class="flex items-center gap-4 justify-between">
				<div class="flex items-center gap-2 max-w-[fit-content] overflow-hidden overflow-ellipsis whitespace-nowrap">
					<Avatar>
						<AvatarImage :src="data.avatar" alt="streamer-profile-image" />
						<AvatarFallback>{{ data.login.slice(0, 2) }}</AvatarFallback>
					</Avatar>
					<div class="flex flex-col">
						<span>{{ data.displayName }}</span>
						<small class="text-xs font-medium leading-none text-muted-foreground">
							Logged as
						</small>
					</div>
				</div>
				<IconLogout class="cursor-pointer" @click="logout.mutate()" />
			</div>
			<div v-else>
				<Button
					variant="secondary" class="w-full" as="a" :href="loginLink"
					:disabled="isLoginLinkError"
				>
					<div class="flex items-center gap-2">
						<span>
							Login
						</span>
						<TwitchIcon class="w-4 h-4 fill-white" />
					</div>
				</Button>
			</div>
		</Transition>
	</div>
</template>

<style scoped>
.v-enter-active,
.v-leave-active {
	transition: opacity 0.5s ease;
}

.v-enter-from,
.v-leave-to {
	opacity: 0;
}
</style>

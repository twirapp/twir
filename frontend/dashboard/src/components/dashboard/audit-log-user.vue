<script setup lang="ts">
import { resolveProfile } from '@/helpers/resolveProfile.js'

interface User {
	profileImageUrl?: string | null
	displayName?: string | null
	login?: string | null
	notFound?: boolean | null
}

const props = defineProps<{
	user: User
	platform?: string | null
}>()

const profile = resolveProfile({
	profileImageUrl: props.user.profileImageUrl,
	displayName: props.user.displayName,
	login: props.user.login,
	platform: props.platform ?? undefined,
	notFound: props.user.notFound,
})
</script>

<template>
	<div class="flex gap-2 items-center shrink-0">
		<div
			v-if="profile.notFound || !profile.url"
			class="flex gap-2 items-center"
		>
			<img
				v-if="profile.avatar"
				class="size-4 rounded-full"
				:src="profile.avatar"
			/>
			<span class="truncate">{{ profile.displayName }}</span>
		</div>
		<a
			v-else
			:href="profile.url"
			target="_blank"
			class="flex gap-2 items-center hover:underline"
		>
			<img
				v-if="profile.avatar"
				class="size-4 rounded-full"
				:src="profile.avatar"
			/>
			<span class="truncate">{{ profile.displayName }}</span>
		</a>
	</div>
</template>

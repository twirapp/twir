<script setup lang="ts">
import { useProfile } from '~/api/user'

const { data: profile } = useProfile()
</script>

<template>
	<a
		v-if="!profile"
		class="flex flex-row px-4 py-2 items-center gap-2 bg-[#5D58F5] text-white rounded-lg font-medium focus-visible:outline-none focus-visible:ring-4 focus-visible:ring-[#5D58F5]/50 cursor-pointer hover:bg-[#6964FF] transition-shadow"
		href=""
	>
		Login
		<SvgoSocialTwitch :fontControlled="false" class="w-5 h-5 fill-white" />
	</a>

	<a v-else href="/dashboard" class="text-white/75 inline-flex gap-x-3 items-center">
		<span class="max-[600px]:hidden">
			{{ profile.displayName }}
		</span>
		<NuxtImg
			:src="profile.avatar"
			alt="avatarImage"
			class="rounded-full w-9 h-9 object-cover"
		/>
	</a>
</template>

<script setup lang="ts">
import type { CommunityUsersOpts } from '~/gql/graphql'

import { useCommunityUsers } from '~/layers/public/api/use-community-users'
import { useStreamerProfile } from '~/layers/public/api/use-streamer-profile'

definePageMeta({
	layout: 'public',
})

const { data: streamerProfile } = await useStreamerProfile()

const communityUsersOpts = computed<CommunityUsersOpts>(() => {
	return {
		channelId: streamerProfile.value?.twitchGetUserByName?.id ?? '',
	}
})

const { data } = await useCommunityUsers(communityUsersOpts)
</script>

<template>
	<div class="flex-wrap w-full border rounded-md" style="background-color: rgb(24, 24, 28)">
		Users
		{{ data }}
	</div>
</template>

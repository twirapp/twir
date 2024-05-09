<script setup lang="ts">
import { computed } from 'vue'

import { useStreamerProfile } from '@/api/use-streamer-profile'
import { useTTSChannelSettings, useTTSUsersSettings } from '@/api/use-tts-settings'
import TableRowsSkeleton from '@/components/TableRowsSkeleton.vue'
import {
	Table,
	TableBody,
	TableCell,
	TableHead,
	TableHeader,
	TableRow,
} from '@/components/ui/table'
import { useTwitchGetUsers } from '@/composables/use-twitch-users'
import UserRow from '@/pages/tts/user-row.vue'

const { data: profile } = useStreamerProfile()
const {
	data: channelSettings,
	isLoading: isChannelSettingsLoading,
} = useTTSChannelSettings()
const {
	data: usersSettings,
	isLoading: isUsersSettingsLoading,
} = useTTSUsersSettings()

const usersIds = computed(() => usersSettings.value?.settings.map(s => s.userId) ?? [])
const { data: users, isLoading: isTwitchUsersLoading } = useTwitchGetUsers(usersIds)

const isLoading = computed(() => {
	return isChannelSettingsLoading.value || isUsersSettingsLoading.value || isTwitchUsersLoading.value
})

const usersWithProfiles = computed(() => {
	if (!users.value?.users) return []

	return users.value.users.map(u => {
		const settings = usersSettings.value?.settings.find(s => s.userId === u.id)
		if (!settings) return null

		return {
			name: u.displayName,
			avatar: u.profileImageUrl,
			...settings,
		}
	}).filter(Boolean)
})
</script>

<template>
	<div class="rounded-md border">
		<Table>
			<TableHeader>
				<TableRow>
					<TableHead class="w-[50px]"></TableHead>
					<TableHead class="w-full">
						User
					</TableHead>
					<TableHead class="w-[100px]">
						Voice
					</TableHead>
					<TableHead class="w-[50px]">
						Rate
					</TableHead>
					<TableHead class="text-right w-[50px]">
						Pitch
					</TableHead>
				</TableRow>
			</TableHeader>
			<Transition name="table-rows" appear mode="out-in">
				<TableBody v-if="!usersWithProfiles || isLoading">
					<TableRowsSkeleton :rows="20" :colspan="5" />
				</TableBody>
				<TableBody v-else-if="!users?.users?.length">
					<TableRow>
						<TableCell :colspan="5">
							<div class="flex items-center justify-center">
								No data
							</div>
						</TableCell>
					</TableRow>
				</TableBody>
				<TableBody v-else>
					<UserRow
						v-if="channelSettings && profile?.twitchGetUserByName"
						:name="profile.twitchGetUserByName.displayName"
						:avatar="profile.twitchGetUserByName.profileImageUrl"
						:pitch="channelSettings.pitch"
						:rate="channelSettings.rate"
						:voice="channelSettings.voice"
					/>

					<UserRow
						v-for="(user) of usersWithProfiles"
						:key="user!.userId"
						:name="user!.name"
						:avatar="user!.avatar"
						:pitch="user!.pitch"
						:rate="user!.rate"
						:voice="user!.voice"
					/>
				</TableBody>
			</Transition>
		</Table>
	</div>
</template>

<style scoped>
.table-rows-enter-active,
.table-rows-leave-active {
	transition: opacity 0.5s ease;
}

.table-rows-enter-from,
.table-rows-leave-to {
	opacity: 0;
}
</style>

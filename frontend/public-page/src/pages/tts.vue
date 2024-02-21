<script setup lang="ts">
import { computed } from 'vue';

import TableRowsSkeleton from '@/components/TableRowsSkeleton.vue';
import {
	Table,
	TableBody,
	TableHead,
	TableHeader,
	TableRow,
} from '@/components/ui/table';
import { useStreamerProfile } from '@/composables/use-streamer-profile';
import { useTTSChannelSettings, useTTSUsersSettings } from '@/composables/use-tts-settings';
import { useTwitchGetUsers } from '@/composables/use-twitch-users';
import UserRow from '@/pages/tts/user-row.vue';


const { data: profile } = useStreamerProfile();
const {
	data: channelSettings,
	isLoading: isChannelSettingsLoading,
} = useTTSChannelSettings();
const {
	data: usersSettings,
	isLoading: isUsersSettingsLoading,
} = useTTSUsersSettings();

const usersIds = computed(() => usersSettings.value?.settings.map(s => s.userId) ?? []);
const { data: users, isLoading: isTwitchUsersLoading } = useTwitchGetUsers(usersIds);

const isLoading = computed(() => {
	return isChannelSettingsLoading.value || isUsersSettingsLoading.value || isTwitchUsersLoading.value;
});

const usersWithProfiles = computed(() => {
	if (!users.value?.users) return [];

	return users.value.users.map(u => {
		const settings = usersSettings.value?.settings.find(s => s.userId === u.id);
		if (!settings) return;

		return {
			name: u.displayName,
			avatar: u.profileImageUrl,
			...settings,
		};
	}).filter(Boolean);
});
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
				<TableBody v-if="isLoading">
					<table-rows-skeleton :rows="20" :colspan="5" />
				</TableBody>
				<TableBody v-else>
					<user-row
						v-if="channelSettings && profile"
						:name="profile?.displayName"
						:avatar="profile?.profileImageUrl"
						:pitch="channelSettings.pitch"
						:rate="channelSettings.rate"
						:voice="channelSettings.voice"
					/>

					<user-row
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

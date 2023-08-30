<script setup lang="ts">
import { computed } from 'vue';
import { useRoute } from 'vue-router';

import {
  useProfile,
  useTTSChannelSettings,
  useTTSUsersSettings,
  useTwitchGetUsers,
} from '@/api/index.js';

const route = useRoute();
const channelName = computed<string>(() => {
  if (typeof route.params.channelName != 'string') {
    return '';
  }
  return route.params.channelName;
});

const { data: profile } = useProfile(channelName);

const channelId = computed<string>(() => {
  if (!profile.value) return '';

  return profile.value.id;
});

const { data: channelSettings } = useTTSChannelSettings(channelId);
const { data: usersSettings } = useTTSUsersSettings(channelId);

const usersIds = computed(() => usersSettings.value?.settings.map(s => s.userId) ?? []);
const { data: users } = useTwitchGetUsers(usersIds);

const usersWithProfiles = computed(() => {
  return users.value?.users.map(u => {
    const settings = usersSettings.value?.settings.find(s => s.userId === u.id);
    if (!settings) return;

    return {
      name: u.displayName,
      avatar: u.profileImageUrl,
      ...settings,
    };
  }).filter(Boolean) ?? [];
});
</script>

<template>
	<div class="overflow-hidden rounded-lg border-gray-200 shadow-lg">
		<table class="w-full border-collapse text-left text-sm text-slate-200 relative">
			<thead class="bg-neutral-700 text-slate-200">
				<tr>
					<th scope="col" class="px-6 py-4 font-medium w-24"></th>
					<th scope="col" class="px-6 py-4 font-medium">
						User
					</th>
					<th scope="col" class="px-6 py-4 font-medium">
						Voice
					</th>
					<th scope="col" class="px-6 py-4 font-medium">
						Rate
					</th>
					<th scope="col" class="px-6 py-4 font-medium">
						Pitch
					</th>
				</tr>
			</thead>
			<tbody class="divide-y divide-neutral-600 border-t border-neutral-600 bg-neutral-700">
				<tr v-if="channelSettings" class="border-b border-purple-400 hover:bg-neutral-600">
					<th class="px-6 py-4">
						<img :src="profile?.profileImageUrl" class="rounded-full w-8 h-8" alt="avatar" />
					</th>
					<th class="px-6 py-4">
						{{ profile?.login }}
					</th>
					<th class="px-6 py-4">
						{{ channelSettings.voice }}
					</th>
					<th class="px-6 py-4">
						{{ channelSettings.rate }}
					</th>
					<th class="px-6 py-4">
						{{ channelSettings.pitch }}
					</th>
				</tr>


				<tr
					v-for="(user, index) of usersWithProfiles" :key="index"
					class="hover:bg-neutral-600"
				>
					<th class="px-6 py-4">
						<img :src="user!.avatar" class="rounded-full w-8 h-8" alt="avatar" />
					</th>
					<th class="px-6 py-4">
						{{ user!.name }}
					</th>
					<th class="px-6 py-4">
						{{ user!.voice }}
					</th>
					<th class="px-6 py-4">
						{{ user!.rate }}
					</th>
					<th class="px-6 py-4">
						{{ user!.pitch }}
					</th>
				</tr>
			</tbody>
		</table>
	</div>
</template>

<style scoped></style>

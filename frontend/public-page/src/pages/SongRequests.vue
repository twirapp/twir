<script setup lang="ts">
import { UseTimeAgo } from '@vueuse/components';
import { computed } from 'vue';
import { useRoute } from 'vue-router';

import { useProfile, useSongsQueue } from '@/api/index.js';
import { convertMillisToTime } from '@/helpers/millisToTime.js';

const route = useRoute();
const channelName = computed<string>(() => {
	if (typeof route.params.channelName != 'string') {
		return '';
	}
	return route.params.channelName;
});

const { data: profile } = useProfile(channelName);

const channelId = computed<string | null>(() => {
	if (!profile.value) return null;

	return profile.value.id;
});
const { data: queue } = useSongsQueue(channelId);
</script>

<template>
	<div class="overflow-hidden rounded-lg border-gray-200 shadow-lg">
		<table class="w-full border-collapse text-left text-sm text-slate-200 relative">
			<thead class="bg-neutral-700 text-slate-200">
				<tr>
					<th scope="col" class="px-6 py-4 font-medium">
						#
					</th>
					<th scope="col" class="px-6 py-4 font-medium">
						Name
					</th>
					<th scope="col" class="px-6 py-4 font-medium">
						Requester
					</th>
					<th scope="col" class="px-6 py-4 font-medium">
						Requested
					</th>
					<th scope="col" class="px-6 py-4 font-medium">
						Duration
					</th>
				</tr>
			</thead>
			<tbody class="divide-y divide-neutral-600 border-t border-neutral-600 bg-neutral-700">
				<tr v-for="(song, index) of queue?.songs" :key="index" class="hover:bg-neutral-600">
					<th class="px-6 py-4 w-1">
						{{ index+1 }}
					</th>
					<th class="px-6 py-4">
						<a href="text-" class="text-purple-200">{{ song.title }}</a>
					</th>
					<th class="px-6 py-4 w-2">
						{{ song.requestedBy }}
					</th>
					<th class="px-6 py-4 w-40">
						<UseTimeAgo v-slot="{ timeAgo }" :time="new Date(Number(song.createdAt))">
							{{ timeAgo }}
						</UseTimeAgo>
					</th>
					<th class="px-6 py-4 w-2">
						{{ convertMillisToTime(song.duration * 1000) }}
					</th>
				</tr>
			</tbody>
		</table>
	</div>
</template>

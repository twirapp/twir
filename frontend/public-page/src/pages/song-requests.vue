<script setup lang="ts">
import { UseTimeAgo } from '@vueuse/components';

import { useSongsQueue } from '@/api/song-requests.js';
import { convertMillisToTime } from '@/helpers.js';

const props = defineProps<{
	channelId: string
	channelName: string
}>();

const { data: queue } = useSongsQueue(props.channelId);
</script>

<template>
	<div class="overflow-auto overflow-y-hidden rounded-lg border-gray-200 shadow-lg">
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
						<a :href="song.songLink" target="_blank" class="text-purple-200">{{ song.title }}</a>
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

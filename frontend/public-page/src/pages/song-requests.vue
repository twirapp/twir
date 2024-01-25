<script setup lang="ts">
import { UseTimeAgo } from '@vueuse/components';

import { convertMillisToTime } from '../helpers';

import TableRowsSkeleton from '@/components/TableRowsSkeleton.vue';
import {
	Table,
	TableBody,
	TableCell,
	TableHead,
	TableHeader,
	TableRow,
} from '@/components/ui/table';
import { useSongsQueue } from '@/composables/use-song-requests';

const { data: queue, isLoading } = useSongsQueue();
</script>

<template>
	<div class="rounded-md border">
		<Table>
			<TableHeader>
				<TableRow>
					<TableHead class="w-[10px]"></TableHead>
					<TableHead class="w-full min-w-[500px]">
						Name
					</TableHead>
					<TableHead class="w-[150px] min-w-[150px]">
						Requested by
					</TableHead>
					<TableHead class="w-[150px] min-w-[150px]"></TableHead>
					<TableHead class="text-right w-[50px]">
						Duration
					</TableHead>
				</TableRow>
			</TableHeader>
			<Transition name="table-rows" appear mode="out-in">
				<TableBody v-if="isLoading">
					<table-rows-skeleton :rows="20" :colspan="5" />
				</TableBody>
				<TableBody v-else>
					<TableRow v-for="(song, idx) in queue?.songs" :key="song.title">
						<TableCell>
							#{{ idx + 1 }}
						</TableCell>
						<TableCell>
							<a :href="song.songLink" target="_blank" class="hover:underline">
								{{ song.title }}
							</a>
						</TableCell>
						<TableCell>
							<a
								:href="`https://twitch.tv/${song.requestedBy}`"
								target="_blank"
								class="hover:underline"
							>
								{{ song.requestedBy }}
							</a>
						</TableCell>
						<TableCell>
							<UseTimeAgo v-slot="{ timeAgo }" :time="new Date(Number(song.createdAt))">
								{{ timeAgo }}
							</UseTimeAgo>
						</TableCell>
						<TableCell class="text-right">
							{{ convertMillisToTime(song.duration * 1000) }}
						</TableCell>
					</TableRow>
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

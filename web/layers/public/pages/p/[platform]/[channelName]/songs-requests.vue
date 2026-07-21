<script setup lang="ts">
import { useSongRequests } from '#layers/public/api/use-song-requests.ts'

import UserCell from '~/components/table/cells/user-cell.vue'
import { convertMillisToTime } from '~/utils/time-utils.ts'

definePageMeta({
	layout: 'public',
})

const { data } = await useSongRequests()
</script>

<template>
	<div class="bg-card w-full flex-wrap rounded-md border">
		<UiTable>
			<UiTableHeader>
				<UiTableRow>
					<UiTableHead class="w-[60%]"> Name </UiTableHead>
					<UiTableHead class="w-[30%]"> Requested by </UiTableHead>
					<UiTableHead class="w-[10%]"> Duration </UiTableHead>
				</UiTableRow>
			</UiTableHeader>
			<UiTableBody>
				<UiTableRow
					v-for="song in data?.songRequestsPublicQueue"
					:key="song.userId"
				>
					<UiTableCell>
						<a
							class="text-white underline"
							:href="song.songLink"
						>
							{{ song.title }}
						</a>
					</UiTableCell>
					<UiTableCell class="font-medium">
						<UserCell
							:name="song.twitchProfile.login"
							:display-name="song.twitchProfile.displayName"
							:avatar="song.twitchProfile.profileImageUrl"
						/>
					</UiTableCell>
					<UiTableCell>
						{{ convertMillisToTime(song.durationSeconds * 1000) }}
					</UiTableCell>
				</UiTableRow>
			</UiTableBody>
		</UiTable>
	</div>
</template>

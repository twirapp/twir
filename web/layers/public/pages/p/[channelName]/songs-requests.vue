<script setup lang="ts">
import UserCell from '~/components/table/cells/user-cell.vue'
import { useSongRequests } from '~/layers/public/api/use-song-requests'
import { convertMillisToTime } from '~/utils/time-utils'

definePageMeta({
	layout: 'public',
})

const { data } = await useSongRequests()
</script>

<template>
	<div class="flex-wrap w-full border rounded-md" style="background-color: rgb(24, 24, 28)">
		<UiTable>
			<UiTableHeader>
				<UiTableRow>
					<UiTableHead class="w-[60%]">
						Name
					</UiTableHead>
					<UiTableHead class="w-[30%]">
						Requested by
					</UiTableHead>
					<UiTableHead class="w-[10%]">
						Duration
					</UiTableHead>
				</UiTableRow>
			</UiTableHeader>
			<UiTableBody>
				<UiTableRow v-for="song in data?.songRequestsPublicQueue" :key="song.userId">
					<UiTableCell>
						<a class="underline text-white" :href="song.songLink">
							{{ song.title }}
						</a>
					</UiTableCell>
					<UiTableCell class="font-medium">
						<UserCell :name="song.twitchProfile.login" :display-name="song.twitchProfile.displayName" :avatar="song.twitchProfile.profileImageUrl" />
					</UiTableCell>
					<UiTableCell>
						{{ convertMillisToTime(song.durationSeconds * 1000) }}
					</UiTableCell>
				</UiTableRow>
			</UiTableBody>
		</UiTable>
	</div>
</template>

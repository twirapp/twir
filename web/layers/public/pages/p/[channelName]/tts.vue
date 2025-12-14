<script setup lang="ts">
import UserCell from '~/components/table/cells/user-cell.vue'
import { useTtsPublicSettings } from '~~/layers/public/api/use-tts-settings'

definePageMeta({
	layout: 'public',
})

const { data } = await useTtsPublicSettings()
</script>

<template>
	<div class="flex-wrap w-full border rounded-md bg-card">
		<UiTable>
			<UiTableHeader>
				<UiTableRow>
					<UiTableHead class="w-[70%]"> User </UiTableHead>
					<UiTableHead class="w-[10%]"> Voice </UiTableHead>
					<UiTableHead class="w-[10%]"> Rate </UiTableHead>
					<UiTableHead class="w-[10%]"> Pitch </UiTableHead>
				</UiTableRow>
			</UiTableHeader>
			<UiTableBody>
				<UiTableRow v-for="setting in data?.ttsPublicUsersSettings" :key="setting.userId">
					<UiTableCell class="font-medium">
						<UserCell
							:name="setting.twitchProfile.login"
							:display-name="setting.twitchProfile.displayName"
							:avatar="setting.twitchProfile.profileImageUrl"
						/>
					</UiTableCell>
					<UiTableCell>
						{{ setting.voice }}
					</UiTableCell>
					<UiTableCell>
						{{ setting.rate }}
					</UiTableCell>
					<UiTableCell>
						{{ setting.pitch }}
					</UiTableCell>
				</UiTableRow>
			</UiTableBody>
		</UiTable>
	</div>
</template>

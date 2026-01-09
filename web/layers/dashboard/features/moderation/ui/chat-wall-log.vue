<script setup lang="ts">
import { ref, watch } from 'vue'


import type { ChatWall } from '#layers/dashboard/api/moderation-chat-wall.ts'

import { useModerationChatWall } from '#layers/dashboard/api/moderation-chat-wall.ts'
import DialogOrSheet from '#layers/dashboard/components/dialog-or-sheet.vue'




const { t } = useI18n()

const props = defineProps<{
	chatWall: ChatWall
}>()

const api = useModerationChatWall()
const { data, executeQuery } = api.useLogs(props.chatWall.id)

const dialogOpened = ref(false)

watch(dialogOpened, (v) => {
	if (v) {
		executeQuery()
	}
})
</script>

<template>
	<UiDialog v-model:open="dialogOpened">
		<UiDialogTrigger as-child>
			<UiButton :disabled="!chatWall.affectedMessages" size="sm">
				{{ t('chatWall.table.affectedMessages') }} ({{ chatWall.affectedMessages }})
			</UiButton>
		</UiDialogTrigger>
		<DialogOrSheet>
			<UiDialogHeader>
				<UiDialogTitle>{{ t('chatWall.table.logs.title') }}</UiDialogTitle>
			</UiDialogHeader>

			<UiTable class="bg-sidebar rounded">
				<UiTableHeader>
					<UiTableRow>
						<UiTableHead class="w-[10%]">
							{{ t('chatWall.table.logs.user') }}
						</UiTableHead>
						<UiTableHead>
							{{ t('chatWall.table.logs.message') }}
						</UiTableHead>
					</UiTableRow>
				</UiTableHeader>
				<UiTableBody>
					<UiTableRow v-for="message of data?.chatWallLogs" :key="message.id">
						<UiTableCell class="w-[10%]">
							<a :href="`https://twitch.tv/${message.twitchProfile.login}`" class="flex items-center gap-2">
								<img :src="message.twitchProfile.profileImageUrl" class="size-6 rounded-full" />
								<span>
									{{ message.twitchProfile.displayName }}
								</span>
							</a>
						</UiTableCell>
						<UiTableCell>
							<span class="wrap-break-word">
								{{ message.text }}
							</span>
						</UiTableCell>
					</UiTableRow>
				</UiTableBody>
			</UiTable>
		</DialogOrSheet>
	</UiDialog>
</template>

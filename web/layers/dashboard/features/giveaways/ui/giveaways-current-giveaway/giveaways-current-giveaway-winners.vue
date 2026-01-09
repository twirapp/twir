<script setup lang="ts">
import { MessageSquareIcon } from 'lucide-vue-next'
import { computed, onMounted, ref, watch } from 'vue'


import type { ChatMessage } from '#layers/dashboard/api/chat-messages.ts'

import { useChatMessagesApi } from '#layers/dashboard/api/chat-messages.ts'
import { useGiveaways } from '~/features/giveaways/composables/giveaways-use-giveaways.ts'
import GiveawaysCurrentGiveawayWinnersProfile from '~/features/giveaways/ui/giveaways-current-giveaway/giveaways-current-giveaway-winners-profile.vue'

const { t } = useI18n()

const { winners } = useGiveaways()

// Chat messages state
const chatMessagesApi = useChatMessagesApi()
const chatMessages = ref<ChatMessage[]>([])
const isLoadingMessages = ref(false)
const selectedWinnerUserId = ref('')

const { executeQuery: refetchMessages } = chatMessagesApi.useQuery(
	{
		perPage: 1000,
	},
	{ manual: true }
)

onMounted(async () => {
	const { data: messages } = await refetchMessages()
	if (!messages?.value?.chatMessages) return
	chatMessages.value = messages.value.chatMessages ?? []
})

const { data: chatMessagesSubscriptionData } = chatMessagesApi.subscribeToChatMessages()
watch(chatMessagesSubscriptionData, (data) => {
	if (data?.chatMessages) {
		chatMessages.value.unshift(data.chatMessages)
	}
})

const filteredMessages = computed(() => {
	return chatMessages.value.filter((msg) => msg.userID === selectedWinnerUserId.value)
})

function handleSelectWinner(winnerId: string) {
	selectedWinnerUserId.value = winnerId
}
</script>

<template>
	<div class="flex flex-col h-full min-h-0">
		<!-- No winners message -->
		<div
			v-if="winners.length === 0"
			class="flex-1 flex items-center justify-center flex-col gap-4 p-4 text-muted-foreground"
		>
			<div class="text-center">
				<p>{{ t('sharedTexts.noData') }}</p>
				<p class="text-sm">
					{{ t('giveaways.chooseWinner') }}
				</p>
			</div>
		</div>

		<!-- Winners content -->
		<template v-else>
			<!-- Winners list -->
			<div class="border-b shrink-0 border-border p-2">
				<div class="flex flex-wrap gap-2">
					<div
						v-for="winner in winners"
						:key="winner.userId"
						class="flex items-center gap-2 p-1 px-2 rounded-md cursor-pointer transition-colors w-full lg:w-auto"
						:class="{
							'bg-muted': winner.userId !== selectedWinnerUserId,
							'bg-primary text-primary-foreground': winner.userId === selectedWinnerUserId,
						}"
						@click="handleSelectWinner(winner.userId)"
					>
						<img
							:src="winner.twitchProfile.profileImageUrl"
							:alt="winner.twitchProfile.displayName"
							class="size-6 rounded-full"
						/>
						<span class="font-medium">{{ winner.twitchProfile.displayName }}</span>
					</div>
				</div>
			</div>

			<!-- Winner's chat messages -->
			<div v-if="selectedWinnerUserId" class="h-0 min-h-0 flex-1 flex flex-col">
				<GiveawaysCurrentGiveawayWinnersProfile :user-id="selectedWinnerUserId" />

				<div class="shrink-0 p-2 border-b border-border">
					<h3 class="text-sm font-medium flex items-center gap-2">
						<MessageSquareIcon class="size-4" />
						Logs
					</h3>
				</div>

				<div class="flex-1 overflow-y-auto h-full">
					<div v-if="isLoadingMessages" class="p-4 text-center text-muted-foreground">
						{{ t('sharedTexts.loading') || 'Loading messages...' }}
					</div>

					<div
						v-else-if="filteredMessages.length === 0"
						class="p-4 text-center text-muted-foreground"
					>
						{{ t('sharedTexts.noData') }}
					</div>

					<div v-else class="p-2 space-y-1">
						<div
							v-for="message in filteredMessages"
							:key="message.id"
							class="py-1 px-2 flex items-start gap-2 hover:bg-muted rounded-sm"
						>
							<span class="text-xs text-muted-foreground whitespace-nowrap shrink-0">{{
								new Date(message.createdAt).toLocaleString()
							}}</span>
							<span class="text-sm wrap-break-word">{{ message.text }}</span>
						</div>
					</div>
				</div>
			</div>

			<div v-else class="flex-1 flex items-center justify-center text-muted-foreground">
				{{
					t('giveaways.currentGiveaway.selectWinner') ||
					'Select a winner to view their chat logs and profile'
				}}
			</div>
		</template>
	</div>
</template>

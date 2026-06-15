<script setup lang="ts">
import type { ChatMessage } from '~~/layers/dashboard/api/chat-messages.js'

import { computed, onMounted, ref, watch } from 'vue'
import { useChatMessagesApi } from '~~/layers/dashboard/api/chat-messages.js'
import { useGiveaways } from '~~/layers/dashboard/features/giveaways/composables/giveaways-use-giveaways.js'
import GiveawaysCurrentGiveawayWinnersProfile from '~~/layers/dashboard/features/giveaways/ui/giveaways-current-giveaway/giveaways-current-giveaway-winners-profile.vue'

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
	<div class="flex h-full min-h-0 flex-col">
		<!-- No winners message -->
		<div
			v-if="winners.length === 0"
			class="text-muted-foreground flex flex-1 flex-col items-center justify-center gap-4 p-4"
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
			<div class="border-border shrink-0 border-b p-2">
				<div class="flex flex-wrap gap-2">
					<div
						v-for="winner in winners"
						:key="winner.userId"
						class="flex w-full cursor-pointer items-center gap-2 rounded-md p-1 px-2 transition-colors lg:w-auto"
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
			<div
				v-if="selectedWinnerUserId"
				class="flex h-0 min-h-0 flex-1 flex-col"
			>
				<GiveawaysCurrentGiveawayWinnersProfile :user-id="selectedWinnerUserId" />

				<div class="border-border shrink-0 border-b p-2">
					<h3 class="flex items-center gap-2 text-sm font-medium">
						<Icon
							name="lucide:message-square"
							class="size-4"
						/>
						Logs
					</h3>
				</div>

				<div class="h-full flex-1 overflow-y-auto">
					<div
						v-if="isLoadingMessages"
						class="text-muted-foreground p-4 text-center"
					>
						{{ t('sharedTexts.loading') || 'Loading messages...' }}
					</div>

					<div
						v-else-if="filteredMessages.length === 0"
						class="text-muted-foreground p-4 text-center"
					>
						{{ t('sharedTexts.noData') }}
					</div>

					<div
						v-else
						class="space-y-1 p-2"
					>
						<div
							v-for="message in filteredMessages"
							:key="message.id"
							class="hover:bg-muted flex items-start gap-2 rounded-sm px-2 py-1"
						>
							<span class="text-muted-foreground shrink-0 text-xs whitespace-nowrap">{{
								new Date(message.createdAt).toLocaleString()
							}}</span>
							<span class="text-sm wrap-break-word">{{ message.text }}</span>
						</div>
					</div>
				</div>
			</div>

			<div
				v-else
				class="text-muted-foreground flex flex-1 items-center justify-center"
			>
				{{
					t('giveaways.currentGiveaway.selectWinner') ||
					'Select a winner to view their chat logs and profile'
				}}
			</div>
		</template>
	</div>
</template>

<script setup lang="ts">
import { MessageSquareIcon } from 'lucide-vue-next'
import { computed, onMounted, ref, watch } from 'vue'

import type { ChatMessage } from '@/api/chat-messages'

import { useChatMessagesApi } from '@/api/chat-messages'
import { ScrollArea } from '@/components/ui/scroll-area'
import { useGiveaways } from '@/features/giveaways/composables/giveaways-use-giveaways.ts'

const { winners } = useGiveaways()

// Chat messages state
const chatMessagesApi = useChatMessagesApi()
const chatMessages = ref<ChatMessage[]>([])
const isLoadingMessages = ref(false)
const selectedWinnerId = ref('')

const { executeQuery: refetchMessages } = chatMessagesApi.useQuery({
	perPage: 1000,
}, { manual: true })

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
	return chatMessages.value.filter(msg => msg.userID === selectedWinnerId.value)
})

function handleSelectWinner(winnerId: string) {
	selectedWinnerId.value = winnerId
}
</script>

<template>
	<div class="flex flex-col h-full">
		<!-- No winners message -->
		<div v-if="winners.length === 0" class="flex-1 flex items-center justify-center flex-col gap-4 p-4 text-muted-foreground">
			<div class="text-center">
				<p>No winners have been chosen yet.</p>
				<p class="text-sm">
					Use the "Choose winners" button to select winners for this giveaway.
				</p>
			</div>
		</div>

		<!-- Winners content -->
		<template v-else>
			<!-- Winners list -->
			<ScrollArea class="border-b border-border p-2">
				<div class="flex flex-wrap gap-2">
					<div
						v-for="winner in winners"
						:key="winner.userId"
						class="flex items-center gap-2 p-2 rounded-md cursor-pointer transition-colors w-full lg:w-auto"
						:class="{
							'bg-muted': winner.userId !== selectedWinnerId,
							'bg-primary text-primary-foreground': winner.userId === selectedWinnerId,
						}"
						@click="handleSelectWinner(winner.userId)"
					>
						<img
							:src="winner.twitchProfile.profileImageUrl"
							:alt="winner.twitchProfile.displayName"
							class="w-8 h-8 rounded-full"
						/>
						<div class="flex flex-col">
							<span class="font-medium">{{ winner.twitchProfile.displayName }}</span>
							<span class="text-xs" :class="{ 'text-muted-foreground': winner.userId !== selectedWinnerId, 'text-primary-foreground/80': winner.userId === selectedWinnerId }">@{{ winner.twitchProfile.login }}</span>
						</div>
					</div>
				</div>
			</ScrollArea>

			<!-- Winner's chat messages -->
			<div v-if="selectedWinnerId" class="flex-1 flex flex-col">
				<div class="p-2 border-b border-border">
					<h3 class="text-sm font-medium flex items-center gap-2">
						<MessageSquareIcon class="size-4" />
						Chat messages
					</h3>
				</div>

				<div class="flex-1 relative overflow-auto">
					<div v-if="isLoadingMessages" class="p-4 text-center text-muted-foreground">
						Loading messages...
					</div>

					<div v-else-if="filteredMessages.length === 0" class="p-4 text-center text-muted-foreground">
						No messages found for this winner
					</div>

					<div v-else class="p-2 space-y-1">
						<div
							v-for="message in filteredMessages"
							:key="message.id"
							class="py-1 px-2 flex items-start gap-2 hover:bg-muted rounded-sm"
						>
							<span class="text-xs text-muted-foreground whitespace-nowrap flex-shrink-0">{{ new Date(message.createdAt).toLocaleString() }}</span>
							<span class="text-sm break-words">{{ message.text }}</span>
						</div>
					</div>
				</div>
			</div>

			<div v-else class="flex-1 flex items-center justify-center text-muted-foreground">
				Select a winner to view their chat messages
			</div>
		</template>
	</div>
</template>

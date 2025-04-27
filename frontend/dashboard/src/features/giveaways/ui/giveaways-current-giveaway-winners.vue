<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { MessageSquareIcon } from 'lucide-vue-next'

import { useGiveaways } from '@/features/giveaways/composables/giveaways-use-giveaways.ts'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { useToast } from '@/components/ui/toast'

// Import chat messages API
import { useChatMessagesApi } from '@/api/chat-messages'

const { winners, currentGiveawayId } = useGiveaways()
const { toast } = useToast()

// Chat messages state
const chatMessagesApi = useChatMessagesApi()
const chatMessages = ref([])
const isLoadingMessages = ref(false)
const selectedWinnerId = ref('')

// Get selected winner's messages
const selectedWinnerMessages = computed(() => {
  if (!selectedWinnerId.value) return []
  return chatMessages.value.filter(msg => msg.userID === selectedWinnerId.value)
})

// Get selected winner
const selectedWinner = computed(() => {
  if (!selectedWinnerId.value) return null
  return winners.value.find(w => w.userId === selectedWinnerId.value)
})

// Load chat messages for a winner
async function loadWinnerMessages(userId) {
  if (!userId) return

  selectedWinnerId.value = userId
  isLoadingMessages.value = true

  try {
    // Fetch chat messages
    const result = await chatMessagesApi.fetchChatMessages({
      userNameLike: selectedWinner.value?.twitchProfile.login,
      perPage: 100
    })

    if (result.data?.chatMessages) {
      chatMessages.value = result.data.chatMessages
    }
  } catch (error) {
    toast({
      variant: 'destructive',
      title: 'Error loading messages',
      description: error instanceof Error ? error.message : 'Unknown error',
    })
  } finally {
    isLoadingMessages.value = false
  }
}

// Subscribe to chat messages
onMounted(() => {
  // Subscribe to chat messages
  const { data: newChatMessages } = chatMessagesApi.subscribeToChatMessages()

  // Watch for new chat messages
  watch(newChatMessages, (newMessage) => {
    if (newMessage?.chatMessages && selectedWinnerId.value) {
      // If the message is from the selected winner, add it to the list
      if (newMessage.chatMessages.userID === selectedWinnerId.value) {
        chatMessages.value = [newMessage.chatMessages, ...chatMessages.value]
      }
    }
  })
})

// Watch for winners changes
watch(winners, (newWinners) => {
  if (newWinners.length > 0 && !selectedWinnerId.value) {
    // Select the first winner by default
    loadWinnerMessages(newWinners[0].userId)
  }
})
</script>

<template>
  <div class="flex flex-col h-full">
    <!-- Winners list -->
    <ScrollArea class="border-b border-border p-2">
      <div class="flex flex-wrap gap-2">
        <div
          v-for="winner in winners"
          :key="winner.userId"
          class="flex items-center gap-2 p-2 rounded-md cursor-pointer transition-colors"
          :class="{
            'bg-muted': winner.userId !== selectedWinnerId,
            'bg-primary text-primary-foreground': winner.userId === selectedWinnerId
          }"
          @click="loadWinnerMessages(winner.userId)"
        >
          <img
            :src="winner.twitchProfile.profileImageUrl"
            :alt="winner.twitchProfile.displayName"
            class="w-8 h-8 rounded-full"
          />
          <div class="flex flex-col">
            <span class="font-medium">{{ winner.twitchProfile.displayName }}</span>
            <span class="text-xs" :class="{'text-muted-foreground': winner.userId !== selectedWinnerId, 'text-primary-foreground/80': winner.userId === selectedWinnerId}">@{{ winner.twitchProfile.login }}</span>
          </div>
        </div>
      </div>
    </ScrollArea>

    <!-- Winner's chat messages -->
    <div v-if="selectedWinnerId" class="flex-1 overflow-hidden flex flex-col">
      <div class="p-2 border-b border-border">
        <h3 class="text-sm font-medium flex items-center gap-2">
          <MessageSquareIcon class="size-4" />
          Chat messages from {{ selectedWinner?.twitchProfile.displayName }}
        </h3>
      </div>

      <ScrollArea class="flex-1">
        <div v-if="isLoadingMessages" class="p-4 text-center text-muted-foreground">
          Loading messages...
        </div>

        <div v-else-if="selectedWinnerMessages.length === 0" class="p-4 text-center text-muted-foreground">
          No messages found for this winner
        </div>

        <div v-else class="p-2 space-y-2">
          <div
            v-for="message in selectedWinnerMessages"
            :key="message.id"
            class="p-2 rounded-md bg-muted"
          >
            <div class="flex items-center gap-2 mb-1">
              <span class="font-medium text-sm">{{ message.userDisplayName }}</span>
              <span class="text-xs text-muted-foreground">{{ new Date(message.createdAt).toLocaleString() }}</span>
            </div>
            <p>{{ message.text }}</p>
          </div>
        </div>
      </ScrollArea>
    </div>

    <div v-else class="flex-1 flex items-center justify-center text-muted-foreground">
      Select a winner to view their chat messages
    </div>
  </div>
</template>

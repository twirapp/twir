<script setup lang="ts">
import { useGiveaways } from '@/features/giveaways/composables/giveaways-use-giveaways.ts'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { ScrollArea } from '@/components/ui/scroll-area'

const { winners } = useGiveaways()
</script>

<template>
  <Card>
    <CardHeader class="pb-2">
      <CardTitle class="text-lg">Winners</CardTitle>
    </CardHeader>
    <CardContent>
      <ScrollArea class="h-24">
        <div class="flex flex-wrap gap-2">
          <div 
            v-for="winner in winners" 
            :key="winner.userId"
            class="flex items-center gap-2 bg-muted p-2 rounded-md"
          >
            <img 
              :src="winner.twitchProfile.profileImageUrl" 
              :alt="winner.twitchProfile.displayName"
              class="w-8 h-8 rounded-full"
            />
            <div class="flex flex-col">
              <span class="font-medium">{{ winner.twitchProfile.displayName }}</span>
              <span class="text-xs text-muted-foreground">@{{ winner.twitchProfile.login }}</span>
            </div>
          </div>
        </div>
      </ScrollArea>
    </CardContent>
  </Card>
</template>

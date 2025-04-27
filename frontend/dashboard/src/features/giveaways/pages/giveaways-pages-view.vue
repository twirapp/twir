<script setup lang="ts">
import { onMounted } from 'vue'
import { useRoute } from 'vue-router'

import { useGiveaways } from '@/features/giveaways/composables/giveaways-use-giveaways.ts'
import GiveawaysCurrentGiveaway from '@/features/giveaways/ui/giveaways-current-giveaway.vue'

const route = useRoute()
const giveawayId = route.params.id as string
const { loadParticipants } = useGiveaways()

onMounted(() => {
  if (giveawayId) {
    loadParticipants(giveawayId)
  }
})
</script>

<template>
  <GiveawaysCurrentGiveaway />
</template>

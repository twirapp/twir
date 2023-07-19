<script setup lang='ts'>
import { useEventsManager } from '@/api/index.js';

const eventsManager = useEventsManager();
const { data: eventsList, isLoading } = eventsManager.getAll({});
</script>

<template>
  <n-card v-for="event of eventsList?.events" :key="event.id">
    {{ event }}
  </n-card>
</template>

<script lang="ts" setup>
import { useStore } from '@nanostores/vue';
import { ModerationUpdateDto } from '@tsuwari/shared';
import { onMounted, ref } from 'vue';

import ModerationComponent from '@/components/Moderation.vue';
import { api } from '@/plugins/api.js';
import { selectedDashboardStore } from '@/stores/userStore';

const settings = ref<ModerationUpdateDto['items']>();

const selectedDashboard = useStore(selectedDashboardStore);

selectedDashboardStore.subscribe(() => {
  getModerationSettings();
});

async function getModerationSettings() {
  const { data } = await api(`/v1/channels/${selectedDashboard.value.channelId}/moderation`);
  settings.value = data;
}

onMounted(() => {
  getModerationSettings();
});
</script>


<template>
  <div 
    class="grid items-start xl:grid-cols-3 lg:grid-cols-2 grid-cols-1 gap-2"
  >
    <div
      v-for="setting, index in settings"
      :key="index"
      class="block rounded-lg card text-white shadow-lg"
    >
      <ModerationComponent :settings="(setting as any)" />
    </div>
  </div>
</template>
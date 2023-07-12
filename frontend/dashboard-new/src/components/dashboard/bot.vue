<script setup lang='ts'>
import { NCard, NAlert, NSkeleton, NButton, NSpin } from 'naive-ui';

import { useBotInfo, useBotJoinPart } from '@/api/index.js';

const { data, isLoading, isFetching } = useBotInfo();

const stateMutation = useBotJoinPart();
async function changeBotState() {
	await stateMutation.mutate(data?.value?.enabled ? 'part' : 'join');
}
</script>

<template>
  <n-card
    title="Bot manage"
    :content-style="{ padding: isLoading ? '10px' : '0px' }"
    :segmented="{
      content: true,
      footer: 'soft'
    }"
  >
    <n-skeleton v-if="!data || isLoading" :sharp="false" />

    <n-alert v-else :type="data?.isMod ? 'success' : 'error'" :bordered="false" class="bot-alert">
      <span v-if="data?.isMod">
        <b>{{ data?.botName }}</b> is a moderator.
      </span>
      <span v-else>
        We have found that the bot is not a moderator on this channel. Please, use <b>/mod {{ data?.botName }}</b>, or some of functionality may work incorrectly.
      </span>
    </n-alert>

    <template #action>
      <n-button
        :type="data?.enabled ? 'error' : 'success'"
        block
        :loading="isFetching || isLoading"
        @click="changeBotState"
      >
        {{ data?.enabled ? 'Leave' : 'Join' }}
      </n-button>
    </template>
  </n-card>
</template>

<style scoped>
.bot-alert {
	border-radius: 0px;
}
</style>

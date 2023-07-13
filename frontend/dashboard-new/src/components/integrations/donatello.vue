<script setup lang='ts'>
import { useTimeout } from '@vueuse/core';
import { NTimeline, NTimelineItem, NText, NInputGroup, NButton, NInput } from 'naive-ui';

import { useDonatelloIntegration } from '@/api/index.js';
import DonatelloSVG from '@/assets/icons/integrations/donatello.svg';
import CopyInput from '@/components/copyInput.vue';
import WithSettings from '@/components/integrations/variants/withSettings.vue';
import { copyToClipBoard } from '@/helpers/index.js';

const { data: donatelloData } = useDonatelloIntegration();
const { start: copyStart, ready: isCopyReady } = useTimeout(2000, { controls: true });

async function copy() {
	console.log(donatelloData?.value?.integrationId);
	if (donatelloData?.value?.integrationId) {
		await copyToClipBoard(donatelloData?.value?.integrationId);
		copyStart();
	}
}

const webhookUrl = `${window.location.origin}/api/webhooks/integrations/donatello`;
</script>

<template>
  <with-settings name="Donatello">
    <template #icon>
      <DonatelloSVG style="width: 50px; margin-top: -15px; margin-bottom: -15px" />
    </template>
    <template #settings>
      <n-timeline>
        <n-timeline-item type="info" title="Step 1">
          <n-text>
            Go to
            <a
              href="https://donatello.to/panel/settings"
              target="_blank"
              class="donatello-link"
            >
              https://donatello.to/panel/settings
            </a>
            and scroll to "Вихідний API" section
          </n-text>
        </n-timeline-item>
        <n-timeline-item type="info" title="Step 2">
          <n-text>Copy api key and paste into "Api Key" input</n-text>
          <copy-input :text="donatelloData?.integrationId ?? ''" style="margin-top: 5px" />
        </n-timeline-item>
        <n-timeline-item type="info" title="Step 3">
          <n-text>Copy link and paste into "Link" field</n-text>
          <copy-input :text="webhookUrl" style="margin-top: 5px" />
        </n-timeline-item>
      </n-timeline>
    </template>
  </with-settings>
</template>

<style scoped>
.donatello-link {
	color: #41c489;
	text-decoration: none
}
</style>

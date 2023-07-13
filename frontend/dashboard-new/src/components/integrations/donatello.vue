<script setup lang='ts'>
import { NTimeline, NTimelineItem, NText } from 'naive-ui';

import { useDonatelloIntegration } from '@/api/index.js';
import DonatelloSVG from '@/assets/icons/integrations/donatello.svg';
import CopyInput from '@/components/copyInput.vue';
import WithSettings from '@/components/integrations/variants/withSettings.vue';

const { data: donatelloData } = useDonatelloIntegration();

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
              class="link"
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
.link {
	color: #41c489;
	text-decoration: none
}
</style>

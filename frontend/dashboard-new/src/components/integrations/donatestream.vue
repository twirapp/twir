<script setup lang='ts'>
import { NTimeline, NTimelineItem, NText } from 'naive-ui';
import { computed } from 'vue';

import { useDonateStreamIntegration } from '@/api/index.js';
import DonateStreamSVG from '@/assets/icons/integrations/donate.stream.svg';
import CopyInput from '@/components/copyInput.vue';
import WithSettings from '@/components/integrations/variants/withSettings.vue';

const integration = useDonateStreamIntegration();
const { data } = integration.useGetData();
const { mutateAsync } = integration.usePost();

const currentPageUrl = `${window.location.origin}/api/webhooks/integrations/donatestream/`;
const webhookUrl = computed(() => {
	return `${currentPageUrl}/${data.value?.integrationId}`;
});
</script>

<template>
  <with-settings name="Donate.stream">
    <template #icon>
      <DonateStreamSVG style="width: 50px" />
    </template>
    <template #settings>
      <n-timeline>
        <n-timeline-item type="info" title="Step 1">
          <n-text>
            Paste that link into input on the
            <a
              href="https://lk.donate.stream/settings/api-key"
              target="_blank"
              class="donatello-link"
            >
              https://lk.donate.stream/settings/api-key
            </a>
            <copy-input :text="webhookUrl" style="margin-top: 5px" />
          </n-text>
        </n-timeline-item>
        <n-timeline-item type="info" title="Step 2">
          <n-text>Copy api key and paste into "Api Key" input</n-text>
          <copy-input :text="''" style="margin-top: 5px" />
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

<script setup lang="ts">
import { NA, NTimeline, NTimelineItem, NText } from 'naive-ui';

import { useDonatelloIntegration } from '@/api/index.js';
import DonatelloSVG from '@/assets/integrations/donatello.svg?use';
import CopyInput from '@/components/copyInput.vue';
import DonateDescription from '@/components/integrations/helpers/donateDescription.vue';
import WithSettings from '@/components/integrations/variants/withSettings.vue';

const { data: donatelloData } = useDonatelloIntegration();

const webhookUrl = `${window.location.origin}/api/webhooks/integrations/donatello`;
</script>

<template>
	<with-settings
		title="Donatello"
		:icon="DonatelloSVG"
		icon-width="48px"
	>
		<template #description>
			<donate-description />
		</template>
		<template #settings>
			<n-timeline>
				<n-timeline-item type="info" title="Step 1">
					<n-text>
						Go to
						<n-a
							href="https://donatello.to/panel/settings"
							target="_blank"
						>
							https://donatello.to/panel/settings
						</n-a>
						and scroll to "Вихідний API" section
					</n-text>
				</n-timeline-item>
				<n-timeline-item type="info" title="Step 2">
					<n-text>Copy api key and paste into "Api Key" input</n-text>
					<copy-input :text="donatelloData?.integrationId ?? ''" class="mt-1" />
				</n-timeline-item>
				<n-timeline-item type="info" title="Step 3">
					<n-text>Copy link and paste into "Link" field</n-text>
					<copy-input :text="webhookUrl" class="mt-1" />
				</n-timeline-item>
			</n-timeline>
		</template>
	</with-settings>
</template>

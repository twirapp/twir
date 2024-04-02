<script setup lang="ts">
import { NTimeline, NTimelineItem, NText, NInput, NButton, NInputGroup, NA } from 'naive-ui';
import { computed, ref } from 'vue';

import { useDonateStreamIntegration } from '@/api/index.js';
import DonateStreamSVG from '@/assets/integrations/donatestream.svg?use';
import CopyInput from '@/components/copyInput.vue';
import DonateDescription from '@/components/integrations/helpers/donateDescription.vue';
import WithSettings from '@/components/integrations/variants/withSettings.vue';

const integration = useDonateStreamIntegration();
const { data } = integration.useGetData();
const { mutateAsync } = integration.usePost();

const currentPageUrl = `${window.location.origin}/api/webhooks/integrations/donatestream`;
const webhookUrl = computed(() => {
	return `${currentPageUrl}/${data.value?.integrationId}`;
});

const secret = ref('');

async function saveSecret() {
	if (!secret.value) return;
	await mutateAsync(secret.value);
}
</script>

<template>
	<with-settings
		title="Donate.stream"
		:icon="DonateStreamSVG"
		icon-width="9rem"
	>
		<template #description>
			<donate-description />
		</template>
		<template #settings>
			<n-timeline>
				<n-timeline-item type="info" title="Step 1">
					<n-text>
						Paste that link into input on the
						<n-a
							href="https://lk.donate.stream/settings/api"
							target="_blank"
						>
							https://lk.donate.stream/settings/api
						</n-a>
						<copy-input :text="webhookUrl" class="mt-1" />
					</n-text>
				</n-timeline-item>
				<n-timeline-item type="info" title="Step 2">
					<n-text>
						Paste the
						<n-a
							href="https://i.imgur.com/OtW97pV.png"
							target="_blank"
						>
							secret key
						</n-a>
						from page and click SAVE
					</n-text>
					<n-input-group>
						<n-input
							v-model:value="secret" type="text" size="small"
							placeholder="secret from page"
						/>
						<n-button size="small" secondary type="success" @click="saveSecret">
							Save
						</n-button>
					</n-input-group>
				</n-timeline-item>
				<n-timeline-item type="info" title="Step 3">
					<n-text>Back to donate.stream and click "confirm" button</n-text>
				</n-timeline-item>
			</n-timeline>
		</template>
	</with-settings>
</template>

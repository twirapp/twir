<script setup lang='ts'>
import { NTimeline, NTimelineItem, NText, NInput, NButton, NInputGroup } from 'naive-ui';
import { computed, ref } from 'vue';
import { useI18n } from 'vue-i18n';

import { useDonateStreamIntegration } from '@/api/index.js';
import DonateStreamSVG from '@/assets/icons/integrations/donate.stream.svg?component';
import CopyInput from '@/components/copyInput.vue';
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

const { t } = useI18n();
</script>

<template>
	<with-settings
		title="Donate.stream"
		:icon="DonateStreamSVG"
		icon-width="100px"
		:description="t('integrations.donateServicesInfo', {
			events: t('sidebar.events').toLocaleLowerCase(),
			chatAlerts: t('sidebar.chatAlerts').toLocaleLowerCase(),
			overlaysRegistry: t('sidebar.overlaysRegistry').toLocaleLowerCase(),
		})"
	>
		<template #settings>
			<n-timeline>
				<n-timeline-item type="info" title="Step 1">
					<n-text>
						Paste that link into input on the
						<a
							href="https://lk.donate.stream/settings/api"
							target="_blank"
							class="link"
						>
							https://lk.donate.stream/settings/api
						</a>
						<copy-input :text="webhookUrl" style="margin-top: 5px" />
					</n-text>
				</n-timeline-item>
				<n-timeline-item type="info" title="Step 2">
					<n-text>
						Paste the <a
							href="https://i.imgur.com/OtW97pV.png"
							target="_blank"
							class="link"
						>
							secret key
						</a>
						from page and click SAVE
					</n-text>
					<n-input-group>
						<n-input v-model:value="secret" type="text" size="small" placeholder="secret from page" />
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

<style scoped>
.link {
	color: #41c489;
	text-decoration: none
}
</style>

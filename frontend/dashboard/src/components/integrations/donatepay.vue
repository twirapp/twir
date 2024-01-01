<script setup lang='ts'>
import { NInputGroup, NButton, NInput, NFormItem } from 'naive-ui';
import { ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import { useDonatepayIntegration } from '@/api/index.js';
import DonatePaySVG from '@/assets/icons/integrations/donatepay.svg?component';
import WithSettings from '@/components/integrations/variants/withSettings.vue';

function redirectToGetApiKey() {
	window.open('https://donatepay.ru/page/api', '_blank');
}

const manager = useDonatepayIntegration();
const { data } = manager.useGetData();
const { mutateAsync } = manager.usePost();

const apiKey = ref<string>('');

watch(data, (value) => {
	if (value?.apiKey) {
		apiKey.value = value.apiKey;
	}
});

async function save() {
	await mutateAsync(apiKey.value);
}

const { t } = useI18n();
</script>

<template>
	<with-settings
		title="Donatepay"
		:save="save"
		:icon="DonatePaySVG"
		icon-width="80px"
		:description="<span v-html='t('integrations.donateServicesInfo', {
			events: t('sidebar.events').toLocaleLowerCase(),
			chatAlerts: t('sidebar.chatAlerts').toLocaleLowerCase(),
			overlaysRegistry: t('sidebar.overlaysRegistry').toLocaleLowerCase(),
		})'>"
	>
		<template #settings>
			<n-form-item label="Api key">
				<n-input-group>
					<n-button secondary type="info" href="qweqwe" @click="redirectToGetApiKey">
						Get api key
					</n-button>
					<n-input
						v-model:value="apiKey"
						type="password"
						placeholder="Api key"
						show-password-on="click"
					/>
				</n-input-group>
			</n-form-item>
		</template>
	</with-settings>
</template>


<script setup lang='ts'>
import { NInput, NFormItem } from 'naive-ui';
import { ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import { useValorantIntegration } from '@/api/index.js';
import ValorantSVG from '@/assets/icons/integrations/valorant.svg?component';
import WithSettings from '@/components/integrations/variants/withSettings.vue';

const manager = useValorantIntegration();
const { data } = manager.useGetData();
const { mutateAsync } = manager.usePost();

const userName = ref<string>('');

watch(data, (value) => {
	if (value?.userName) {
		userName.value = value.userName;
	}
});

async function save() {
	await mutateAsync(userName.value);
}

const { t } = useI18n();
</script>

<template>
	<with-settings
		title="Valorant"
		:save="save"
		:icon="ValorantSVG"
		:description="t('integrations.valorant.info')"
	>
		<template #settings>
			<n-form-item label="Valorant username with tag">
				<n-input
					v-model:value="userName"
					type="text"
					placeholder="username#tag"
				/>
			</n-form-item>
		</template>
	</with-settings>
</template>


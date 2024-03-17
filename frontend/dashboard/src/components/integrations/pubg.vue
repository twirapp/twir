<script setup lang="ts">
import { NFormItem, NInput } from 'naive-ui';
import { ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import { usePubgIntegration } from '@/api/integrations/pubg';
import PubgLogoSVG from '@/assets/integrations/pubg.svg?use';
import withSettings from '@/components/integrations/variants/withSettings.vue';
import { useNaiveDiscrete } from '@/composables/use-naive-discrete';

const { notification } = useNaiveDiscrete();
const { t } = useI18n();


const manager = usePubgIntegration();
const { data } = manager.useGetData();
const updater = manager.usePut();

const nickname = ref('');
async function save() {
	try {
		await updater.mutateAsync(nickname.value);
		notification.success({ title: t('sharedTexts.saved'), duration: 2500 });
	} catch (err) {
		notification.error({ title: t('sharedTexts.errorOnSave'), duration: 2500 });
	}
}

watch(data, () => {
	nickname.value = data.value ?? '';
});

</script>

<template>
	<with-settings title="PUBG" :icon="PubgLogoSVG" icon-width="3rem" :save="save">
		<template #description>
			<i18n-t
				keypath="integrations.pubg.description"
			>
				<b class="variable">$(pubg.fppsolorank)</b>
				<b class="variable">$(pubg.totalwins)</b>
			</i18n-t>
		</template>
		<template #settings>
			<n-form-item label="PUBG nickname">
				<n-input v-model:value="nickname" placeholder="Nickname" />
			</n-form-item>
		</template>
	</with-settings>
</template>

<style scoped>
</style>

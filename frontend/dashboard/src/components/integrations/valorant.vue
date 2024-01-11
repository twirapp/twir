<script setup lang="ts">
import { NInput, NFormItem, useThemeVars } from 'naive-ui';
import { ref, watch } from 'vue';

import { useValorantIntegration } from '@/api/index.js';
import ValorantSVG from '@/assets/integrations/valorant.svg?use';
import WithSettings from '@/components/integrations/variants/withSettings.vue';

const themeVars = useThemeVars();

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
</script>

<template>
	<with-settings
		title="Valorant"
		:save="save"
		:icon="ValorantSVG"
	>
		<template #description>
			<i18n-t
				keypath="integrations.valorant.info"
			>
				<b class="variable">$(valorant.profile.elo)</b>
				<b class="variable">$(valorant.profile.tier)</b>
			</i18n-t>
		</template>
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

<style scoped>
.variable {
	color: v-bind('themeVars.successColor');
}
</style>

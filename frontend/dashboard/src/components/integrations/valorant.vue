<script setup lang='ts'>
import { NInput, NFormItem } from 'naive-ui';
import { ref, watch } from 'vue';

import { useValorantIntegration } from '@/api/index.js';
import ValorantSVG from '@/assets/icons/integrations/valorant.svg';
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
	console.log(userName.value);
	await mutateAsync(userName.value);
}
</script>

<template>
	<with-settings name="Valorant" :save="save">
		<template #icon>
			<ValorantSVG style="width: 50px; display: flex" />
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


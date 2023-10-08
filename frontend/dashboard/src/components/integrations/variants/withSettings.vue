<script setup lang='ts'>
import { IconSettings } from '@tabler/icons-vue';
import { NTooltip, NButton, NModal, NSpace } from 'naive-ui';
import { ref } from 'vue';
import { FunctionalComponent } from 'vue/dist/vue.js';
import { useI18n } from 'vue-i18n';

import { useUserAccessFlagChecker } from '@/api/index.js';

const props = defineProps<{
	name: string,
	save?: () => void | Promise<void>
	description?: string
}>();

defineSlots<{
	icon: FunctionalComponent<any>
	settings: FunctionalComponent<any>
}>();

const showSettings = ref(false);
const modalWidth = '600px';

async function callSave() {
	await props.save?.();
	showSettings.value = false;
}

const userCanManageIntegrations = useUserAccessFlagChecker('MANAGE_INTEGRATIONS');

const { t } = useI18n();
</script>

<template>
	<tr>
		<td>
			<n-tooltip trigger="hover" placement="left">
				<template #trigger>
					<slot name="icon" />
				</template>
				{{ name }}
			</n-tooltip>
		</td>
		<td>
			<span v-if="description">{{ description }}</span>
		</td>
		<td>
			<n-button
				:disabled="!userCanManageIntegrations" strong secondary type="info"
				@click="showSettings = true"
			>
				<IconSettings />
				{{ t('sharedButtons.settings') }}
			</n-button>
		</td>
	</tr>

	<n-modal
		v-model:show="showSettings"
		:mask-closable="false"
		:segmented="true"
		preset="card"
		:title="name"
		class="modal"
		:style="{
			width: modalWidth,
			position: 'fixed',
			left: `calc(50% - ${modalWidth}/2)`,
			top: '50px',
		}"
	>
		<template #header>
			{{ name }}
		</template>
		<slot name="settings" />

		<template #action>
			<n-space justify="end">
				<n-button secondary @click="showSettings = false">
					{{ t('sharedButtons.close') }}
				</n-button>
				<n-button v-if="save" secondary type="success" @click="callSave">
					{{ t('sharedButtons.save') }}
				</n-button>
			</n-space>
		</template>
	</n-modal>
</template>

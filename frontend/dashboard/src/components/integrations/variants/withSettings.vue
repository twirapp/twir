<script setup lang='ts'>
import { IconSettings } from '@tabler/icons-vue';
import { NTooltip, NButton, NModal, NSpace, NSkeleton } from 'naive-ui';
import { onUnmounted, ref } from 'vue';
import { FunctionalComponent } from 'vue/dist/vue.js';
import { useI18n } from 'vue-i18n';

import { useUserAccessFlagChecker } from '@/api/index.js';

const props = withDefaults(defineProps<{
	name: string,
	save?: () => void | Promise<void>,
	description?: string,
	isLoading?: boolean,
	modalWidth?: string,
}>(), {
	modalWidth: '600px',
});

defineSlots<{
	icon: FunctionalComponent<any>
	settings: FunctionalComponent<any>
	content?: FunctionalComponent<any>
}>();

const showSettings = ref(false);

async function callSave() {
	await props.save?.();
}

const userCanManageIntegrations = useUserAccessFlagChecker('MANAGE_INTEGRATIONS');

const { t } = useI18n();

onUnmounted(() => showSettings.value = false);
</script>

<template>
	<tr>
		<td v-if="isLoading" colspan="3">
			<n-skeleton height="40px" width="100%" :sharp="false" />
		</td>
		<template v-else>
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
				<slot name="content" />
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
		</template>
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
			top: '5%',
			bottom: '5%'
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

<script setup lang="ts">
import { IconSettings } from '@tabler/icons-vue';
import { NButton, NModal, NSpace } from 'naive-ui';
import { FunctionalComponent, onUnmounted, ref } from 'vue';
import { useI18n } from 'vue-i18n';

import { useUserAccessFlagChecker } from '@/api';
import Card from '@/components/card/card.vue';

const props = withDefaults(defineProps<{
	title: string
	icon?: FunctionalComponent
	iconFill?: string
	save?: () => void | Promise<void>
	modalWidth?: string,
	iconWidth?: string,
	isLoading?: boolean,
	saveDisabled?: boolean
}>(), {
	modalWidth: '600px',
});

defineSlots<{
	settings: FunctionalComponent,
	customDescriptionSlot?: FunctionalComponent,
	additionalFooter?: FunctionalComponent
	description?: FunctionalComponent | string
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
	<card
		:title="title"
		style="height: 100%;"
		:with-stroke="false"
		:icon="icon"
		:icon-fill="iconFill"
		:icon-width="iconWidth"
		:is-loading="isLoading"
	>
		<template #content>
			<slot name="customDescriptionSlot" />
			<slot class="description" name="description" />
		</template>

		<template #footer>
			<div class="flex justify-between flex-wrap w-full items-center gap-2">
				<div>
					<n-button
						:disabled="!userCanManageIntegrations"
						secondary
						size="large"
						@click="showSettings = true"
					>
						<div class="flex gap-1">
							<span>{{ t('sharedButtons.settings') }}</span>
							<IconSettings />
						</div>
					</n-button>
				</div>

				<div>
					<slot name="additionalFooter" />
				</div>
			</div>
		</template>
	</card>

	<n-modal
		v-model:show="showSettings"
		:mask-closable="false"
		:segmented="true"
		preset="card"
		:title="title"
		class="modal"
		:style="{
			width: modalWidth,
			top: '5%',
			bottom: '5%'
		}"
	>
		<template #header>
			{{ title }}
		</template>
		<slot name="settings" />

		<template #action>
			<n-space justify="end">
				<n-button secondary @click="showSettings = false">
					{{ t('sharedButtons.close') }}
				</n-button>
				<n-button v-if="save" secondary type="success" :disabled="saveDisabled" @click="callSave">
					{{ t('sharedButtons.save') }}
				</n-button>
			</n-space>
		</template>
	</n-modal>
</template>

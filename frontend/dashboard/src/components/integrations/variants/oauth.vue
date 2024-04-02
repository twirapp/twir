<script setup lang="ts">
import { IconLogin, IconLogout, IconSettings } from '@tabler/icons-vue';
import { NButton, NAvatar, useThemeVars, NModal, NSpace } from 'naive-ui';
import { FunctionalComponent, ref } from 'vue';
import { useI18n } from 'vue-i18n';

import { useUserAccessFlagChecker } from '@/api';
import Card from '@/components/card/card.vue';

const themeVars = useThemeVars();

const props = withDefaults(defineProps<{
	title: string,
	isLoading?: boolean
	data: { userName?: string, avatar?: string } | undefined
	logout: () => any
	authLink?: string,
	icon: FunctionalComponent<any>
	iconWidth?: string
	iconColor?: string
	withSettings?: boolean
	save?: () => any | Promise<any>
}>(), {
	authLink: '',
	description: '',
});

defineSlots<{
	settings?: FunctionalComponent,
	description?: FunctionalComponent | string,
}>();

const showSettings = ref(false);

async function login() {
	if (!props.authLink) return;

	window.open(props.authLink, 'Twir connect integration', 'width=800,height=600');
}

async function saveSettings() {
	await props.save?.();
	showSettings.value = false;
}

const userCanManageIntegrations = useUserAccessFlagChecker('MANAGE_INTEGRATIONS');

const { t } = useI18n();
</script>

<template>
	<card
		:title="title"
		style="height: 100%;"
		:with-stroke="false"
		:icon="icon"
		:icon-width="iconWidth"
		:is-loading="isLoading"
	>
		<template #content>
			<slot class="description" name="description" />
		</template>

		<template #footer>
			<div class="flex justify-between flex-wrap items-center gap-1 w-full">
				<div class="flex gap-2 flex-wrap">
					<n-button
						v-if="withSettings"
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
					<n-button
						:disabled="!userCanManageIntegrations || !authLink"
						secondary
						size="large"
						:type="data?.userName ? 'error' : 'success'"
						@click="data?.userName ? logout() : login()"
					>
						<div class="flex gap-1">
							<span>
								{{ t(`sharedButtons.${data?.userName ? 'logout' : 'login'}`) }}
							</span>
							<IconLogout v-if="data?.userName" />
							<IconLogin v-else />
						</div>
					</n-button>
				</div>
				<div
					v-if="data?.userName"
					class="flex gap-2 rounded-[var(--n-border-radius)]"
					:style="{ backgroundColor: themeVars.buttonColor2 }"
				>
					<div class="flex items-center gap-2 h-full px-4 py-2">
						<n-avatar v-if="data?.avatar" round :src="data?.avatar" class="h-6 w-6" />
						<span>{{ data.userName }}</span>
					</div>
				</div>
			</div>
		</template>
	</card>

	<n-modal
		v-if="withSettings"
		v-model:show="showSettings"
		:mask-closable="false"
		:segmented="true"
		preset="card"
		:title="title"
		class="modal"
		:style="{
			width: '50vw',
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
				<n-button v-if="save" secondary type="success" @click="saveSettings">
					{{ t('sharedButtons.save') }}
				</n-button>
			</n-space>
		</template>
	</n-modal>
</template>

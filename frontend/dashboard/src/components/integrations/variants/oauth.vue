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
			<div class="footer">
				<div class="flex gap-2 flex-wrap">
					<n-button
						v-if="withSettings"
						:disabled="!userCanManageIntegrations"
						secondary
						size="large"
						@click="showSettings = true"
					>
						<div style="display: flex; gap: 4px;">
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
						<div class="button-content">
							<span>
								{{ t(`sharedButtons.${data?.userName ? 'logout' : 'login'}`) }}
							</span>
							<IconLogout v-if="data?.userName" />
							<IconLogin v-else />
						</div>
					</n-button>
				</div>
				<div v-if="data?.userName" class="profile">
					<n-avatar v-if="data?.avatar" round :src="data?.avatar" />
					<span>{{ data.userName }}</span>
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

<style scoped>
.footer {
	display: flex;
	justify-content: space-between;
	align-items: center;
	gap: 5px;
	width: 100%;
	flex-wrap: wrap;
}

.button-content {
	display: flex;
	gap: 4px;
}

.profile {
	display: flex;
	align-items: center;
	padding: 10px;
	background-color: v-bind('themeVars.buttonColor2');
	border-radius: 4px;
	gap: 8px;
}

.description :deep(a) {
	color: v-bind('themeVars.successColor');
	text-decoration: none;
}
</style>

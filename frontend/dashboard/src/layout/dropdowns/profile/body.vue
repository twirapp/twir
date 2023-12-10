<script setup lang="ts">
import { NSwitch, NButton, useMessage, NText, NDivider } from 'naive-ui';
import { useI18n } from 'vue-i18n';

import { useProfile, useUser } from '@/api';

const { data: profile } = useProfile();

const userManager = useUser();
const updateUser = userManager.useUpdate();
const regenerateUserApiKey = userManager.useRegenerateApiKey();

const { t } = useI18n();
const message = useMessage();

async function changeLandingVisibilty() {
	if (!profile.value) return;

	await updateUser.mutateAsync({
		hideOnLandingPage: !profile.value.hideOnLandingPage,
	});

	message.success(t('sharedTexts.saved'));
}

async function callRegenerateKey() {
	await regenerateUserApiKey.mutateAsync();

	message.success(t('navbar.profile.regenerateApiKey.info'), {
		duration: 10000,
		closable: true,
	});
}
</script>

<template>
	<div class="profile-body">
		<div style="display: flex; gap: 12px;">
			<n-switch
				:value="!profile?.hideOnLandingPage"
				:disabled="updateUser.isLoading.value"
				@change="changeLandingVisibilty"
			/>
			<n-text>{{ t('navbar.profile.showMeOnLanding') }}</n-text>
		</div>

		<n-divider style="margin: 0" />

		<div style="display: flex; flex-direction: column; gap: 4px">
			<n-button
				:disabled="regenerateUserApiKey.isLoading.value"
				secondary
				type="warning"
				@click="callRegenerateKey"
			>
				{{ t('navbar.profile.regenerateApiKey.button') }}
			</n-button>
			<n-text style="font-size: 11px;" depth="3">
				{{ t('navbar.profile.regenerateApiKey.info') }}
			</n-text>
		</div>
	</div>
</template>

<style scoped>
.profile-body {
	display: flex;
	flex-direction: column;
	padding: 8px;
	gap: 12px;
}
</style>

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
	<div class="flex flex-col p-2 gap-3">
		<div class="flex gap-3">
			<n-switch
				:value="!profile?.hideOnLandingPage"
				:disabled="updateUser.isLoading.value"
				@change="changeLandingVisibilty"
			/>
			<n-text>{{ t('navbar.profile.showMeOnLanding') }}</n-text>
		</div>

		<n-divider class="m-0" />

		<div class="flex flex-col gap-1">
			<n-button
				:disabled="regenerateUserApiKey.isLoading.value"
				secondary
				type="warning"
				@click="callRegenerateKey"
			>
				{{ t('navbar.profile.regenerateApiKey.button') }}
			</n-button>
			<n-text class="text-xs" depth="3">
				{{ t('navbar.profile.regenerateApiKey.info') }}
			</n-text>
		</div>
	</div>
</template>

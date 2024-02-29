<script setup lang="ts">
import { IconRefresh } from '@tabler/icons-vue';
import { NButton, NCard, NFormItem, NInput, NSwitch, NText, useMessage } from 'naive-ui';
import { useI18n } from 'vue-i18n';

import { useProfile, useUser } from '@/api';

const { data: profile } = useProfile();

const userManager = useUser();
const updateUser = userManager.useUpdate();
const regenerateUserApiKey = userManager.useRegenerateApiKey();

const { t } = useI18n();
const message = useMessage();

async function changeLandingVisibility() {
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
	<div class="account-body">
		<n-card title="Private" size="small" bordered>
			<div style="display: flex; gap: 12px;">
				<n-text>{{ t('navbar.profile.showMeOnLanding') }}</n-text>
				<n-switch
					:value="!profile?.hideOnLandingPage"
					:disabled="updateUser.isLoading.value"
					@change="changeLandingVisibility"
				/>
			</div>
		</n-card>

		<n-card title="Api" size="small" bordered>
			<n-form-item label="Key">
				<div class="flex gap-1 w-full flex-wrap">
					<n-input
						type="password"
						show-password-on="click"
						:value="profile?.apiKey"
						:maxlength="8"
						class="flex-1"
					/>
					<n-button secondary type="warning" class="min-w-[150px] sm:w-full" @click="callRegenerateKey">
						<div class="flex items-center gap-1">
							<IconRefresh style="height: 20px; width: 20px" />
							{{ t('navbar.profile.regenerateApiKey.button') }}
						</div>
					</n-button>
				</div>
			</n-form-item>
			<n-text style="font-size: 13px;" depth="3">
				{{ t('navbar.profile.regenerateApiKey.info') }}
			</n-text>
		</n-card>
	</div>
</template>

<style scoped>
.account-body {
	display: flex;
	flex-direction: column;
	padding: 8px;
	gap: 12px;
}

</style>

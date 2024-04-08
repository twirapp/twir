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
	<div class="flex flex-col gap-3">
		<n-card title="Private" size="small" bordered>
			<div class="flex gap-3">
				<n-text>{{ t('userSettings.account.showMeOnLanding') }}</n-text>
				<n-switch
					:value="!profile?.hideOnLandingPage"
					:disabled="updateUser.isLoading.value"
					@update-value="changeLandingVisibility"
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
							<IconRefresh class="h-5 w-5" />
							{{ t('userSettings.account.regenerateApiKey.button') }}
						</div>
					</n-button>
				</div>
			</n-form-item>
			<n-text class="text-sx" depth="3">
				{{ t('userSettings.account.regenerateApiKey.info') }}
			</n-text>
		</n-card>
	</div>
</template>

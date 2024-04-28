<script setup lang="ts">
import { IconRefresh } from '@tabler/icons-vue';
import { NButton, NCard, NFormItem, NInput, NSwitch, NText } from 'naive-ui';
import { storeToRefs } from 'pinia';
import { useI18n } from 'vue-i18n';

import { useProfile, useUserSettings } from '@/api';
import { useToast } from '@/components/ui/toast';

const { data: profile } = storeToRefs(useProfile());

const userManager = useUserSettings();
const updateUser = userManager.useUserUpdateMutation();
const regenerateUserApiKey = userManager.useApiKeyGenerateMutation();

const { t } = useI18n();
const toast = useToast();

async function changeLandingVisibility() {
	if (!profile.value) return;

	await updateUser.executeMutation({
		opts: {
			hideOnLandingPage: !profile.value.hideOnLandingPage,
		},
	});

	toast.toast({
		title: t('sharedTexts.saved'),
		duration: 1500,
		variant: 'success',
	});
}

async function callRegenerateKey() {
	await regenerateUserApiKey.executeMutation({});

	toast.toast({
		title: t('sharedTexts.saved'),
		variant: 'success',
	});
}
</script>

<template>
	<div class="flex flex-col gap-6">
		<div class="flex flex-col gap-6">
			<h4 class="scroll-m-20 text-xl font-semibold tracking-tight">
				Private
			</h4>
			<n-card size="small" bordered>
				<div class="flex gap-3">
					<n-text>{{ t('userSettings.account.showMeOnLanding') }}</n-text>
					<n-switch
						:value="!profile?.hideOnLandingPage"
						:disabled="updateUser.fetching.value"
						@update-value="changeLandingVisibility"
					/>
				</div>
			</n-card>
		</div>

		<div class="flex flex-col gap-6">
			<h4 class="scroll-m-20 text-xl font-semibold tracking-tight">
				Api
			</h4>

			<n-card size="small" bordered>
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
	</div>
</template>

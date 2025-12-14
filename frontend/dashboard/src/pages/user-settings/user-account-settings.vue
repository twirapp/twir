<script setup lang="ts">
import { IconRefresh } from '@tabler/icons-vue'
import { NButton, NCard, NFormItem, NInput, NSwitch, NText } from 'naive-ui'
import { useI18n } from 'vue-i18n'

import { useProfile, useUserSettings } from '@/api'
import { toast } from 'vue-sonner'

const { data: profile } = useProfile()

const userManager = useUserSettings()
const updateUser = userManager.useUserUpdateMutation()
const regenerateUserApiKey = userManager.useApiKeyGenerateMutation()

const { t } = useI18n()

async function changeLandingVisibility() {
	if (!profile.value) return

	await updateUser.executeMutation({
		opts: {
			hideOnLandingPage: !profile.value.hideOnLandingPage,
		},
	})

	toast.success(t('sharedTexts.saved'), {
		duration: 1500,
	})
}

async function callRegenerateKey() {
	await regenerateUserApiKey.executeMutation({})

	toast.success(t('sharedTexts.saved'))
}
</script>

<template>
	<div class="flex flex-col gap-6">
		<div class="flex flex-col gap-6">
			<h4 class="scroll-m-20 text-xl font-semibold tracking-tight">Private</h4>
			<NCard size="small" bordered>
				<div class="flex gap-3">
					<NText>{{ t('userSettings.account.showMeOnLanding') }}</NText>
					<NSwitch
						:value="!profile?.hideOnLandingPage"
						:disabled="updateUser.fetching.value"
						@update-value="changeLandingVisibility"
					/>
				</div>
			</NCard>
		</div>

		<div class="flex flex-col gap-6">
			<h4 class="scroll-m-20 text-xl font-semibold tracking-tight">Api</h4>

			<NCard size="small" bordered>
				<NFormItem label="Key">
					<div class="flex gap-1 w-full flex-wrap">
						<NInput
							type="password"
							show-password-on="click"
							:value="profile?.apiKey"
							:maxlength="8"
							class="flex-1"
						/>
						<NButton
							secondary
							type="warning"
							class="min-w-[150px] sm:w-full"
							@click="callRegenerateKey"
						>
							<div class="flex items-center gap-1">
								<IconRefresh class="h-5 w-5" />
								{{ t('userSettings.account.regenerateApiKey.button') }}
							</div>
						</NButton>
					</div>
				</NFormItem>
				<NText class="text-sx" depth="3">
					{{ t('userSettings.account.regenerateApiKey.info') }}
				</NText>
			</NCard>
		</div>
	</div>
</template>

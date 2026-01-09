<script setup lang="ts">
import { CopyIcon, EyeIcon, EyeOffIcon, RefreshCwIcon } from 'lucide-vue-next'
import { ref } from 'vue'


import { useUserSettings } from '#layers/dashboard/api/auth'





import { toast } from 'vue-sonner'

const { user: profile } = storeToRefs(useDashboardAuth())

const userManager = useUserSettings()
const updateUser = userManager.useUserUpdateMutation()
const regenerateUserApiKey = userManager.useApiKeyGenerateMutation()

const { t } = useI18n()

const showApiKey = ref(false)

async function changeLandingVisibility() {
	if (!profile.value) return

	await updateUser.executeMutation({
		opts: {
			hideOnLandingPage: !profile.value.hideOnLandingPage,
		},
	})

	await executeQuery({ requestPolicy: 'network-only' })

	toast.success(t('sharedTexts.saved'), {
		duration: 1500,
	})
}

async function callRegenerateKey() {
	const result = await regenerateUserApiKey.executeMutation({})

	if (result.error) {
		toast.error('Failed to regenerate API key')
		return
	}

	await executeQuery({ requestPolicy: 'network-only' })

	toast.success(t('sharedTexts.saved'))
}

async function copyApiKey() {
	if (!profile.value?.apiKey) return

	try {
		await navigator.clipboard.writeText(profile.value.apiKey)
		toast.success(t('sharedTexts.copied'), {
			duration: 1500,
		})
	} catch (err) {
		toast.error('Failed to copy')
	}
}
</script>

<template>
	<div class="flex flex-col gap-6">
		<div class="flex flex-col gap-6">
			<h4 class="scroll-m-20 text-xl font-semibold tracking-tight">Private</h4>
			<UiCard class="p-2">
				<UiCardContent>
					<div class="flex items-center gap-2">
						<UiSwitch
							:model-value="!profile?.hideOnLandingPage"
							@update:model-value="changeLandingVisibility"
						/>
						<UiLabel>{{ t('userSettings.account.showMeOnLanding') }}</UiLabel>
					</div>
				</UiCardContent>
			</UiCard>
		</div>

		<div class="flex flex-col gap-6">
			<h4 class="scroll-m-20 text-xl font-semibold tracking-tight">Api</h4>

			<UiCard>
				<UiCardHeader>
					<UiCardTitle>Key</UiCardTitle>
				</UiCardHeader>
				<UiCardContent class="space-y-4">
					<div class="flex gap-2 w-full flex-wrap">
						<UiInput
							:type="showApiKey ? 'text' : 'password'"
							:model-value="profile?.apiKey ?? ''"
							class="flex-1"
							readonly
						/>
						<UiButton variant="outline" size="icon" type="button" @click="showApiKey = !showApiKey">
							<EyeIcon v-if="!showApiKey" />
							<EyeOffIcon v-else />
						</UiButton>
						<UiButton variant="outline" size="icon" type="button" @click="copyApiKey">
							<CopyIcon />
						</UiButton>
						<UiButton variant="outline" class="min-w-37.5 sm:w-full" @click="callRegenerateKey">
							<RefreshCwIcon />
							{{ t('userSettings.account.regenerateApiKey.button') }}
						</UiButton>
					</div>
					<p class="text-sm text-muted-foreground">
						{{ t('userSettings.account.regenerateApiKey.info') }}
					</p>
				</UiCardContent>
			</UiCard>
		</div>
	</div>
</template>

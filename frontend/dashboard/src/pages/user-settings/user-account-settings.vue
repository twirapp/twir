<script setup lang="ts">
import { CopyIcon, EyeIcon, EyeOffIcon, RefreshCwIcon } from 'lucide-vue-next'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'

import { useProfile, useUserSettings } from '@/api/auth'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Switch } from '@/components/ui/switch'
import { toast } from 'vue-sonner'

const { data: profile, executeQuery } = useProfile()

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
			<Card class="p-2">
				<CardContent>
					<div class="flex items-center gap-2">
						<Switch
							:model-value="!profile?.hideOnLandingPage"
							@update:model-value="changeLandingVisibility"
						/>
						<Label>{{ t('userSettings.account.showMeOnLanding') }}</Label>
					</div>
				</CardContent>
			</Card>
		</div>

		<div class="flex flex-col gap-6">
			<h4 class="scroll-m-20 text-xl font-semibold tracking-tight">Api</h4>

			<Card>
				<CardHeader>
					<CardTitle>Key</CardTitle>
				</CardHeader>
				<CardContent class="space-y-4">
					<div class="flex gap-2 w-full flex-wrap">
						<Input
							:type="showApiKey ? 'text' : 'password'"
							:model-value="profile?.apiKey ?? ''"
							class="flex-1"
							readonly
						/>
						<Button variant="outline" size="icon" type="button" @click="showApiKey = !showApiKey">
							<EyeIcon v-if="!showApiKey" />
							<EyeOffIcon v-else />
						</Button>
						<Button variant="outline" size="icon" type="button" @click="copyApiKey">
							<CopyIcon />
						</Button>
						<Button variant="outline" class="min-w-37.5 sm:w-full" @click="callRegenerateKey">
							<RefreshCwIcon />
							{{ t('userSettings.account.regenerateApiKey.button') }}
						</Button>
					</div>
					<p class="text-sm text-muted-foreground">
						{{ t('userSettings.account.regenerateApiKey.info') }}
					</p>
				</CardContent>
			</Card>
		</div>
	</div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

import { useProfile, useUserSettings } from '~~/layers/dashboard/api/auth'
import { Badge } from '@/components/ui/badge'
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
const regenerateChannelApiKey = userManager.useChannelApiKeyGenerateMutation()

const { t } = useI18n()

const showChannelApiKey = ref(false)
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

async function callRegenerateChannelKey() {
	const result = await regenerateChannelApiKey.executeMutation({})

	if (result.error) {
		toast.error('Failed to regenerate API key')
		return
	}

	await executeQuery({ requestPolicy: 'network-only' })

	toast.success(t('sharedTexts.saved'))
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

async function copyChannelApiKey() {
	if (!profile.value?.channelApiKey) return

	try {
		await navigator.clipboard.writeText(profile.value.channelApiKey)
		toast.success(t('sharedTexts.copied'), {
			duration: 1500,
		})
	} catch (err) {
		toast.error('Failed to copy')
	}
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
					<CardTitle>{{ t('userSettings.account.channelApiKey.title') }}</CardTitle>
				</CardHeader>
				<CardContent class="space-y-4">
					<div class="flex gap-2 w-full flex-wrap">
						<Input
							:type="showChannelApiKey ? 'text' : 'password'"
							:model-value="profile?.channelApiKey ?? ''"
							class="flex-1"
							readonly
						/>
						<Button variant="outline" size="icon" type="button" @click="showChannelApiKey = !showChannelApiKey">
							<Icon name="lucide:eye" v-if="!showChannelApiKey" />
							<Icon name="lucide:eye-off" v-else />
						</Button>
						<Button variant="outline" size="icon" type="button" @click="copyChannelApiKey">
							<Icon name="lucide:copy"  />
						</Button>
						<Button variant="outline" class="min-w-37.5 sm:w-full" @click="callRegenerateChannelKey">
							<Icon name="lucide:refresh-cw"  />
							{{ t('userSettings.account.channelApiKey.button') }}
						</Button>
					</div>
					<p class="text-sm text-muted-foreground">
						{{ t('userSettings.account.channelApiKey.info') }}
					</p>
				</CardContent>
			</Card>

			<Card>
				<CardHeader>
					<CardTitle class="flex items-center gap-2">
						{{ t('userSettings.account.legacyApiKey.title') }}
						<Badge variant="destructive">
							{{ t('userSettings.account.legacyApiKey.badge') }}
						</Badge>
					</CardTitle>
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
							<Icon name="lucide:eye" v-if="!showApiKey" />
							<Icon name="lucide:eye-off" v-else />
						</Button>
						<Button variant="outline" size="icon" type="button" @click="copyApiKey">
							<Icon name="lucide:copy"  />
						</Button>
						<Button variant="outline" class="min-w-37.5 sm:w-full" @click="callRegenerateKey">
							<Icon name="lucide:refresh-cw"  />
							{{ t('userSettings.account.regenerateApiKey.button') }}
						</Button>
					</div>
					<p class="text-sm text-muted-foreground">
						{{ t('userSettings.account.legacyApiKey.info') }}
					</p>
				</CardContent>
			</Card>
		</div>
	</div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'

import { useUserAccessFlagChecker } from '~~/layers/dashboard/api/auth'
import { useDotaApi } from '~~/layers/dashboard/api/dota'
import DialogOrSheet from '~~/layers/dashboard/components/dialog-or-sheet.vue'
import Card from '~~/layers/dashboard/components/card/card.vue'
import { useDota } from '~~/layers/dashboard/features/modules/composables/use-dota'
import { Button } from '@/components/ui/button'
import { CardContent, CardHeader, CardTitle, Card as UICard } from '@/components/ui/card'
import CopyInput from '@/components/ui/copy-input/CopyInput.vue'
import { Dialog, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import {
	FormControl,
	FormDescription,
	FormField,
	FormItem,
	FormLabel,
	FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Switch } from '@/components/ui/switch'
import { ChannelRolePermissionEnum } from '~/gql/graphql.js'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()

const showSettings = ref(false)

const {
	form,
	handleSubmit,
	isLoading,
	fetching,
	settings,
	linkSteam,
	unlinkSteam,
	resetSession,
	regenerateGsiToken,
} = useDota()

const dotaApi = useDotaApi()
const { data: steamAuthLinkData } = dotaApi.useQueryDotaSteamAuthLink()

const canManageModules = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageModules)

onMounted(async () => {
	if (!route.query.dotaSteamCallback) return

	const queryString = window.location.search
	await router.replace({ query: {} })
	await linkSteam(queryString)
})

function downloadGsiConfig() {
	const config = settings.value?.gsiConfig
	if (!config) return

	const blob = new Blob([config], { type: 'text/plain' })
	const url = URL.createObjectURL(blob)
	const link = document.createElement('a')
	link.href = url
	link.download = 'gamestate_integration_twir.cfg'
	link.click()
	URL.revokeObjectURL(url)
}
</script>

<template>
	<Card
		:title="t('modules.dota.title')"
		:is-loading="fetching"
		icon="simple-icons:dota2"
		icon-height="30px"
		icon-width="30px"
		:description="t('modules.dota.description')"
	>
		<template #footer>
			<Button
				class="flex gap-2 items-center"
				variant="secondary"
				:disabled="!canManageModules"
				@click="showSettings = !showSettings"
			>
				{{ t('sharedTexts.settings') }}
				<Icon name="lucide:settings" class="size-4" />
			</Button>
		</template>
	</Card>

	<Dialog v-model:open="showSettings">
		<DialogOrSheet>
			<DialogHeader>
				<DialogTitle>{{ t('modules.dota.settings.title') }}</DialogTitle>
			</DialogHeader>

			<form class="space-y-6 py-4" @submit.prevent="handleSubmit">
				<FormField v-slot="{ componentField }" name="enabled">
					<FormItem class="flex flex-row items-center justify-between rounded-lg border p-4">
						<div class="space-y-0.5">
							<FormLabel class="text-base">
								{{ t('modules.dota.settings.enabled.label') }}
							</FormLabel>
							<FormDescription>
								{{ t('modules.dota.settings.enabled.description') }}
							</FormDescription>
						</div>
						<FormControl>
							<Switch
								:model-value="componentField.modelValue"
								@update:model-value="componentField['onUpdate:modelValue']"
							/>
						</FormControl>
					</FormItem>
				</FormField>

				<UICard>
					<CardHeader>
						<CardTitle>{{ t('modules.dota.steam.title') }}</CardTitle>
					</CardHeader>
					<CardContent class="space-y-4">
						<div
							v-if="settings?.steamAccountId"
							class="flex items-center justify-between gap-4"
						>
							<div class="flex items-center gap-3 min-w-0">
								<img
									v-if="settings.steamProfile?.avatar"
									:src="settings.steamProfile.avatar"
									alt=""
									class="size-10 rounded-full"
								>
								<div class="min-w-0">
									<div class="font-medium truncate">
										{{ settings.steamProfile?.name || settings.steamAccountId }}
									</div>
									<a
										v-if="settings.steamProfile"
										:href="settings.steamProfile.profileUrl"
										target="_blank"
										class="text-sm text-muted-foreground hover:underline"
									>
										{{ settings.steamAccountId }}
									</a>
								</div>
							</div>
							<Button
								type="button"
								variant="destructive"
								:disabled="!canManageModules"
								@click="unlinkSteam"
							>
								{{ t('modules.dota.steam.unlink') }}
							</Button>
						</div>
						<div v-else class="space-y-2">
							<p class="text-sm text-muted-foreground">
								{{ t('modules.dota.steam.notLinked') }}
							</p>
							<Button
								v-if="steamAuthLinkData?.dotaSteamAuthLink"
								as="a"
								:href="steamAuthLinkData.dotaSteamAuthLink"
								type="button"
								class="flex gap-2 items-center"
							>
								<Icon name="simple-icons:steam" class="size-4" />
								{{ t('modules.dota.steam.link') }}
							</Button>
						</div>
					</CardContent>
				</UICard>

				<UICard>
					<CardHeader>
						<CardTitle>{{ t('modules.dota.gsi.title') }}</CardTitle>
					</CardHeader>
					<CardContent class="space-y-4">
						<div class="space-y-2">
							<FormLabel>{{ t('modules.dota.gsi.token') }}</FormLabel>
							<div class="flex gap-2">
								<CopyInput :text="settings?.gsiToken ?? ''" type="password" />
								<Button
									type="button"
									variant="secondary"
									:disabled="!canManageModules"
									@click="regenerateGsiToken"
								>
									{{ t('modules.dota.gsi.regenerate') }}
								</Button>
							</div>
						</div>
						<p class="text-sm text-muted-foreground">
							{{ t('modules.dota.gsi.instructions') }}
						</p>
						<Button
							v-if="settings?.gsiConfig"
							type="button"
							variant="secondary"
							class="flex gap-2 items-center"
							@click="downloadGsiConfig"
						>
							<Icon name="lucide:download" class="size-4" />
							{{ t('modules.dota.gsi.download') }}
						</Button>
						<p v-else class="text-sm text-destructive">
							{{ t('modules.dota.gsi.notConfigured') }}
						</p>
					</CardContent>
				</UICard>

				<UICard>
					<CardHeader>
						<CardTitle>{{ t('modules.dota.mmr.title') }}</CardTitle>
					</CardHeader>
					<CardContent class="space-y-4">
						<div class="grid grid-cols-2 gap-4">
							<FormField v-slot="{ componentField }" name="mmr">
								<FormItem>
									<FormLabel>{{ t('modules.dota.mmr.mmr.label') }}</FormLabel>
									<FormControl>
										<Input v-bind="componentField" type="number" />
									</FormControl>
									<FormMessage />
								</FormItem>
							</FormField>

							<FormField v-slot="{ componentField }" name="mmrDelta">
								<FormItem>
									<FormLabel>{{ t('modules.dota.mmr.mmrDelta.label') }}</FormLabel>
									<FormControl>
										<Input v-bind="componentField" type="number" :min="1" :max="100" />
									</FormControl>
									<FormDescription>
										{{ t('modules.dota.mmr.mmrDelta.description') }}
									</FormDescription>
									<FormMessage />
								</FormItem>
							</FormField>
						</div>

						<div class="flex items-center justify-between rounded-lg border p-4">
							<div>
								{{ t('modules.dota.mmr.session') }}:
								<span class="font-medium">
									{{ settings?.sessionWins ?? 0 }}W - {{ settings?.sessionLosses ?? 0 }}L
								</span>
							</div>
							<Button
								type="button"
								variant="secondary"
								:disabled="!canManageModules"
								@click="resetSession"
							>
								{{ t('modules.dota.mmr.resetSession') }}
							</Button>
						</div>
					</CardContent>
				</UICard>

				<UICard>
					<CardHeader>
						<CardTitle>{{ t('modules.dota.predictions.title') }}</CardTitle>
					</CardHeader>
					<CardContent class="space-y-4">
						<FormField v-slot="{ componentField }" name="predictionSettings.enabled">
							<FormItem class="flex flex-row items-center justify-between rounded-lg border p-4">
								<div class="space-y-0.5">
									<FormLabel class="text-base">
										{{ t('modules.dota.predictions.enabled.label') }}
									</FormLabel>
									<FormDescription>
										{{ t('modules.dota.predictions.enabled.description') }}
									</FormDescription>
								</div>
								<FormControl>
									<Switch
										:model-value="componentField.modelValue"
										@update:model-value="componentField['onUpdate:modelValue']"
									/>
								</FormControl>
							</FormItem>
						</FormField>

						<FormField v-slot="{ componentField }" name="predictionSettings.titleTemplate">
							<FormItem>
								<FormLabel>{{ t('modules.dota.predictions.titleTemplate.label') }}</FormLabel>
								<FormControl>
									<Input v-bind="componentField" />
								</FormControl>
								<FormMessage />
							</FormItem>
						</FormField>

						<FormField v-slot="{ componentField }" name="predictionSettings.windowSeconds">
							<FormItem>
								<FormLabel>{{ t('modules.dota.predictions.windowSeconds.label') }}</FormLabel>
								<FormControl>
									<Input v-bind="componentField" type="number" :min="30" :max="1800" />
								</FormControl>
								<FormDescription>
									{{ t('modules.dota.predictions.windowSeconds.description') }}
								</FormDescription>
								<FormMessage />
							</FormItem>
						</FormField>
					</CardContent>
				</UICard>

				<UICard>
					<CardHeader>
						<CardTitle>{{ t('modules.dota.chatEvents.title') }}</CardTitle>
					</CardHeader>
					<CardContent class="space-y-4">
						<div
							v-for="eventKey in ['matchStarted', 'matchEnded', 'roshanKilled', 'aegisPickup'] as const"
							:key="eventKey"
							class="space-y-2 rounded-lg border p-4"
						>
							<FormField v-slot="{ componentField }" :name="`chatEvents.${eventKey}.enabled`">
								<FormItem class="flex flex-row items-center justify-between">
									<FormLabel class="text-base">
										{{ t(`modules.dota.chatEvents.${eventKey}.label`) }}
									</FormLabel>
									<FormControl>
										<Switch
											:model-value="componentField.modelValue"
											@update:model-value="componentField['onUpdate:modelValue']"
										/>
									</FormControl>
								</FormItem>
							</FormField>

							<FormField v-slot="{ componentField }" :name="`chatEvents.${eventKey}.template`">
								<FormItem>
									<FormControl>
										<Input v-bind="componentField" />
									</FormControl>
									<FormDescription>
										{{ t(`modules.dota.chatEvents.${eventKey}.variables`) }}
									</FormDescription>
									<FormMessage />
								</FormItem>
							</FormField>

							<FormField v-slot="{ componentField }" :name="`chatEvents.${eventKey}.cooldown`">
								<FormItem>
									<FormLabel>{{ t('modules.dota.chatEvents.cooldown.label') }}</FormLabel>
									<FormControl>
										<Input v-bind="componentField" type="number" :min="0" />
									</FormControl>
									<FormMessage />
								</FormItem>
							</FormField>
						</div>
					</CardContent>
				</UICard>

				<UICard>
					<CardHeader>
						<CardTitle>{{ t('modules.dota.commands.title') }}</CardTitle>
					</CardHeader>
					<CardContent>
						<div class="grid grid-cols-2 gap-2">
							<FormField
								v-for="command in ['mmr', 'wl', 'lg', 'gm', 'np', 'wp'] as const"
								:key="command"
								v-slot="{ componentField }"
								:name="`commandsSettings.${command}`"
							>
								<FormItem class="flex flex-row items-center justify-between rounded-lg border p-3">
									<FormLabel class="font-mono">
										!{{ command }}
									</FormLabel>
									<FormControl>
										<Switch
											:model-value="componentField.modelValue"
											@update:model-value="componentField['onUpdate:modelValue']"
										/>
									</FormControl>
								</FormItem>
							</FormField>
						</div>
					</CardContent>
				</UICard>

				<div class="flex justify-end">
					<Button type="submit" :disabled="isLoading || !canManageModules">
						{{ t('modules.dota.settings.buttons.save') }}
					</Button>
				</div>
			</form>
		</DialogOrSheet>
	</Dialog>
</template>

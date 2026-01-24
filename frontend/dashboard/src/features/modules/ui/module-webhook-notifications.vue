<script setup lang="ts">
import { SettingsIcon, Webhook } from 'lucide-vue-next'
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import { useUserAccessFlagChecker } from '@/api/auth'
import Card from '@/components/card/card.vue'
import { Button } from '@/components/ui/button'
import { Card as UICard, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import {
	FormControl,
	FormDescription,
	FormField,
	FormItem,
	FormLabel,
	FormMessage,
} from '@/components/ui/form'
import { Switch } from '@/components/ui/switch'
import { useWebhookNotifications } from '@/features/modules/composables/use-webhook-notifications'
import { ChannelRolePermissionEnum } from '@/gql/graphql'

const { t } = useI18n()

const showSettings = ref(false)

const { handleSubmit, isLoading, fetching, settings, exists } = useWebhookNotifications()

const canManageModules = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageModules)

const isEnabled = computed(() => settings.value?.enabled ?? false)
</script>

<template>
	<Card
		:title="t('modules.webhookNotifications.title')"
		:is-loading="fetching"
		:icon="Webhook"
		icon-height="30px"
		icon-width="30px"
		:description="t('modules.webhookNotifications.description')"
	>
		<template #content>
			<p class="text-sm text-muted-foreground">
				{{
					isEnabled
						? t('modules.webhookNotifications.status.enabled')
						: t('modules.webhookNotifications.status.disabled')
				}}
			</p>
		</template>

		<template #footer>
			<Button
				class="flex gap-2 items-center"
				variant="secondary"
				:disabled="!canManageModules"
				@click="showSettings = !showSettings"
			>
				{{ t('sharedTexts.settings') }}
				<SettingsIcon class="size-4" />
			</Button>
		</template>
	</Card>

	<Dialog v-model:open="showSettings">
		<DialogContent class="sm:max-w-[680px]">
			<DialogHeader>
				<DialogTitle>{{ t('modules.webhookNotifications.settings.title') }}</DialogTitle>
			</DialogHeader>

			<form class="space-y-6" @submit.prevent="handleSubmit">
				<p class="text-sm text-muted-foreground">
					{{ t('modules.webhookNotifications.settings.info') }}
				</p>

				<FormField v-slot="{ value, handleChange }" name="enabled">
					<FormItem class="flex flex-row items-center justify-between rounded-lg border p-4">
						<div class="space-y-0.5">
							<FormLabel class="text-base">
								{{ t('modules.webhookNotifications.settings.enabled.label') }}
							</FormLabel>
							<FormDescription>
								{{ t('modules.webhookNotifications.settings.enabled.description') }}
							</FormDescription>
						</div>
						<FormControl>
							<Switch :model-value="value" @update:model-value="handleChange" />
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>

				<UICard>
					<CardHeader>
						<CardTitle>{{ t('modules.webhookNotifications.settings.github.title') }}</CardTitle>
					</CardHeader>
					<CardContent class="space-y-4">
						<FormField v-slot="{ value, handleChange }" name="githubIssues">
							<FormItem class="flex flex-row items-center justify-between rounded-lg border p-4">
								<div class="space-y-0.5">
									<FormLabel class="text-base">
										{{ t('modules.webhookNotifications.settings.github.issues.label') }}
									</FormLabel>
									<FormDescription>
										{{ t('modules.webhookNotifications.settings.github.issues.description') }}
									</FormDescription>
								</div>
								<FormControl>
									<Switch :model-value="value" @update:model-value="handleChange" />
								</FormControl>
								<FormMessage />
							</FormItem>
						</FormField>

						<FormField v-slot="{ value, handleChange }" name="githubPullRequests">
							<FormItem class="flex flex-row items-center justify-between rounded-lg border p-4">
								<div class="space-y-0.5">
									<FormLabel class="text-base">
										{{ t('modules.webhookNotifications.settings.github.pullRequests.label') }}
									</FormLabel>
									<FormDescription>
										{{ t('modules.webhookNotifications.settings.github.pullRequests.description') }}
									</FormDescription>
								</div>
								<FormControl>
									<Switch :model-value="value" @update:model-value="handleChange" />
								</FormControl>
								<FormMessage />
							</FormItem>
						</FormField>

						<FormField v-slot="{ value, handleChange }" name="githubCommits">
							<FormItem class="flex flex-row items-center justify-between rounded-lg border p-4">
								<div class="space-y-0.5">
									<FormLabel class="text-base">
										{{ t('modules.webhookNotifications.settings.github.commits.label') }}
									</FormLabel>
									<FormDescription>
										{{ t('modules.webhookNotifications.settings.github.commits.description') }}
									</FormDescription>
								</div>
								<FormControl>
									<Switch :model-value="value" @update:model-value="handleChange" />
								</FormControl>
								<FormMessage />
							</FormItem>
						</FormField>
					</CardContent>
				</UICard>

				<div class="flex justify-end">
					<Button type="submit" :disabled="isLoading">
						{{
							exists
								? t('modules.webhookNotifications.settings.buttons.update')
								: t('modules.webhookNotifications.settings.buttons.create')
						}}
					</Button>
				</div>
			</form>
		</DialogContent>
	</Dialog>
</template>

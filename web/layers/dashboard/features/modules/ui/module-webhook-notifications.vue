<script lang="ts" setup>
import { SettingsIcon, Webhook } from 'lucide-vue-next'
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import { useUserAccessFlagChecker } from '@/api/auth'
import Card from '@/components/card/card.vue'
import { Button } from '@/components/ui/button'
import { CardContent, CardHeader, CardTitle, Card as UICard } from '@/components/ui/card'
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
		:description="t('modules.webhookNotifications.description')"
		:icon="Webhook"
		:is-loading="fetching"
		:title="t('modules.webhookNotifications.title')"
		icon-height="30px"
		icon-width="30px"
	>
		<template #content>
			<p class="text-muted-foreground text-sm">
				{{
					isEnabled
						? t('modules.webhookNotifications.status.enabled')
						: t('modules.webhookNotifications.status.disabled')
				}}
			</p>
		</template>

		<template #footer>
			<Button
				:disabled="!canManageModules"
				class="flex items-center gap-2"
				variant="secondary"
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

			<form
				class="space-y-6"
				@submit.prevent="handleSubmit"
			>
				<p class="text-muted-foreground text-sm">
					{{ t('modules.webhookNotifications.settings.info') }}
				</p>

				<FormField
					v-slot="{ value, handleChange }"
					name="enabled"
				>
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
							<Switch
								:model-value="value"
								@update:model-value="handleChange"
							/>
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>

				<UICard>
					<CardHeader>
						<CardTitle>{{ t('modules.webhookNotifications.settings.github.title') }}</CardTitle>
					</CardHeader>
					<CardContent class="space-y-4">
						<!-- Issues -->
						<div class="space-y-3 rounded-lg border p-4">
							<FormField
								v-slot="{ value, handleChange }"
								name="githubIssues"
							>
								<FormItem class="flex flex-row items-center justify-between">
									<div class="space-y-0.5">
										<FormLabel class="text-base">
											{{ t('modules.webhookNotifications.settings.github.issues.label') }}
										</FormLabel>
										<FormDescription>
											{{ t('modules.webhookNotifications.settings.github.issues.description') }}
										</FormDescription>
									</div>
									<FormControl>
										<Switch
											:model-value="value"
											@update:model-value="handleChange"
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							</FormField>
							<div class="space-y-2 border-l pl-4">
								<FormField
									v-slot="{ value, handleChange }"
									name="githubIssuesOnline"
								>
									<FormItem class="flex flex-row items-center justify-between">
										<div class="space-y-0.5">
											<FormLabel>
												{{ t('modules.webhookNotifications.settings.github.issues.online.label') }}
											</FormLabel>
											<FormDescription>
												{{
													t(
														'modules.webhookNotifications.settings.github.issues.online.description'
													)
												}}
											</FormDescription>
										</div>
										<FormControl>
											<Switch
												:model-value="value"
												@update:model-value="handleChange"
											/>
										</FormControl>
										<FormMessage />
									</FormItem>
								</FormField>
								<FormField
									v-slot="{ value, handleChange }"
									name="githubIssuesOffline"
								>
									<FormItem class="flex flex-row items-center justify-between">
										<div class="space-y-0.5">
											<FormLabel>
												{{ t('modules.webhookNotifications.settings.github.issues.offline.label') }}
											</FormLabel>
											<FormDescription>
												{{
													t(
														'modules.webhookNotifications.settings.github.issues.offline.description'
													)
												}}
											</FormDescription>
										</div>
										<FormControl>
											<Switch
												:model-value="value"
												@update:model-value="handleChange"
											/>
										</FormControl>
										<FormMessage />
									</FormItem>
								</FormField>
							</div>
						</div>

						<!-- Pull Requests -->
						<div class="space-y-3 rounded-lg border p-4">
							<FormField
								v-slot="{ value, handleChange }"
								name="githubPullRequests"
							>
								<FormItem class="flex flex-row items-center justify-between">
									<div class="space-y-0.5">
										<FormLabel class="text-base">
											{{ t('modules.webhookNotifications.settings.github.pullRequests.label') }}
										</FormLabel>
										<FormDescription>
											{{
												t('modules.webhookNotifications.settings.github.pullRequests.description')
											}}
										</FormDescription>
									</div>
									<FormControl>
										<Switch
											:model-value="value"
											@update:model-value="handleChange"
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							</FormField>
							<div class="space-y-2 border-l pl-4">
								<FormField
									v-slot="{ value, handleChange }"
									name="githubPullRequestsOnline"
								>
									<FormItem class="flex flex-row items-center justify-between">
										<div class="space-y-0.5">
											<FormLabel>
												{{
													t(
														'modules.webhookNotifications.settings.github.pullRequests.online.label'
													)
												}}
											</FormLabel>
											<FormDescription>
												{{
													t(
														'modules.webhookNotifications.settings.github.pullRequests.online.description'
													)
												}}
											</FormDescription>
										</div>
										<FormControl>
											<Switch
												:model-value="value"
												@update:model-value="handleChange"
											/>
										</FormControl>
										<FormMessage />
									</FormItem>
								</FormField>
								<FormField
									v-slot="{ value, handleChange }"
									name="githubPullRequestsOffline"
								>
									<FormItem class="flex flex-row items-center justify-between">
										<div class="space-y-0.5">
											<FormLabel>
												{{
													t(
														'modules.webhookNotifications.settings.github.pullRequests.offline.label'
													)
												}}
											</FormLabel>
											<FormDescription>
												{{
													t(
														'modules.webhookNotifications.settings.github.pullRequests.offline.description'
													)
												}}
											</FormDescription>
										</div>
										<FormControl>
											<Switch
												:model-value="value"
												@update:model-value="handleChange"
											/>
										</FormControl>
										<FormMessage />
									</FormItem>
								</FormField>
							</div>
						</div>

						<!-- Commits -->
						<div class="space-y-3 rounded-lg border p-4">
							<FormField
								v-slot="{ value, handleChange }"
								name="githubCommits"
							>
								<FormItem class="flex flex-row items-center justify-between">
									<div class="space-y-0.5">
										<FormLabel class="text-base">
											{{ t('modules.webhookNotifications.settings.github.commits.label') }}
										</FormLabel>
										<FormDescription>
											{{ t('modules.webhookNotifications.settings.github.commits.description') }}
										</FormDescription>
									</div>
									<FormControl>
										<Switch
											:model-value="value"
											@update:model-value="handleChange"
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							</FormField>
							<div class="space-y-2 border-l pl-4">
								<FormField
									v-slot="{ value, handleChange }"
									name="githubCommitsOnline"
								>
									<FormItem class="flex flex-row items-center justify-between">
										<div class="space-y-0.5">
											<FormLabel>
												{{ t('modules.webhookNotifications.settings.github.commits.online.label') }}
											</FormLabel>
											<FormDescription>
												{{
													t(
														'modules.webhookNotifications.settings.github.commits.online.description'
													)
												}}
											</FormDescription>
										</div>
										<FormControl>
											<Switch
												:model-value="value"
												@update:model-value="handleChange"
											/>
										</FormControl>
										<FormMessage />
									</FormItem>
								</FormField>
								<FormField
									v-slot="{ value, handleChange }"
									name="githubCommitsOffline"
								>
									<FormItem class="flex flex-row items-center justify-between">
										<div class="space-y-0.5">
											<FormLabel>
												{{
													t('modules.webhookNotifications.settings.github.commits.offline.label')
												}}
											</FormLabel>
											<FormDescription>
												{{
													t(
														'modules.webhookNotifications.settings.github.commits.offline.description'
													)
												}}
											</FormDescription>
										</div>
										<FormControl>
											<Switch
												:model-value="value"
												@update:model-value="handleChange"
											/>
										</FormControl>
										<FormMessage />
									</FormItem>
								</FormField>
							</div>
						</div>
					</CardContent>
				</UICard>

				<div class="flex justify-end">
					<Button
						:disabled="isLoading"
						type="submit"
					>
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

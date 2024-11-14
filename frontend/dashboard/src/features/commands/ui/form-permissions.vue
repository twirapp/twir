<script setup lang="ts">
import { ExternalLink } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'

import TwitchUsersSelect from '@/components/twitchUsers/twitch-users-select.vue'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Checkbox } from '@/components/ui/checkbox'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Label } from '@/components/ui/label'
import { useCommandEditV2 } from '@/features/commands/composables/use-command-edit-v2'

const { t } = useI18n()

const { channelRoles } = useCommandEditV2()

function confirmLeaving(e: Event, onLeave: () => any) {
	e.preventDefault()

	// eslint-disable-next-line no-alert
	if (confirm('You are leaving the page. Are you sure?')) {
		return onLeave()
	}
}
</script>

<template>
	<div class="flex flex-col gap-4">
		<Card>
			<CardHeader>
				<CardTitle>{{ t('commands.modal.permissions.divider') }}</CardTitle>
				<CardDescription>If none of roles selected - then command available for all users</CardDescription>
			</CardHeader>
			<CardContent class="flex flex-col gap-4">
				<div class="flex flex-col gap-2">
					<Label>
						<span>Role restriction</span>
						<RouterLink
							v-slot="{ href, navigate }"
							custom
							to="/dashboard/community?tab=permissions"
						>
							<a
								:href="href"
								class="ml-1 inline-flex text-xs underline"
								@click="(e) => confirmLeaving(e, navigate)"
							>
								manage roles <ExternalLink class="size-4" />
							</a>
						</RouterLink>
					</Label>
					<div class="grid grid-cols-1 md:grid-cols-2 gap-1 xl:max-w-[50%]">
						<FormField
							v-for="(role, index) in channelRoles?.roles" v-slot="{ value, handleChange }"
							:key="role.id"
							type="checkbox"
							:value="role.id"
							:unchecked-value="false"
							name="rolesIds"
						>
							<div v-if="index === 0" class="role">
								<Checkbox id="allRoles" :checked="!value?.length" disabled />
								<Label for="allRoles" class="capitalize">Everyone</Label>
							</div>

							<FormItem class="space-y-0">
								<FormLabel class="flex flex-row items-center gap-2 space-y-0 bg-accent px-3 py-2 rounded-md leading-5">
									<FormControl>
										<Checkbox
											:checked="value?.includes(role.id)"
											@update:checked="handleChange"
										/>
									</FormControl>
									<span>{{ role.name }}</span>
								</FormLabel>
							</FormItem>
						</FormField>
					</div>
				</div>

				<FormField v-slot="{ componentField }" name="allowedUsersIds">
					<FormItem>
						<FormLabel>Exceptions</FormLabel>
						<FormControl>
							<TwitchUsersSelect
								multiple
								:initial="componentField.modelValue"
								:model-value="componentField.modelValue"
								@update:model-value="componentField['onUpdate:modelValue']"
							/>
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>

				<FormField v-slot="{ componentField }" name="deniedUsersIds">
					<FormItem>
						<FormLabel>Blocked users</FormLabel>
						<FormControl>
							<TwitchUsersSelect
								multiple
								:initial="componentField.modelValue"
								:model-value="componentField.modelValue"
								@update:model-value="componentField['onUpdate:modelValue']"
							/>
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>
			</CardContent>
		</Card>

		<Card>
			<CardHeader>
				<CardTitle>{{ t('commands.modal.restrictions.name') }}</CardTitle>
			</CardHeader>
		</Card>
	</div>
</template>

<style scoped>
.role {
	@apply flex flex-row items-center gap-2 space-y-0 bg-accent px-3 py-2 rounded-md leading-5
}
</style>

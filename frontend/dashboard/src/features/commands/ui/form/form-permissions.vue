<script setup lang="ts">
import { ExternalLink } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'

import TwitchUsersSelect from '@/components/twitchUsers/twitch-users-select.vue'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import FormRolesSelector from '@/features/commands/ui/form-roles-selector.vue'

const { t } = useI18n()

function confirmLeaving(e: Event, onLeave: () => any) {
	e.preventDefault()

	// eslint-disable-next-line no-alert
	if (confirm('You are leaving the page. Are you sure?')) {
		return onLeave()
	}
}
</script>

<template>
	<Card>
		<CardHeader>
			<CardTitle>{{ t('commands.modal.permissions.divider') }}</CardTitle>
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

				<FormRolesSelector field-name="rolesIds" />
			</div>

			<FormField v-slot="{ componentField }" name="allowedUsersIds">
				<FormItem>
					<FormLabel>{{ t('commands.modal.permissions.exceptions') }}</FormLabel>
					<FormControl>
						<TwitchUsersSelect
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
					<FormLabel>{{ t('commands.modal.permissions.blocked') }}</FormLabel>
					<FormControl>
						<TwitchUsersSelect
							:initial="componentField.modelValue"
							:model-value="componentField.modelValue"
							@update:model-value="componentField['onUpdate:modelValue']"
						/>
					</FormControl>
					<FormMessage />
				</FormItem>
			</FormField>

			<div class="flex flex-col gap-2">
				<FormField v-slot="{ componentField }" name="requiredWatchTime">
					<FormItem class="required-block">
						<FormLabel>{{ t('commands.modal.restrictions.watchTime') }}</FormLabel>
						<FormControl>
							<Input class="w-fit" type="number" v-bind="componentField" />
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>

				<FormField v-slot="{ componentField }" name="requiredMessages">
					<FormItem class="required-block">
						<FormLabel>{{ t('commands.modal.restrictions.messages') }}</FormLabel>
						<FormControl>
							<Input class="w-fit" type="number" v-bind="componentField" />
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>

				<FormField v-slot="{ componentField }" name="requiredUsedChannelPoints">
					<FormItem class="required-block">
						<FormLabel>{{ t('commands.modal.restrictions.channelsPoints') }}</FormLabel>
						<FormControl>
							<Input class="w-fit" type="number" v-bind="componentField" />
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>
			</div>
		</CardContent>
	</Card>
</template>

<style scoped>
.required-block {
	@apply flex flex-row flex-wrap justify-between items-center;
}
</style>

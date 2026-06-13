<script setup lang="ts">
import { EditIcon, ShieldUser } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'

import TwitchUsersSelect from '@/components/twitchUsers/twitch-users-select.vue'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import {
	FormControl,
	FormDescription,
	FormField,
	FormItem,
	FormLabel,
	FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import FormRolesSelector from '@/features/commands/ui/form-roles-selector.vue'
import CommunityRolesModal from '@/features/community-roles/community-roles-modal.vue'

const { t } = useI18n()
</script>

<template>
	<Card>
		<CardHeader class="flex flex-row place-content-center">
			<CardTitle class="flex items-center gap-2">
				<ShieldUser />
				{{ t('commands.modal.permissions.divider') }}
			</CardTitle>
		</CardHeader>
		<CardContent class="flex flex-col gap-4 pt-4">
			<div class="flex flex-col gap-2">
				<Label class="flex gap-1">
					<span>Role restriction</span>
					<CommunityRolesModal>
						<template #trigger>
							<span class="flex flex-row gap-1 items-center cursor-pointer underline">
								{{ t('sidebar.roles') }}
								<EditIcon class="size-4" />
							</span>
						</template>
					</CommunityRolesModal>
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
					<FormDescription>
						Users, who can bypass roles restriction for that command.
					</FormDescription>
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
					<FormDescription> Users, who cannot use that command. </FormDescription>
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
@reference '@/assets/index.css';

.required-block {
	@apply flex flex-row flex-wrap justify-between items-center;
}
</style>

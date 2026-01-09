<script setup lang="ts">
import { EditIcon, ShieldUser } from 'lucide-vue-next'


import TwitchUsersSelect from '#layers/dashboard/components/twitchUsers/twitch-users-select.vue'




import FormRolesSelector from '~/features/commands/ui/form-roles-selector.vue'
import CommunityRolesModal from '~/features/community-roles/community-roles-modal.vue'

const { t } = useI18n()
</script>

<template>
	<UiCard>
		<UiCardHeader class="flex flex-row place-content-center">
			<UiCardTitle class="flex items-center gap-2">
				<ShieldUser />
				{{ t('commands.modal.permissions.divider') }}
			</UiCardTitle>
		</UiCardHeader>
		<UiCardContent class="flex flex-col gap-4 pt-4">
			<div class="flex flex-col gap-2">
				<UiLabel class="flex gap-1">
					<span>Role restriction</span>
					<CommunityRolesModal>
						<template #trigger>
							<span class="flex flex-row gap-1 items-center cursor-pointer underline">
								{{ t('sidebar.roles') }}
								<EditIcon class="size-4" />
							</span>
						</template>
					</CommunityRolesModal>
				</UiLabel>

				<FormRolesSelector field-name="rolesIds" />
			</div>

			<UiFormField v-slot="{ componentField }" name="allowedUsersIds">
				<UiFormItem>
					<UiFormLabel>{{ t('commands.modal.permissions.exceptions') }}</UiFormLabel>
					<UiFormControl>
						<TwitchUsersSelect
							:initial="componentField.modelValue"
							:model-value="componentField.modelValue"
							@update:model-value="componentField['onUpdate:modelValue']"
						/>
					</UiFormControl>
					<UiFormMessage />
					<UiFormDescription>
						Users, who can bypass roles restriction for that command.
					</UiFormDescription>
				</UiFormItem>
			</UiFormField>

			<UiFormField v-slot="{ componentField }" name="deniedUsersIds">
				<UiFormItem>
					<UiFormLabel>{{ t('commands.modal.permissions.blocked') }}</UiFormLabel>
					<UiFormControl>
						<TwitchUsersSelect
							:initial="componentField.modelValue"
							:model-value="componentField.modelValue"
							@update:model-value="componentField['onUpdate:modelValue']"
						/>
					</UiFormControl>
					<UiFormMessage />
					<UiFormDescription> Users, who cannot use that command. </UiFormDescription>
				</UiFormItem>
			</UiFormField>

			<div class="flex flex-col gap-2">
				<UiFormField v-slot="{ componentField }" name="requiredWatchTime">
					<UiFormItem class="required-block">
						<UiFormLabel>{{ t('commands.modal.restrictions.watchTime') }}</UiFormLabel>
						<UiFormControl>
							<UiInput class="w-fit" type="number" v-bind="componentField" />
						</UiFormControl>
						<UiFormMessage />
					</UiFormItem>
				</UiFormField>

				<UiFormField v-slot="{ componentField }" name="requiredMessages">
					<UiFormItem class="required-block">
						<UiFormLabel>{{ t('commands.modal.restrictions.messages') }}</UiFormLabel>
						<UiFormControl>
							<UiInput class="w-fit" type="number" v-bind="componentField" />
						</UiFormControl>
						<UiFormMessage />
					</UiFormItem>
				</UiFormField>

				<UiFormField v-slot="{ componentField }" name="requiredUsedChannelPoints">
					<UiFormItem class="required-block">
						<UiFormLabel>{{ t('commands.modal.restrictions.channelsPoints') }}</UiFormLabel>
						<UiFormControl>
							<UiInput class="w-fit" type="number" v-bind="componentField" />
						</UiFormControl>
						<UiFormMessage />
					</UiFormItem>
				</UiFormField>
			</div>
		</UiCardContent>
	</UiCard>
</template>

<style scoped>
@reference '~/assets/css/tailwind.css';

.required-block {
	@apply flex flex-row flex-wrap justify-between items-center;
}
</style>

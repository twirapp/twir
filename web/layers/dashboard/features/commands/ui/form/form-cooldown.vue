<script setup lang="ts">
import { EditIcon, RefreshCcwIcon } from 'lucide-vue-next'






import FormRolesSelector from '~/features/commands/ui/form-roles-selector.vue'
import CommunityRolesModal from '~/features/community-roles/community-roles-modal.vue'

const { t } = useI18n()
</script>

<template>
	<UiCard>
		<UiCardHeader class="flex flex-row place-content-center flex-wrap">
			<UiCardTitle class="flex items-center gap-2">
				<RefreshCcwIcon />
				{{ t('commands.modal.cooldown.label') }}
			</UiCardTitle>
		</UiCardHeader>
		<UiCardContent class="pt-4">
			<div class="flex flex-col gap-4">
				<UiFormField v-slot="{ componentField }" name="cooldown">
					<UiFormItem>
						<UiFormLabel class="flex gap-2">
							{{ t('commands.modal.cooldown.value') }}
						</UiFormLabel>
						<UiFormControl>
							<UiInput type="number" v-bind="componentField" />
						</UiFormControl>
						<UiFormMessage />
					</UiFormItem>
				</UiFormField>

				<UiFormField v-slot="{ componentField }" name="cooldownType">
					<UiFormItem>
						<UiFormLabel>{{ t('commands.modal.cooldown.type.name') }}</UiFormLabel>
						<UiFormControl>
							<UiSelect v-bind="componentField">
								<UiFormControl>
									<UiSelectTrigger>
										<UiSelectValue />
									</UiSelectTrigger>
								</UiFormControl>
								<UiSelectContent>
									<UiSelectGroup>
										<UiSelectItem value="GLOBAL">
											{{ t('commands.modal.cooldown.type.global') }}
										</UiSelectItem>
										<UiSelectItem value="PER_USER">
											{{ t('commands.modal.cooldown.type.user') }}
										</UiSelectItem>
									</UiSelectGroup>
								</UiSelectContent>
							</UiSelect>
						</UiFormControl>
						<UiFormMessage />
					</UiFormItem>
				</UiFormField>

				<div class="flex flex-col gap-2">
					<span class="inline-flex gap-1">
						Affected roles
						<CommunityRolesModal>
							<template #trigger>
								<span class="flex flex-row gap-1 items-center cursor-pointer underline">
									{{ t('sidebar.roles') }}
									<EditIcon class="size-4" />
								</span>
							</template>
						</CommunityRolesModal>
					</span>
					<FormRolesSelector field-name="cooldownRolesIds" hide-broadcaster />
				</div>
			</div>
		</UiCardContent>
	</UiCard>
</template>

<style scoped></style>

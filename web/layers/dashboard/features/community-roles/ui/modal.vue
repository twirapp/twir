<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { onMounted, toRaw } from 'vue'

import { z } from 'zod'

import type { ChannelRolesQuery, RolesCreateOrUpdateOpts } from '~/gql/graphql'

import { PERMISSIONS_FLAGS } from '#layers/dashboard/api/auth'
import { useRoles } from '#layers/dashboard/api/roles'
import UsersMultiSearch from '#layers/dashboard/components/twitchUsers/twitch-users-select.vue'





import { toast } from 'vue-sonner'
import { ChannelRolePermissionEnum } from '~/gql/graphql'

const props = defineProps<{
	role?: ChannelRolesQuery['roles'][number] | null
}>()

const emit = defineEmits<{
	close: []
}>()

const { t } = useI18n()

const formSchema = toTypedSchema(
	z.object({
		name: z.string().min(1, t('roles.validations.nameRequired')).max(50),
		permissions: z.array(z.nativeEnum(ChannelRolePermissionEnum)),
		users: z.array(z.string()),
		settings: z.object({
			requiredMessages: z.number().min(0).max(99999999),
			requiredUserChannelPoints: z.number().min(0).max(999999999999),
			requiredWatchTime: z.number().min(0).max(99999999),
		}),
	})
)

const initialValues = {
	name: '',
	permissions: [],
	users: [],
	settings: {
		requiredMessages: 0,
		requiredUserChannelPoints: 0,
		requiredWatchTime: 0,
	},
}

const { handleSubmit, setValues } = useForm({
	validationSchema: formSchema,
	initialValues,
	keepValuesOnUnmount: true,
	validateOnMount: false,
})

onMounted(() => {
	if (!props.role) return
	const raw = structuredClone(toRaw(props.role))
	setValues({
		name: raw.name,
		permissions: raw.permissions,
		settings: raw.settings,
		users: props.role.users.map((u) => u.id),
	})
})

const rolesManager = useRoles()
const rolesUpdater = rolesManager.useRolesUpdateMutation()
const rolesCreator = rolesManager.useRolesCreateMutation()

const onSubmit = handleSubmit(async (formData) => {
	if (props.role?.id) {
		await rolesUpdater.executeMutation({
			id: props.role.id,
			opts: formData as RolesCreateOrUpdateOpts,
		})
	} else {
		await rolesCreator.executeMutation({
			opts: formData as RolesCreateOrUpdateOpts,
		})
	}

	toast.success(t('sharedTexts.saved'), {
		duration: 1500,
	})

	emit('close')
})
</script>

<template>
	<form>
		<div class="grid gap-6">
			<UiFormField v-slot="{ componentField }" name="name">
				<UiFormItem>
					<UiFormLabel>{{ t('sharedTexts.name') }}</UiFormLabel>
					<UiFormControl>
						<UiInput v-bind="componentField" />
					</UiFormControl>
					<UiFormMessage />
				</UiFormItem>
			</UiFormField>

			<UiSeparator />

			<div class="space-y-2">
				<h4 class="font-medium">
					{{ t('roles.modal.accessToUsers') }}
				</h4>
				<UiFormField v-slot="{ componentField }" name="users">
					<UiFormItem>
						<UsersMultiSearch
							:model-value="componentField.modelValue"
							@update:model-value="componentField['onUpdate:modelValue']"
						/>
						<UiFormMessage />
					</UiFormItem>
				</UiFormField>
			</div>

			<UiSeparator />

			<div class="space-y-2">
				<h4 class="font-medium">
					{{ t('roles.modal.accessByStats') }}
				</h4>
				<div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
					<UiFormField v-slot="{ componentField }" name="settings.requiredWatchTime">
						<UiFormItem>
							<UiFormLabel>{{ t('roles.modal.requiredWatchTime') }}</UiFormLabel>
							<UiFormControl>
								<UiInput type="number" v-bind="componentField" min="0" max="99999999" />
							</UiFormControl>
							<UiFormMessage />
						</UiFormItem>
					</UiFormField>

					<UiFormField v-slot="{ componentField }" name="settings.requiredMessages">
						<UiFormItem>
							<UiFormLabel>{{ t('roles.modal.requiredMessages') }}</UiFormLabel>
							<UiFormControl>
								<UiInput type="number" v-bind="componentField" min="0" max="99999999" />
							</UiFormControl>
							<UiFormMessage />
						</UiFormItem>
					</UiFormField>

					<UiFormField v-slot="{ componentField }" name="settings.requiredUserChannelPoints">
						<UiFormItem>
							<UiFormLabel>{{ t('roles.modal.requiredChannelPoints') }}</UiFormLabel>
							<UiFormControl>
								<UiInput type="number" v-bind="componentField" min="0" max="999999999999" />
							</UiFormControl>
							<UiFormMessage />
						</UiFormItem>
					</UiFormField>
				</div>
			</div>

			<UiSeparator />

			<div class="space-y-2">
				<h4 class="font-medium">
					{{ t('roles.modal.permissions') }}
				</h4>
				<div class="grid grid-cols-1 sm:grid-cols-2 gap-2">
					<UiFormField v-slot="{ value, handleChange }" name="permissions">
						<template v-for="(permission, index) of PERMISSIONS_FLAGS" :key="index">
							<div v-if="permission === 'delimiter'" class="col-span-2" />
							<UiFormItem v-else class="flex flex-row items-start space-x-3 space-y-0">
								<UiFormControl>
									<UiCheckbox
										:model-value="value?.includes(permission.perm)"
										:disabled="
											value?.some(
												(p: ChannelRolePermissionEnum) =>
													p === ChannelRolePermissionEnum.CanAccessDashboard
											) && permission.perm !== ChannelRolePermissionEnum.CanAccessDashboard
										"
										@update:model-value="
											(checked: boolean | 'indeterminate') => {
												if (checked) {
													handleChange([...(value || []), permission.perm])
												} else {
													handleChange(
														value?.filter(
															(p: ChannelRolePermissionEnum) => p !== permission.perm
														) || []
													)
												}
											}
										"
									/>
								</UiFormControl>
								<UiFormLabel class="font-normal">
									{{ permission.description }}
								</UiFormLabel>
							</UiFormItem>
						</template>
					</UiFormField>
				</div>
			</div>

			<UiButton type="submit" class="w-full" @click="onSubmit">
				{{ t('sharedButtons.save') }}
			</UiButton>
		</div>
	</form>
</template>

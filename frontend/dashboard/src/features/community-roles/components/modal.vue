<script setup lang='ts'>
import {
	type FormInst,
	type FormItemRule,
	type FormRules,
	NButton,
	NCheckbox,
	NCheckboxGroup,
	NDivider,
	NForm,
	NFormItem,
	NGrid,
	NGridItem,
	NInput,
	NInputNumber
} from 'naive-ui'
import { onMounted, ref, toRaw } from 'vue'
import { useI18n } from 'vue-i18n'

import { PERMISSIONS_FLAGS } from '@/api/index.js'
import { useRoles } from '@/api/roles'
import UsersMultiSearch from '@/components/twitchUsers/multiple.vue'
import { useToast } from '@/components/ui/toast'
import {
	ChannelRolePermissionEnum,
	type ChannelRolesQuery,
	type RolesCreateOrUpdateOpts
} from '@/gql/graphql'

const props = defineProps<{
	role?: ChannelRolesQuery['roles'][number] | null
}>()
defineEmits<{
	close: []
}>()

const { t } = useI18n()

const formRef = ref<FormInst | null>(null)
const formValue = ref<RolesCreateOrUpdateOpts>({
	name: '',
	permissions: [],
	users: [],
	settings: {
		requiredMessages: 0,
		requiredUserChannelPoints: 0,
		requiredWatchTime: 0
	}
})

const rules: FormRules = {
	name: {
		trigger: ['input', 'blur'],
		validator: (_: FormItemRule, value: string) => {
			if (!value || !value.length) {
				return new Error(t('roles.validations.nameRequired'))
			}

			return true
		}
	}
}

onMounted(() => {
	if (!props.role) return

	const raw = structuredClone(toRaw(props.role))
	formValue.value.name = raw.name
	formValue.value.permissions = raw.permissions
	formValue.value.settings = raw.settings
	if (props.role.users.length) {
		formValue.value.users = props.role.users.map(u => u.id)
	}
})

const rolesManager = useRoles()
const rolesUpdater = rolesManager.useRolesUpdateMutation()
const rolesCreator = rolesManager.useRolesCreateMutation()

const toast = useToast()

async function save() {
	if (!formRef.value || !formValue.value) return
	await formRef.value.validate()

	const data = formValue.value

	if (props.role?.id) {
		await rolesUpdater.executeMutation({
			id: props.role.id,
			opts: data
		})
	} else {
		await rolesCreator.executeMutation({
			opts: {
				...data
			}
		})
	}

	toast.toast({
		title: t('sharedTexts.saved'),
		duration: 1500
	})
}
</script>

<template>
	<NForm ref="formRef" :model="formValue" :rules="rules">
		<NFormItem :label="t('sharedTexts.name')" path="name" show-require-mark>
			<NInput v-model:value="formValue.name" />
		</NFormItem>

		<NDivider>{{ t('roles.modal.accessToUsers') }}</NDivider>

		<UsersMultiSearch v-model="formValue.users" />

		<NDivider>{{ t('roles.modal.accessByStats') }}</NDivider>

		<NGrid cols="1 s:2 m:2 l:2" responsive="screen" :x-gap="5">
			<NGridItem :span="1">
				<NFormItem :label="t('roles.modal.requiredWatchTime')">
					<NInputNumber
						v-model:value="formValue.settings!.requiredWatchTime"
						:min="0" :max="99999999"
					/>
				</NFormItem>
			</NGridItem>

			<NGridItem :span="1">
				<NFormItem :label="t('roles.modal.requiredMessages')">
					<NInputNumber
						v-model:value="formValue.settings!.requiredMessages"
						:min="0"
						:max="99999999"
					/>
				</NFormItem>
			</NGridItem>

			<NGridItem :span="1">
				<NFormItem :label="t('roles.modal.requiredChannelPoints')">
					<NInputNumber
						v-model:value="formValue.settings!.requiredUserChannelPoints"
						:min="0"
						:max="999999999999"
					/>
				</NFormItem>
			</NGridItem>
		</NGrid>

		<NDivider>{{ t('roles.modal.permissions') }}</NDivider>

		<NCheckboxGroup v-model:value="formValue.permissions">
			<NGrid cols="1 s:2 m:2 l:2" responsive="screen" :x-gap="5">
				<NGridItem
					v-for="(permission, index) of PERMISSIONS_FLAGS"
					:key="index"
					:span="permission === 'delimiter' ? 2 : 1"
				>
					<NCheckbox
						v-if="permission !== 'delimiter'"
						:disabled="formValue.permissions.some(p => p === ChannelRolePermissionEnum.CanAccessDashboard)
							&& permission.perm !== ChannelRolePermissionEnum.CanAccessDashboard
						"
						:value="permission.perm"
						:label="permission.description"
					/>
				</NGridItem>
			</NGrid>
		</NCheckboxGroup>

		<NDivider />

		<NButton secondary type="success" block class="mt-3.5" @click="save">
			{{ t('sharedButtons.save') }}
		</NButton>
	</NForm>
</template>

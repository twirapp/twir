<script setup lang='ts'>
import {
	type FormInst,
	type FormRules,
	type FormItemRule,
  NInput,
  NForm,
  NFormItem,
  NDivider,
  NGrid,
  NGridItem,
  NInputNumber,
  NCheckbox,
  NCheckboxGroup,
  NButton,
} from 'naive-ui';
import { ref, onMounted, toRaw, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import { useRolesManager, PERMISSIONS_FLAGS } from '@/api/index.js';
import { type EditableRole } from '@/components/roles/types.js';
import UsersMultiSearch from '@/components/twitchUsers/multiple.vue';

const props = defineProps<{
	role?: EditableRole | null
}>();
const emits = defineEmits<{
	close: []
}>();

const formRef = ref<FormInst | null>(null);
const formValue = ref<EditableRole>({
	name: '',
	type: '',
	permissions: [],
	users: [],
	settings: {
		requiredMessages: 0,
		requiredUserChannelPoints: 0,
		requiredWatchTime: 0,
	},
});

const rules: FormRules = {
	name: {
		trigger: ['input', 'blur'],
		validator: (_: FormItemRule, value: string) => {
			if (!value || !value.length) {
				return new Error(t('roles.validations.nameRequired'));
			}

			return true;
		},
	},
};

const searchUsersIds = ref<string[]>([]);
watch(searchUsersIds, (value) => {
	formValue.value.users = value.map((id) => ({ userId: id }));
});

onMounted(() => {
	if (!props.role) return;
	formValue.value = structuredClone(toRaw(props.role));
	searchUsersIds.value = formValue.value.users.map((u) => u.userId) ?? [];
});


const rolesManager = useRolesManager();
const rolesUpdater = rolesManager.update;
const rolesCreator = rolesManager.create;

async function save() {
	if (!formRef.value || !formValue.value) return;
	await formRef.value.validate();

	const data = formValue.value;

	if (data.id) {
		await rolesUpdater.mutateAsync({
			id: data.id,
			role: data,
		});
	} else {
		await rolesCreator.mutateAsync(data);
	}

	emits('close');
}

const { t } = useI18n();
</script>

<template>
	<n-form ref="formRef" :model="formValue" :rules="rules">
		<n-form-item :label="t('sharedTexts.name')" path="name" show-require-mark>
			<n-input v-model:value="formValue.name" />
		</n-form-item>

		<n-divider>{{ t('roles.modal.accessToUsers') }}</n-divider>

		<users-multi-search v-model="searchUsersIds" />

		<n-divider>{{ t('roles.modal.accessByStats') }}</n-divider>

		<n-grid cols="1 s:2 m:2 l:2" responsive="screen" :x-gap="5">
			<n-grid-item :span="1">
				<n-form-item :label="t('roles.modal.requiredWatchTime')">
					<n-input-number
						v-model:value="formValue.settings!.requiredWatchTime"
						:min="0" :max="99999999"
					/>
				</n-form-item>
			</n-grid-item>

			<n-grid-item :span="1">
				<n-form-item :label="t('roles.modal.requiredMessages')">
					<n-input-number
						v-model:value="formValue.settings!.requiredMessages"
						:min="0"
						:max="99999999"
					/>
				</n-form-item>
			</n-grid-item>

			<n-grid-item :span="1">
				<n-form-item :label="t('roles.modal.requiredChannelPoints')">
					<n-input-number
						v-model:value="formValue.settings!.requiredUserChannelPoints"
						:min="0"
						:max="999999999999"
					/>
				</n-form-item>
			</n-grid-item>
		</n-grid>

		<n-divider>{{ t('roles.modal.permissions') }}</n-divider>

		<n-checkbox-group v-model:value="formValue.permissions">
			<n-grid cols="1 s:2 m:2 l:2" responsive="screen" :x-gap="5">
				<n-grid-item
					v-for="(permission) of Object.entries(PERMISSIONS_FLAGS)"
					:key="permission[0]"
					:span="1"
				>
					<n-checkbox
						:disabled="formValue.permissions.some(p => p === 'CAN_ACCESS_DASHBOARD') &&
							permission[0] !== 'CAN_ACCESS_DASHBOARD'
						"
						:value="permission[0]"
						:label="permission[1]"
						:style="{ display: permission[1] == '' ? 'none' : undefined }"
					/>
				</n-grid-item>
			</n-grid>
		</n-checkbox-group>

		<n-divider />

		<n-button secondary type="success" block class="mt-3.5" @click="save">
			{{ t('sharedButtons.save') }}
		</n-button>
	</n-form>
</template>

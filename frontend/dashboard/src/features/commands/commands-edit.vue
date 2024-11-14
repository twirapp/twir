<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { h, onMounted, ref, toRaw } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'

import { formSchema, useCommandEditV2 } from './composables/use-command-edit-v2'

import FormGeneral from '@/features/commands/ui/form-general.vue'
import FormPermissions from '@/features/commands/ui/form-permissions.vue'
import PageLayout, { type PageLayoutTab } from '@/layout/page-layout.vue'

const route = useRoute()
const { t } = useI18n()

const { findCommand, submit } = useCommandEditV2()

const loading = ref(true)

const { resetForm, handleSubmit, setValues, errors } = useForm({
	validationSchema: toTypedSchema(formSchema),
})

onMounted(async () => {
	resetForm()

	if (typeof route.params.id === 'string') {
		const command = await findCommand(route.params.id)
		if (command) {
			setValues(toRaw(command))
		}
	}

	loading.value = false
})

const onSubmit = handleSubmit(submit)

const tabs: PageLayoutTab[] = [
	{
		name: 'general',
		title: 'General',
		component: h(FormGeneral),
	},
	{
		name: 'permissions',
		title: 'Permissions',
		component: h(FormPermissions),
	},
]
</script>

<template>
	<form
		:class="{ 'blur-sm': loading }"
		@submit="onSubmit"
	>
		<PageLayout :tabs="tabs" :active-tab="tabs[0].name">
			<template #title>
				create/edit
			</template>
		</PageLayout>
	</form>
</template>

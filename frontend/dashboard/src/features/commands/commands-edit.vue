<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { computed, onMounted, ref, toRaw } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'

import { formSchema, useCommandEditV2 } from './composables/use-command-edit-v2'

import { Button } from '@/components/ui/button'
import FormConditions from '@/features/commands/ui/form/form-conditions.vue'
import FormCooldown from '@/features/commands/ui/form/form-cooldown.vue'
import FormExpiration from '@/features/commands/ui/form/form-expiration.vue'
import FormGeneral from '@/features/commands/ui/form/form-general.vue'
import FormPermissions from '@/features/commands/ui/form/form-permissions.vue'
import FormResponses from '@/features/commands/ui/form/form-responses.vue'
import PageLayout from '@/layout/page-layout.vue'

const route = useRoute()
const { t } = useI18n()

const { findCommand, submit } = useCommandEditV2()

const loading = ref(true)

const { handleSubmit, setValues, values } = useForm({
	validationSchema: toTypedSchema(formSchema),
	initialValues: {
		enabled: true,
		aliases: [],
		responses: [
			{
				text: '',
				twitchCategoriesIds: [],
			},
		],
		description: '',
		rolesIds: [],
		deniedUsersIds: [],
		allowedUsersIds: [],
		requiredMessages: 0,
		requiredUsedChannelPoints: 0,
		requiredWatchTime: 0,
		cooldown: 0,
		cooldownType: 'GLOBAL',
		cooldownRolesIds: [],
		isReply: true,
		visible: true,
		keepResponsesOrder: true,
		onlineOnly: false,
		enabledCategories: [],
		expiresType: null,
		expiresAt: null,
	},
	keepValuesOnUnmount: true,
})

const title = ref('')

onMounted(async () => {
	if (typeof route.query.commandIdForCopy === 'string') {
		const command = await findCommand(route.query.commandIdForCopy)
		if (command) {
			setValues(toRaw({
				...command,
				id: undefined,
				module: undefined,
				name: '',
				responses: command.responses.map(r => ({ text: r.text, twitchCategoriesIds: [] })),
				aliases: [],
			}))
			loading.value = false
			return
		}
	}

	if (typeof route.params.id !== 'string') {
		return
	}

	const command = await findCommand(route.params.id)
	if (command) {
		setValues(toRaw(command))
		title.value = command.name
	}

	loading.value = false
})

const onSubmit = handleSubmit(submit)

const backButton = computed(() => {
	if (values.module === 'CUSTOM') {
		return '/dashboard/commands/custom'
	}

	return '/dashboard/commands/builtin'
})
</script>

<template>
	<form :class="{ 'blur-sm': loading }" @submit="onSubmit">
		<PageLayout stickyHeader show-back :back-redirect-to="backButton">
			<template #title>
				<span v-if="route.params.id === 'create'">Create</span>
				<span v-else>Edit "{{ title }}"</span>
			</template>

			<template #action>
				<Button type="submit" :loading="loading">
					{{ t('sharedButtons.save') }}
				</Button>
			</template>

			<template #content>
				<div class="flex flex-col gap-4">
					<FormGeneral />
					<FormResponses />
					<FormCooldown />
					<FormConditions />
					<FormPermissions />
					<FormExpiration />
				</div>
			</template>
		</PageLayout>
	</form>
</template>
